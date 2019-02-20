// package: meta
// file: meta.proto

import * as jspb from "google-protobuf";

export class PartitionId extends jspb.Message {
  getNamespace(): string;
  setNamespace(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PartitionId.AsObject;
  static toObject(includeInstance: boolean, msg: PartitionId): PartitionId.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: PartitionId, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PartitionId;
  static deserializeBinaryFromReader(message: PartitionId, reader: jspb.BinaryReader): PartitionId;
}

export namespace PartitionId {
  export type AsObject = {
    namespace: string,
  }
}

export class PathElement extends jspb.Message {
  getKind(): string;
  setKind(value: string): void;

  hasId(): boolean;
  clearId(): void;
  getId(): number;
  setId(value: number): void;

  hasName(): boolean;
  clearName(): void;
  getName(): string;
  setName(value: string): void;

  getIdtypeCase(): PathElement.IdtypeCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PathElement.AsObject;
  static toObject(includeInstance: boolean, msg: PathElement): PathElement.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: PathElement, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PathElement;
  static deserializeBinaryFromReader(message: PathElement, reader: jspb.BinaryReader): PathElement;
}

export namespace PathElement {
  export type AsObject = {
    kind: string,
    id: number,
    name: string,
  }

  export enum IdtypeCase {
    IDTYPE_NOT_SET = 0,
    ID = 2,
    NAME = 3,
  }
}

export class Key extends jspb.Message {
  hasPartitionid(): boolean;
  clearPartitionid(): void;
  getPartitionid(): PartitionId | undefined;
  setPartitionid(value?: PartitionId): void;

  clearPathList(): void;
  getPathList(): Array<PathElement>;
  setPathList(value: Array<PathElement>): void;
  addPath(value?: PathElement, index?: number): PathElement;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Key.AsObject;
  static toObject(includeInstance: boolean, msg: Key): Key.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Key, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Key;
  static deserializeBinaryFromReader(message: Key, reader: jspb.BinaryReader): Key;
}

export namespace Key {
  export type AsObject = {
    partitionid?: PartitionId.AsObject,
    pathList: Array<PathElement.AsObject>,
  }
}

export class Value extends jspb.Message {
  getId(): number;
  setId(value: number): void;

  getType(): ValueType;
  setType(value: ValueType): void;

  getDoublevalue(): number;
  setDoublevalue(value: number): void;

  getInt64value(): number;
  setInt64value(value: number): void;

  getStringvalue(): string;
  setStringvalue(value: string): void;

  getTimestampvalue(): Uint8Array | string;
  getTimestampvalue_asU8(): Uint8Array;
  getTimestampvalue_asB64(): string;
  setTimestampvalue(value: Uint8Array | string): void;

  getBooleanvalue(): boolean;
  setBooleanvalue(value: boolean): void;

  getBytesvalue(): Uint8Array | string;
  getBytesvalue_asU8(): Uint8Array;
  getBytesvalue_asB64(): string;
  setBytesvalue(value: Uint8Array | string): void;

  hasKeyvalue(): boolean;
  clearKeyvalue(): void;
  getKeyvalue(): Key | undefined;
  setKeyvalue(value?: Key): void;

  getUint64value(): number;
  setUint64value(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Value.AsObject;
  static toObject(includeInstance: boolean, msg: Value): Value.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Value, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Value;
  static deserializeBinaryFromReader(message: Value, reader: jspb.BinaryReader): Value;
}

export namespace Value {
  export type AsObject = {
    id: number,
    type: ValueType,
    doublevalue: number,
    int64value: number,
    stringvalue: string,
    timestampvalue: Uint8Array | string,
    booleanvalue: boolean,
    bytesvalue: Uint8Array | string,
    keyvalue?: Key.AsObject,
    uint64value: number,
  }
}

export class SchemaField extends jspb.Message {
  getId(): number;
  setId(value: number): void;

  getName(): string;
  setName(value: string): void;

  getType(): ValueType;
  setType(value: ValueType): void;

  getComment(): string;
  setComment(value: string): void;

  hasEditor(): boolean;
  clearEditor(): void;
  getEditor(): SchemaFieldEditorInfo | undefined;
  setEditor(value?: SchemaFieldEditorInfo): void;

  getReadonly(): boolean;
  setReadonly(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SchemaField.AsObject;
  static toObject(includeInstance: boolean, msg: SchemaField): SchemaField.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: SchemaField, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SchemaField;
  static deserializeBinaryFromReader(message: SchemaField, reader: jspb.BinaryReader): SchemaField;
}

export namespace SchemaField {
  export type AsObject = {
    id: number,
    name: string,
    type: ValueType,
    comment: string,
    editor?: SchemaFieldEditorInfo.AsObject,
    readonly: boolean,
  }
}

export class SchemaFieldEditorInfo extends jspb.Message {
  getDisplayname(): string;
  setDisplayname(value: string): void;

  getType(): SchemaFieldEditorInfoType;
  setType(value: SchemaFieldEditorInfoType): void;

  getEditorreadonly(): boolean;
  setEditorreadonly(value: boolean): void;

  getForeigntype(): string;
  setForeigntype(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SchemaFieldEditorInfo.AsObject;
  static toObject(includeInstance: boolean, msg: SchemaFieldEditorInfo): SchemaFieldEditorInfo.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: SchemaFieldEditorInfo, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SchemaFieldEditorInfo;
  static deserializeBinaryFromReader(message: SchemaFieldEditorInfo, reader: jspb.BinaryReader): SchemaFieldEditorInfo;
}

export namespace SchemaFieldEditorInfo {
  export type AsObject = {
    displayname: string,
    type: SchemaFieldEditorInfoType,
    editorreadonly: boolean,
    foreigntype: string,
  }
}

export class SchemaKindEditor extends jspb.Message {
  getSingular(): string;
  setSingular(value: string): void;

  getPlural(): string;
  setPlural(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SchemaKindEditor.AsObject;
  static toObject(includeInstance: boolean, msg: SchemaKindEditor): SchemaKindEditor.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: SchemaKindEditor, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SchemaKindEditor;
  static deserializeBinaryFromReader(message: SchemaKindEditor, reader: jspb.BinaryReader): SchemaKindEditor;
}

export namespace SchemaKindEditor {
  export type AsObject = {
    singular: string,
    plural: string,
  }
}

export class SchemaKind extends jspb.Message {
  clearFieldsList(): void;
  getFieldsList(): Array<SchemaField>;
  setFieldsList(value: Array<SchemaField>): void;
  addFields(value?: SchemaField, index?: number): SchemaField;

  hasEditor(): boolean;
  clearEditor(): void;
  getEditor(): SchemaKindEditor | undefined;
  setEditor(value?: SchemaKindEditor): void;

  clearIndexesList(): void;
  getIndexesList(): Array<SchemaIndex>;
  setIndexesList(value: Array<SchemaIndex>): void;
  addIndexes(value?: SchemaIndex, index?: number): SchemaIndex;

  clearAncestorsList(): void;
  getAncestorsList(): Array<string>;
  setAncestorsList(value: Array<string>): void;
  addAncestors(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SchemaKind.AsObject;
  static toObject(includeInstance: boolean, msg: SchemaKind): SchemaKind.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: SchemaKind, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SchemaKind;
  static deserializeBinaryFromReader(message: SchemaKind, reader: jspb.BinaryReader): SchemaKind;
}

export namespace SchemaKind {
  export type AsObject = {
    fieldsList: Array<SchemaField.AsObject>,
    editor?: SchemaKindEditor.AsObject,
    indexesList: Array<SchemaIndex.AsObject>,
    ancestorsList: Array<string>,
  }
}

export class SchemaIndex extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getType(): SchemaIndexType;
  setType(value: SchemaIndexType): void;

  hasComputed(): boolean;
  clearComputed(): void;
  getComputed(): SchemaComputedIndex | undefined;
  setComputed(value?: SchemaComputedIndex): void;

  hasField(): boolean;
  clearField(): void;
  getField(): string;
  setField(value: string): void;

  getValueCase(): SchemaIndex.ValueCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SchemaIndex.AsObject;
  static toObject(includeInstance: boolean, msg: SchemaIndex): SchemaIndex.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: SchemaIndex, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SchemaIndex;
  static deserializeBinaryFromReader(message: SchemaIndex, reader: jspb.BinaryReader): SchemaIndex;
}

export namespace SchemaIndex {
  export type AsObject = {
    name: string,
    type: SchemaIndexType,
    computed?: SchemaComputedIndex.AsObject,
    field: string,
  }

  export enum ValueCase {
    VALUE_NOT_SET = 0,
    COMPUTED = 3,
    FIELD = 4,
  }
}

export class SchemaComputedIndex extends jspb.Message {
  hasFnv64a(): boolean;
  clearFnv64a(): void;
  getFnv64a(): SchemaComputedIndexFnv64a | undefined;
  setFnv64a(value?: SchemaComputedIndexFnv64a): void;

  hasFnv64aPair(): boolean;
  clearFnv64aPair(): void;
  getFnv64aPair(): SchemaComputedIndexFnv64aPair | undefined;
  setFnv64aPair(value?: SchemaComputedIndexFnv64aPair): void;

  getAlgorithmCase(): SchemaComputedIndex.AlgorithmCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SchemaComputedIndex.AsObject;
  static toObject(includeInstance: boolean, msg: SchemaComputedIndex): SchemaComputedIndex.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: SchemaComputedIndex, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SchemaComputedIndex;
  static deserializeBinaryFromReader(message: SchemaComputedIndex, reader: jspb.BinaryReader): SchemaComputedIndex;
}

export namespace SchemaComputedIndex {
  export type AsObject = {
    fnv64a?: SchemaComputedIndexFnv64a.AsObject,
    fnv64aPair?: SchemaComputedIndexFnv64aPair.AsObject,
  }

  export enum AlgorithmCase {
    ALGORITHM_NOT_SET = 0,
    FNV64A = 1,
    FNV64A_PAIR = 2,
  }
}

export class SchemaComputedIndexFnv64a extends jspb.Message {
  getField(): string;
  setField(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SchemaComputedIndexFnv64a.AsObject;
  static toObject(includeInstance: boolean, msg: SchemaComputedIndexFnv64a): SchemaComputedIndexFnv64a.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: SchemaComputedIndexFnv64a, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SchemaComputedIndexFnv64a;
  static deserializeBinaryFromReader(message: SchemaComputedIndexFnv64a, reader: jspb.BinaryReader): SchemaComputedIndexFnv64a;
}

export namespace SchemaComputedIndexFnv64a {
  export type AsObject = {
    field: string,
  }
}

export class SchemaComputedIndexFnv64aPair extends jspb.Message {
  getField1(): string;
  setField1(value: string): void;

  getField2(): string;
  setField2(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SchemaComputedIndexFnv64aPair.AsObject;
  static toObject(includeInstance: boolean, msg: SchemaComputedIndexFnv64aPair): SchemaComputedIndexFnv64aPair.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: SchemaComputedIndexFnv64aPair, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SchemaComputedIndexFnv64aPair;
  static deserializeBinaryFromReader(message: SchemaComputedIndexFnv64aPair, reader: jspb.BinaryReader): SchemaComputedIndexFnv64aPair;
}

export namespace SchemaComputedIndexFnv64aPair {
  export type AsObject = {
    field1: string,
    field2: string,
  }
}

export class Schema extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getKindsMap(): jspb.Map<string, SchemaKind>;
  clearKindsMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Schema.AsObject;
  static toObject(includeInstance: boolean, msg: Schema): Schema.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Schema, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Schema;
  static deserializeBinaryFromReader(message: Schema, reader: jspb.BinaryReader): Schema;
}

export namespace Schema {
  export type AsObject = {
    name: string,
    kindsMap: Array<[string, SchemaKind.AsObject]>,
  }
}

export class GetSchemaRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetSchemaRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetSchemaRequest): GetSchemaRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetSchemaRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetSchemaRequest;
  static deserializeBinaryFromReader(message: GetSchemaRequest, reader: jspb.BinaryReader): GetSchemaRequest;
}

export namespace GetSchemaRequest {
  export type AsObject = {
  }
}

export class GetSchemaResponse extends jspb.Message {
  hasSchema(): boolean;
  clearSchema(): void;
  getSchema(): Schema | undefined;
  setSchema(value?: Schema): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetSchemaResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetSchemaResponse): GetSchemaResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetSchemaResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetSchemaResponse;
  static deserializeBinaryFromReader(message: GetSchemaResponse, reader: jspb.BinaryReader): GetSchemaResponse;
}

export namespace GetSchemaResponse {
  export type AsObject = {
    schema?: Schema.AsObject,
  }
}

export class MetaListEntitiesRequest extends jspb.Message {
  getStart(): Uint8Array | string;
  getStart_asU8(): Uint8Array;
  getStart_asB64(): string;
  setStart(value: Uint8Array | string): void;

  getLimit(): number;
  setLimit(value: number): void;

  getKindname(): string;
  setKindname(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MetaListEntitiesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: MetaListEntitiesRequest): MetaListEntitiesRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: MetaListEntitiesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MetaListEntitiesRequest;
  static deserializeBinaryFromReader(message: MetaListEntitiesRequest, reader: jspb.BinaryReader): MetaListEntitiesRequest;
}

export namespace MetaListEntitiesRequest {
  export type AsObject = {
    start: Uint8Array | string,
    limit: number,
    kindname: string,
  }
}

export class MetaListEntitiesResponse extends jspb.Message {
  getNext(): Uint8Array | string;
  getNext_asU8(): Uint8Array;
  getNext_asB64(): string;
  setNext(value: Uint8Array | string): void;

  getMoreresults(): boolean;
  setMoreresults(value: boolean): void;

  clearEntitiesList(): void;
  getEntitiesList(): Array<MetaEntity>;
  setEntitiesList(value: Array<MetaEntity>): void;
  addEntities(value?: MetaEntity, index?: number): MetaEntity;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MetaListEntitiesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: MetaListEntitiesResponse): MetaListEntitiesResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: MetaListEntitiesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MetaListEntitiesResponse;
  static deserializeBinaryFromReader(message: MetaListEntitiesResponse, reader: jspb.BinaryReader): MetaListEntitiesResponse;
}

export namespace MetaListEntitiesResponse {
  export type AsObject = {
    next: Uint8Array | string,
    moreresults: boolean,
    entitiesList: Array<MetaEntity.AsObject>,
  }
}

export class MetaEntity extends jspb.Message {
  hasKey(): boolean;
  clearKey(): void;
  getKey(): Key | undefined;
  setKey(value?: Key): void;

  clearValuesList(): void;
  getValuesList(): Array<Value>;
  setValuesList(value: Array<Value>): void;
  addValues(value?: Value, index?: number): Value;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MetaEntity.AsObject;
  static toObject(includeInstance: boolean, msg: MetaEntity): MetaEntity.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: MetaEntity, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MetaEntity;
  static deserializeBinaryFromReader(message: MetaEntity, reader: jspb.BinaryReader): MetaEntity;
}

export namespace MetaEntity {
  export type AsObject = {
    key?: Key.AsObject,
    valuesList: Array<Value.AsObject>,
  }
}

export enum ValueType {
  UNKNOWN = 0,
  DOUBLE = 1,
  INT64 = 2,
  STRING = 3,
  TIMESTAMP = 4,
  BOOLEAN = 5,
  BYTES = 6,
  KEY = 7,
  UINT64 = 8,
}

export enum SchemaFieldEditorInfoType {
  DEFAULT = 0,
  PASSWORD = 1,
  LOOKUP = 2,
}

export enum SchemaIndexType {
  UNSPECIFIED = 0,
  MEMORY = 1,
}

