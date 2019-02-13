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
  }
}

export class Field extends jspb.Message {
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
  getEditor(): FieldEditorInfo | undefined;
  setEditor(value?: FieldEditorInfo): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Field.AsObject;
  static toObject(includeInstance: boolean, msg: Field): Field.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Field, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Field;
  static deserializeBinaryFromReader(message: Field, reader: jspb.BinaryReader): Field;
}

export namespace Field {
  export type AsObject = {
    id: number,
    name: string,
    type: ValueType,
    comment: string,
    editor?: FieldEditorInfo.AsObject,
  }
}

export class FieldEditorInfo extends jspb.Message {
  getDisplayname(): string;
  setDisplayname(value: string): void;

  getType(): FieldEditorInfoType;
  setType(value: FieldEditorInfoType): void;

  getReadonly(): boolean;
  setReadonly(value: boolean): void;

  getForeigntype(): string;
  setForeigntype(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FieldEditorInfo.AsObject;
  static toObject(includeInstance: boolean, msg: FieldEditorInfo): FieldEditorInfo.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: FieldEditorInfo, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FieldEditorInfo;
  static deserializeBinaryFromReader(message: FieldEditorInfo, reader: jspb.BinaryReader): FieldEditorInfo;
}

export namespace FieldEditorInfo {
  export type AsObject = {
    displayname: string,
    type: FieldEditorInfoType,
    readonly: boolean,
    foreigntype: string,
  }
}

export class KindEditor extends jspb.Message {
  getSingular(): string;
  setSingular(value: string): void;

  getPlural(): string;
  setPlural(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): KindEditor.AsObject;
  static toObject(includeInstance: boolean, msg: KindEditor): KindEditor.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: KindEditor, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): KindEditor;
  static deserializeBinaryFromReader(message: KindEditor, reader: jspb.BinaryReader): KindEditor;
}

export namespace KindEditor {
  export type AsObject = {
    singular: string,
    plural: string,
  }
}

export class Kind extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  clearFieldsList(): void;
  getFieldsList(): Array<Field>;
  setFieldsList(value: Array<Field>): void;
  addFields(value?: Field, index?: number): Field;

  hasEditor(): boolean;
  clearEditor(): void;
  getEditor(): KindEditor | undefined;
  setEditor(value?: KindEditor): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Kind.AsObject;
  static toObject(includeInstance: boolean, msg: Kind): Kind.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Kind, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Kind;
  static deserializeBinaryFromReader(message: Kind, reader: jspb.BinaryReader): Kind;
}

export namespace Kind {
  export type AsObject = {
    name: string,
    fieldsList: Array<Field.AsObject>,
    editor?: KindEditor.AsObject,
  }
}

export class Schema extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  clearKindsList(): void;
  getKindsList(): Array<Kind>;
  setKindsList(value: Array<Kind>): void;
  addKinds(value?: Kind, index?: number): Kind;

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
    kindsList: Array<Kind.AsObject>,
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
  TYPEDOUBLE = 0,
  TYPEINT64 = 1,
  TYPESTRING = 2,
  TYPETIMESTAMP = 3,
  TYPEBOOLEAN = 4,
  TYPEBYTES = 5,
  TYPEKEY = 6,
}

export enum FieldEditorInfoType {
  DEFAULT = 0,
  PASSWORD = 1,
  LOOKUP = 2,
}

