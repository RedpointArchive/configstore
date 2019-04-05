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

type ConfigstoreMetaServiceMetaGet = {
  readonly methodName: string;
  readonly service: typeof ConfigstoreMetaService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof meta_pb.MetaGetEntityRequest;
  readonly responseType: typeof meta_pb.MetaGetEntityResponse;
};

type ConfigstoreMetaServiceMetaUpdate = {
  readonly methodName: string;
  readonly service: typeof ConfigstoreMetaService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof meta_pb.MetaUpdateEntityRequest;
  readonly responseType: typeof meta_pb.MetaUpdateEntityResponse;
};

type ConfigstoreMetaServiceMetaCreate = {
  readonly methodName: string;
  readonly service: typeof ConfigstoreMetaService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof meta_pb.MetaCreateEntityRequest;
  readonly responseType: typeof meta_pb.MetaCreateEntityResponse;
};

type ConfigstoreMetaServiceMetaDelete = {
  readonly methodName: string;
  readonly service: typeof ConfigstoreMetaService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof meta_pb.MetaDeleteEntityRequest;
  readonly responseType: typeof meta_pb.MetaDeleteEntityResponse;
};

type ConfigstoreMetaServiceGetDefaultPartitionId = {
  readonly methodName: string;
  readonly service: typeof ConfigstoreMetaService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof meta_pb.GetDefaultPartitionIdRequest;
  readonly responseType: typeof meta_pb.GetDefaultPartitionIdResponse;
};

export class ConfigstoreMetaService {
  static readonly serviceName: string;
  static readonly GetSchema: ConfigstoreMetaServiceGetSchema;
  static readonly MetaList: ConfigstoreMetaServiceMetaList;
  static readonly MetaGet: ConfigstoreMetaServiceMetaGet;
  static readonly MetaUpdate: ConfigstoreMetaServiceMetaUpdate;
  static readonly MetaCreate: ConfigstoreMetaServiceMetaCreate;
  static readonly MetaDelete: ConfigstoreMetaServiceMetaDelete;
  static readonly GetDefaultPartitionId: ConfigstoreMetaServiceGetDefaultPartitionId;
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
  metaGet(
    requestMessage: meta_pb.MetaGetEntityRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: meta_pb.MetaGetEntityResponse|null) => void
  ): UnaryResponse;
  metaGet(
    requestMessage: meta_pb.MetaGetEntityRequest,
    callback: (error: ServiceError|null, responseMessage: meta_pb.MetaGetEntityResponse|null) => void
  ): UnaryResponse;
  metaUpdate(
    requestMessage: meta_pb.MetaUpdateEntityRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: meta_pb.MetaUpdateEntityResponse|null) => void
  ): UnaryResponse;
  metaUpdate(
    requestMessage: meta_pb.MetaUpdateEntityRequest,
    callback: (error: ServiceError|null, responseMessage: meta_pb.MetaUpdateEntityResponse|null) => void
  ): UnaryResponse;
  metaCreate(
    requestMessage: meta_pb.MetaCreateEntityRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: meta_pb.MetaCreateEntityResponse|null) => void
  ): UnaryResponse;
  metaCreate(
    requestMessage: meta_pb.MetaCreateEntityRequest,
    callback: (error: ServiceError|null, responseMessage: meta_pb.MetaCreateEntityResponse|null) => void
  ): UnaryResponse;
  metaDelete(
    requestMessage: meta_pb.MetaDeleteEntityRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: meta_pb.MetaDeleteEntityResponse|null) => void
  ): UnaryResponse;
  metaDelete(
    requestMessage: meta_pb.MetaDeleteEntityRequest,
    callback: (error: ServiceError|null, responseMessage: meta_pb.MetaDeleteEntityResponse|null) => void
  ): UnaryResponse;
  getDefaultPartitionId(
    requestMessage: meta_pb.GetDefaultPartitionIdRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: meta_pb.GetDefaultPartitionIdResponse|null) => void
  ): UnaryResponse;
  getDefaultPartitionId(
    requestMessage: meta_pb.GetDefaultPartitionIdRequest,
    callback: (error: ServiceError|null, responseMessage: meta_pb.GetDefaultPartitionIdResponse|null) => void
  ): UnaryResponse;
}

