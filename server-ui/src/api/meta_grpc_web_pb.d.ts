import * as grpcWeb from 'grpc-web';

import {
  GetDefaultPartitionIdRequest,
  GetDefaultPartitionIdResponse,
  GetSchemaRequest,
  GetSchemaResponse,
  MetaCreateEntityRequest,
  MetaCreateEntityResponse,
  MetaDeleteEntityRequest,
  MetaDeleteEntityResponse,
  MetaGetEntityRequest,
  MetaGetEntityResponse,
  MetaListEntitiesRequest,
  MetaListEntitiesResponse,
  MetaUpdateEntityRequest,
  MetaUpdateEntityResponse} from './meta_pb';

export class ConfigstoreMetaServiceClient {
  constructor (hostname: string,
               credentials: null | { [index: string]: string; },
               options: null | { [index: string]: string; });

  getSchema(
    request: GetSchemaRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: GetSchemaResponse) => void
  ): grpcWeb.ClientReadableStream<GetSchemaResponse>;

  metaList(
    request: MetaListEntitiesRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: MetaListEntitiesResponse) => void
  ): grpcWeb.ClientReadableStream<MetaListEntitiesResponse>;

  metaGet(
    request: MetaGetEntityRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: MetaGetEntityResponse) => void
  ): grpcWeb.ClientReadableStream<MetaGetEntityResponse>;

  metaUpdate(
    request: MetaUpdateEntityRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: MetaUpdateEntityResponse) => void
  ): grpcWeb.ClientReadableStream<MetaUpdateEntityResponse>;

  metaCreate(
    request: MetaCreateEntityRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: MetaCreateEntityResponse) => void
  ): grpcWeb.ClientReadableStream<MetaCreateEntityResponse>;

  metaDelete(
    request: MetaDeleteEntityRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: MetaDeleteEntityResponse) => void
  ): grpcWeb.ClientReadableStream<MetaDeleteEntityResponse>;

  getDefaultPartitionId(
    request: GetDefaultPartitionIdRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: GetDefaultPartitionIdResponse) => void
  ): grpcWeb.ClientReadableStream<GetDefaultPartitionIdResponse>;

}

export class ConfigstoreMetaServicePromiseClient {
  constructor (hostname: string,
               credentials: null | { [index: string]: string; },
               options: null | { [index: string]: string; });

  getSchema(
    request: GetSchemaRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<GetSchemaResponse>;

  metaList(
    request: MetaListEntitiesRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<MetaListEntitiesResponse>;

  metaGet(
    request: MetaGetEntityRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<MetaGetEntityResponse>;

  metaUpdate(
    request: MetaUpdateEntityRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<MetaUpdateEntityResponse>;

  metaCreate(
    request: MetaCreateEntityRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<MetaCreateEntityResponse>;

  metaDelete(
    request: MetaDeleteEntityRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<MetaDeleteEntityResponse>;

  getDefaultPartitionId(
    request: GetDefaultPartitionIdRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<GetDefaultPartitionIdResponse>;

}

