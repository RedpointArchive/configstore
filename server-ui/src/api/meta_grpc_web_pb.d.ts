import * as grpcWeb from 'grpc-web';
import {
  GetSchemaRequest,
  GetSchemaResponse,
  Key,
  MetaEntity,
  MetaListEntitiesRequest,
  MetaListEntitiesResponse,
  PartitionId,
  PathElement,
  Schema,
  KindsEntry,
  SchemaComputedIndex,
  SchemaComputedIndexFnv64a,
  SchemaComputedIndexFnv64aPair,
  SchemaField,
  SchemaFieldEditorInfo,
  SchemaIndex,
  SchemaKind,
  SchemaKindEditor,
  Value} from './meta_pb';

export class ConfigstoreMetaServiceClient {
  constructor (hostname: string,
               credentials: null | { [index: string]: string; },
               options: null | { [index: string]: string; });

  getSchema(
    request: GetSchemaRequest,
    metadata: grpcWeb.Metadata,
    callback: (err: grpcWeb.Error,
               response: GetSchemaResponse) => void
  ): grpcWeb.ClientReadableStream<GetSchemaResponse>;

  metaList(
    request: MetaListEntitiesRequest,
    metadata: grpcWeb.Metadata,
    callback: (err: grpcWeb.Error,
               response: MetaListEntitiesResponse) => void
  ): grpcWeb.ClientReadableStream<MetaListEntitiesResponse>;

}

export class ConfigstoreMetaServicePromiseClient {
  constructor (hostname: string,
               credentials: null | { [index: string]: string; },
               options: null | { [index: string]: string; });

  getSchema(
    request: GetSchemaRequest,
    metadata: grpcWeb.Metadata
  ): Promise<GetSchemaResponse>;

  metaList(
    request: MetaListEntitiesRequest,
    metadata: grpcWeb.Metadata
  ): Promise<MetaListEntitiesResponse>;

}

