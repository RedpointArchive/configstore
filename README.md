# configstore

_configstore is a Go server that lets you define your schema in JSON, and automatically generates a gRPC API with List/Get/Create/Update/Delete methods to access and modify those entities in Google Cloud Firestore._

Basically, configstore is a server that provides a gRPC CRUD interface for the kinds of entities that you describe. You just specify the schema for kinds as JSON, and configstore produces a `.proto` file you can use to generate clients in various languages.

configstore is entirely stateless; all data is stored on Google Cloud Firestore. This makes it easier to run in highly available configurations, without managing the state yourself.

## Usage

configstore is a very early prototype, and isn't ready for production usage.

You can bulid and test the image by running

```
docker build . --tag=configstore
```

You can then run configstore with an invocation similar to the following:

```
docker run --rm -p 13389:13389 -p 13390:13390 -v your_schema.json:/schema.json -e CONFIGSTORE_GOOGLE_CLOUD_PROJECT_ID="your-cloud-project" -e CONFIGSTORE_GRPC_PORT=13389 -e CONFIGSTORE_HTTP_PORT=13390 -e CONFIGSTORE_SCHEMA_PATH="/schema.json" -v your_service_account.json:/adc.json -e GOOGLE_APPLICATION_CREDENTIALS=/adc.json --name=configstore configstore
```

## Screenshots

![server](https://github.com/hach-que/configstore/raw/master/screenshots/server.PNG)

![client](https://github.com/hach-que/configstore/raw/master/screenshots/client.PNG)

## License

```
MIT License

Copyright (c) 2019 June Rhodes

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```
