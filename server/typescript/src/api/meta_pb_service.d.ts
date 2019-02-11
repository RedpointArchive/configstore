// package: meta
// file: meta.proto

import * as meta_pb from "./meta_pb";
import {grpc} from "@improbable-eng/grpc-web";

type ConfigstoreMetaServiceGetSchema = {
  readonly methodName: string;
  readonly service: typeof ConfigstoreMetaService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof meta_pb.GetSchemaRequest;
  readonly responseType: typeof meta_pb.GetSchemaResponse;
};

type ConfigstoreMetaServiceMetaList = {
  readonly methodName: string;
  readonly service: typeof ConfigstoreMetaService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof meta_pb.MetaListEntitiesRequest;
  readonly responseType: typeof meta_pb.MetaListEntitiesResponse;
};

export class ConfigstoreMetaService {
  static readonly serviceName: string;
  static readonly GetSchema: ConfigstoreMetaServiceGetSchema;
  static readonly MetaList: ConfigstoreMetaServiceMetaList;
}

export type ServiceError = { message: string, code: number; metadata: grpc.Metadata }
export type Status = { details: string, code: number; metadata: grpc.Metadata }

interface UnaryResponse {
  cancel(): void;
}
interface ResponseStream<T> {
  cancel(): void;
  on(type: 'data', handler: (message: T) => void): ResponseStream<T>;
  on(type: 'end', handler: () => void): ResponseStream<T>;
  on(type: 'status', handler: (status: Status) => void): ResponseStream<T>;
}
interface RequestStream<T> {
  write(message: T): RequestStream<T>;
  end(): void;
  cancel(): void;
  on(type: 'end', handler: () => void): RequestStream<T>;
  on(type: 'status', handler: (status: Status) => void): RequestStream<T>;
}
interface BidirectionalStream<ReqT, ResT> {
  write(message: ReqT): BidirectionalStream<ReqT, ResT>;
  end(): void;
  cancel(): void;
  on(type: 'data', handler: (message: ResT) => void): BidirectionalStream<ReqT, ResT>;
  on(type: 'end', handler: () => void): BidirectionalStream<ReqT, ResT>;
  on(type: 'status', handler: (status: Status) => void): BidirectionalStream<ReqT, ResT>;
}

export class ConfigstoreMetaServiceClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  getSchema(
    requestMessage: meta_pb.GetSchemaRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: meta_pb.GetSchemaResponse|null) => void
  ): UnaryResponse;
  getSchema(
    requestMessage: meta_pb.GetSchemaRequest,
    callback: (error: ServiceError|null, responseMessage: meta_pb.GetSchemaResponse|null) => void
  ): UnaryResponse;
  metaList(
    requestMessage: meta_pb.MetaListEntitiesRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: meta_pb.MetaListEntitiesResponse|null) => void
  ): UnaryResponse;
  metaList(
    requestMessage: meta_pb.MetaListEntitiesRequest,
    callback: (error: ServiceError|null, responseMessage: meta_pb.MetaListEntitiesResponse|null) => void
  ): UnaryResponse;
}

