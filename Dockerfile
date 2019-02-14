# install Go
FROM ubuntu:18.04 AS go_install

ENV GOROOT=/usr/local/go
ENV GOPATH=/go
ENV PATH="/go/bin:/usr/local/go/bin:${PATH}"

RUN apt update
RUN apt install -y wget
RUN apt install -y git

RUN cd /tmp && \
  wget https://dl.google.com/go/go1.11.linux-amd64.tar.gz && \
  tar -xvf go1.11.linux-amd64.tar.gz && \
  mv go /usr/local

RUN go get -u github.com/golang/protobuf/protoc-gen-go

# install protoc
FROM ubuntu:18.04 AS protoc_install

RUN apt update
RUN apt install -y curl
RUN apt install -y unzip

WORKDIR /tmp
RUN mkdir /protoc
RUN curl -OL https://github.com/google/protobuf/releases/download/v3.6.0/protoc-3.6.0-linux-x86_64.zip
RUN unzip protoc-3.6.0-linux-x86_64.zip -d /protoc

# generate meta protocol buffers for Go
FROM ubuntu:18.04 AS meta_generate_go

ENV GOROOT=/usr/local/go
ENV GOPATH=/go
ENV PATH="/go/bin:/usr/local/go/bin:${PATH}"

COPY --from=protoc_install /protoc /protoc
COPY --from=go_install /go /go
COPY /server/meta.proto /meta.proto
RUN mkdir /meta
RUN /protoc/bin/protoc --go_out=plugins=grpc:meta ./meta.proto
RUN sed -i 's/package meta/package main/' /meta/meta.pb.go

# generate meta protocol buffers for TypeScript
FROM ubuntu:18.04 AS meta_generate_ts

RUN apt update
RUN apt install -y curl
RUN apt install -y unzip
RUN apt install -y gnupg2
RUN curl -sL https://deb.nodesource.com/setup_10.x | bash -
RUN apt-get install -y nodejs
RUN curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | apt-key add -
RUN echo "deb https://dl.yarnpkg.com/debian/ stable main" | tee /etc/apt/sources.list.d/yarn.list
RUN apt-get update && apt-get install yarn

RUN yarn global add ts-protoc-gen google-protobuf

COPY --from=protoc_install /protoc /protoc
COPY /server/meta.proto /meta.proto
WORKDIR /
RUN mkdir -p /src/api
RUN /protoc/bin/protoc \
  --plugin="protoc-gen-ts=$(yarn global bin)/protoc-gen-ts" \
  --js_out="import_style=commonjs,binary:src/api" \
  --ts_out="service=true:src/api" \
  meta.proto

# build server
FROM ubuntu:18.04 AS build_server

ENV GOROOT=/usr/local/go
ENV GOPATH=/go
ENV PATH="/go/bin:/usr/local/go/bin:${PATH}"

COPY --from=go_install /usr/local/go /usr/local/go
COPY --from=go_install /go /go

COPY server /src
COPY --from=meta_generate_go /meta/meta.pb.go /src/meta.pb.go
WORKDIR /src
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -mod vendor -o /server

# prerun client
FROM ubuntu:18.04 AS prerun_client

ENV CONFIGSTORE_GOOGLE_CLOUD_PROJECT_ID="configstore-test-001"
ENV CONFIGSTORE_GRPC_PORT="13389"
ENV CONFIGSTORE_HTTP_PORT="13390"
ENV CONFIGSTORE_SCHEMA_PATH="/schema.json"

COPY --from=build_server /server /server
COPY server/schema.json /schema.json
RUN /server -generate > /client.go

COPY client /client-src
RUN mv /client.go /client-src/client.go

# build server UI
FROM ubuntu:18.04 AS build_server_ui

RUN apt update
RUN apt install -y curl
RUN apt install -y unzip
RUN apt install -y gnupg2
RUN curl -sL https://deb.nodesource.com/setup_10.x | bash -
RUN apt-get install -y nodejs
RUN curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | apt-key add -
RUN echo "deb https://dl.yarnpkg.com/debian/ stable main" | tee /etc/apt/sources.list.d/yarn.list
RUN apt-get update && apt-get install yarn

COPY server-ui /src
WORKDIR /src
RUN yarn
COPY --from=meta_generate_ts /src/api /src/api
RUN yarn build

# test server & client
FROM google/cloud-sdk:234.0.0 AS test

ENV PATH="/root/google-cloud-sdk/bin:${PATH}"
ENV GOROOT=/usr/local/go
ENV GOPATH=/go
ENV PATH="/go/bin:/usr/local/go/bin:${PATH}"
ENV CONFIGSTORE_GOOGLE_CLOUD_PROJECT_ID="configstore-test-001"
ENV CONFIGSTORE_GRPC_PORT="13389"
ENV CONFIGSTORE_HTTP_PORT="13390"
ENV CONFIGSTORE_SCHEMA_PATH="/schema.json"
ENV FIRESTORE_EMULATOR_HOST=127.0.0.1:8432

RUN apt-get update
RUN apt-get remove -y --purge \
  kubectl \
  google-cloud-sdk \
  google-cloud-sdk-app-engine-grpc \
  google-cloud-sdk-pubsub-emulator \
  google-cloud-sdk-app-engine-go \
  google-cloud-sdk-cloud-build-local \
  google-cloud-sdk-datastore-emulator \
  google-cloud-sdk-app-engine-python \
  google-cloud-sdk-cbt \
  google-cloud-sdk-bigtable-emulator \
  google-cloud-sdk-app-engine-python-extras \
  google-cloud-sdk-datalab \
  google-cloud-sdk-app-engine-java
RUN curl https://sdk.cloud.google.com | bash
RUN gcloud components install beta
RUN gcloud components install cloud-firestore-emulator

COPY --from=go_install /usr/local/go /usr/local/go
COPY --from=go_install /go /go

RUN go get -v \
  "github.com/rs/xid" \
  "github.com/golang/protobuf/ptypes/timestamp" \
  "google.golang.org/grpc" \
  "gotest.tools/assert"
COPY --from=build_server /server /server
COPY --from=prerun_client /schema.json /schema.json
COPY --from=prerun_client /client-src /client-src
COPY extra/adc.json /adc.json
ENV GOOGLE_APPLICATION_CREDENTIALS=/adc.json
RUN gcloud beta emulators firestore start --host-port=127.0.0.1:8432 & \
  sleep 5 && /server & \
  sleep 5 && cd /client-src && go test

# final image
FROM scratch AS final

COPY --from=test /server /server
COPY --from=build_server_ui /src/build /server-ui
ENTRYPOINT [ "/server" ]