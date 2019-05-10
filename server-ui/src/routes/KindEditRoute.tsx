import React, { useState } from "react";
import { RouteComponentProps } from "react-router";
import {
  GetSchemaResponse,
  ValueType,
  MetaGetEntityRequest,
  MetaGetEntityResponse,
  SchemaField,
  Value,
  MetaEntity,
  MetaUpdateEntityRequest,
  MetaCreateEntityRequest,
  SchemaFieldEditorInfo,
  MetaOperation,
  Key,
  PathElement,
  GetDefaultPartitionIdResponse,
  GetDefaultPartitionIdRequest,
  PartitionId
} from "../api/meta_pb";
import { g, deserializeKey, prettifyKey, c, serializeKey } from "../core";
import { Link } from "react-router-dom";
import { useAsync } from "react-async";
import { ConfigstoreMetaServicePromiseClient } from "../api/meta_grpc_web_pb";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faSpinner } from "@fortawesome/free-solid-svg-icons";
import { PendingTransactionContext, PendingTransaction } from "../App";
import { createGrpcPromiseClient } from "../svcHost";
import Datetime from "react-datetime";
import { Timestamp } from "google-protobuf/google/protobuf/timestamp_pb";
import moment from "moment";
import { KeySelect } from "../KeySelect";
import { Address4, Address6 } from "ip-address";

export interface KindEditRouteMatch {
  kind: string;
  id: string;
  idx: string;
}

export interface KindEditRouteProps
  extends RouteComponentProps<KindEditRouteMatch> {
  schema: GetSchemaResponse;
}

const getKind = async (props: any): Promise<MetaGetEntityResponse> => {
  const client = createGrpcPromiseClient(ConfigstoreMetaServicePromiseClient);
  const req = new MetaGetEntityRequest();
  req.setKindname(props.kind);
  req.setKey(deserializeKey(props.key));
  return await client.svc.metaGet(req, client.meta);
};

const getDefaultPartitionId = async (
  props: any
): Promise<GetDefaultPartitionIdResponse> => {
  const client = createGrpcPromiseClient(ConfigstoreMetaServicePromiseClient);
  const req = new GetDefaultPartitionIdRequest();
  return await client.svc.getDefaultPartitionId(req, client.meta);
};

function isPendingDelete(
  pendingTransaction: PendingTransaction,
  entity: MetaEntity
) {
  for (const operation of pendingTransaction.operations) {
    if (operation.hasDeleterequest()) {
      const deleteRequest = g(operation.getDeleterequest());
      if (
        serializeKey(g(deleteRequest.getKey())) ==
        serializeKey(g(entity.getKey()))
      ) {
        return true;
      }
    }
  }
  return false;
}

function getPendingUpdate(pendingTransaction: PendingTransaction, key: Key) {
  let idx = 0;
  for (const operation of pendingTransaction.operations) {
    if (operation.hasUpdaterequest()) {
      const updateRequest = g(operation.getUpdaterequest());
      if (
        serializeKey(g(g(updateRequest.getEntity()).getKey())) ==
        serializeKey(key)
      ) {
        return {
          id: `${idx}`,
          entity: MetaEntity.deserializeBinary(
            g(updateRequest.getEntity()).serializeBinary()
          ),
          operation: operation
        };
      }
    }
    idx++;
  }
  return null;
}

function getPendingCreate(
  pendingTransaction: PendingTransaction,
  idxStr: string | undefined | null
) {
  if (idxStr === "" || idxStr === undefined || idxStr === null) {
    return null;
  }
  const targetIdx = parseInt(idxStr);
  let idx = 0;
  for (const operation of pendingTransaction.operations) {
    if (operation.hasCreaterequest()) {
      const createRequest = g(operation.getCreaterequest());
      if (idx === targetIdx) {
        return {
          id: `${idx}`,
          entity: MetaEntity.deserializeBinary(
            g(createRequest.getEntity()).serializeBinary()
          ),
          operation: operation
        };
      }
    }
    idx++;
  }
  return null;
}

function getConditionalField<T>(
  entity: { value: MetaEntity },
  field: SchemaField,
  def: T,
  grab: (value: Value) => T
): T {
  if (entity === undefined) {
    return def;
  }
  const v = entity.value
    .getValuesList()
    .filter(x => x.getId() === field.getId())[0];
  if (v === undefined) {
    return def;
  }

  return grab(v);
}

function setConditionalField(
  entity: { value: MetaEntity },
  field: SchemaField,
  value: Value
) {
  const valuesList = entity.value.getValuesList();
  for (let i = 0; i < valuesList.length; i++) {
    const existingValue = valuesList[i];
    if (existingValue.getId() === field.getId()) {
      valuesList.splice(i, 1);
      break;
    }
  }
  value.setId(field.getId());
  value.setType(field.getType());
  valuesList.push(value);
  entity.value.setValuesList(valuesList);
}

export const KindEditRoute = (props: KindEditRouteProps) => (
  <PendingTransactionContext.Consumer>
    {value => <KindEditRealRoute {...props} pendingTransaction={value} />}
  </PendingTransactionContext.Consumer>
);

const KindEditRealRoute = (
  props: KindEditRouteProps & { pendingTransaction: PendingTransaction }
) => {
  const isCreate = props.location.pathname.startsWith(
    `/kind/${props.match.params.kind}/create`
  );

  const kindSchema = g(props.schema.getSchema())
    .getKindsMap()
    .get(props.match.params.kind);
  if (kindSchema === undefined) {
    return <>No such kind.</>;
  }

  let useDefaultValue = false;
  if (isCreate) {
    const pendingCreate = getPendingCreate(
      props.pendingTransaction,
      props.match.params.idx
    );
    useDefaultValue = pendingCreate === null;
  }

  const [editableValue, setEditableValue] = useState<
    { value: MetaEntity } | undefined
  >(useDefaultValue ? { value: new MetaEntity() } : undefined);
  const [isSaving] = useState<boolean>(false);
  const [saveError] = useState<any | undefined>(undefined);

  let errorDisplay = null;
  if (saveError !== undefined) {
    errorDisplay = (
      <div className="alert alert-danger" role="alert">
        {JSON.stringify(saveError)}
      </div>
    );
  }

  let isValid = true;

  let header = `Create Entity: ${props.match.params.kind}`;
  let defaultPartitionIdNamespace = "";
  if (isCreate) {
    const pendingCreate = getPendingCreate(
      props.pendingTransaction,
      props.match.params.idx
    );
    if (pendingCreate !== null) {
      if (editableValue === undefined) {
        setEditableValue({ value: pendingCreate.entity });
      }
    } else {
      const response = useAsync<GetDefaultPartitionIdResponse>({
        promiseFn: getDefaultPartitionId
      } as any);
      if (response.isLoading) {
        return (
          <>
            <div className="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-4 border-bottom">
              <h1 className="h2">{header}</h1>
              <div className="btn-toolbar mb-2 mb-md-0" />
            </div>
            {errorDisplay}
            <FontAwesomeIcon icon={faSpinner} spin /> Loading data...
          </>
        );
      } else if (response.error !== undefined) {
        return (
          <>
            <div className="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-4 border-bottom">
              <h1 className="h2">{header}</h1>
              <div className="btn-toolbar mb-2 mb-md-0" />
            </div>
            {errorDisplay}
            {JSON.stringify(response.error)}
          </>
        );
      } else if (response.data !== undefined) {
        defaultPartitionIdNamespace = response.data.getNamespace();
      }
    }
  } else {
    header = `Edit Entity: ${prettifyKey(
      deserializeKey(props.match.params.id)
    )}`;
    const pendingUpdate = getPendingUpdate(
      props.pendingTransaction,
      deserializeKey(props.match.params.id)
    );
    if (pendingUpdate === null) {
      const response = useAsync<MetaGetEntityResponse>({
        promiseFn: getKind,
        watch: props.match.params.kind,
        kind: props.match.params.kind,
        key: props.match.params.id
      } as any);
      if (response.isLoading) {
        return (
          <>
            <div className="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-4 border-bottom">
              <h1 className="h2">{header}</h1>
              <div className="btn-toolbar mb-2 mb-md-0" />
            </div>
            {errorDisplay}
            <FontAwesomeIcon icon={faSpinner} spin /> Loading data...
          </>
        );
      } else if (response.error !== undefined) {
        return (
          <>
            <div className="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-4 border-bottom">
              <h1 className="h2">{header}</h1>
              <div className="btn-toolbar mb-2 mb-md-0" />
            </div>
            {errorDisplay}
            {JSON.stringify(response.error)}
          </>
        );
      } else if (response.data !== undefined) {
        if (editableValue === undefined) {
          const getEntity = response.data.getEntity();
          if (getEntity !== undefined) {
            setEditableValue({ value: getEntity });
          } else {
            setEditableValue({ value: new MetaEntity() });
          }
        }
      }
    } else {
      if (editableValue === undefined) {
        setEditableValue({ value: pendingUpdate.entity });
      }
    }
  }

  if (editableValue === undefined) {
    return (
      <>
        <div className="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-4 border-bottom">
          <h1 className="h2">{header}</h1>
          <div className="btn-toolbar mb-2 mb-md-0" />
        </div>
        {errorDisplay}
        <FontAwesomeIcon icon={faSpinner} spin /> Loading data...
      </>
    );
  }

  const onSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    if (!isValid) {
      // prevent submit if form not valid
      return;
    }

    // Automatically fill in default values.
    for (const field of kindSchema.getFieldsList()) {
      const editor = c(field.getEditor(), new SchemaFieldEditorInfo());
      switch (field.getType()) {
        case ValueType.DOUBLE:
        case ValueType.INT64:
        case ValueType.UINT64:
          {
            const value = getConditionalField(
              editableValue,
              field,
              undefined,
              value => {
                switch (field.getType()) {
                  case ValueType.DOUBLE:
                    return value.getDoublevalue();
                  case ValueType.INT64:
                    return value.getInt64value();
                  case ValueType.UINT64:
                    return value.getUint64value();
                  default:
                    return undefined;
                }
              }
            );
            for (const validator of editor.getValidatorsList()) {
              if (validator.hasDefault()) {
                const def = g(validator.getDefault());
                const defValue = g(def.getValue());
                switch (defValue.getType()) {
                  case ValueType.DOUBLE:
                  case ValueType.INT64:
                  case ValueType.UINT64:
                    if (value === undefined) {
                      setConditionalField(editableValue, field, defValue);
                    }
                    break;
                }
              }
            }
          }
          break;
        case ValueType.STRING:
          const value = getConditionalField(
            editableValue,
            field,
            undefined,
            value => {
              switch (field.getType()) {
                case ValueType.STRING:
                  return value.getStringvalue();
                default:
                  return undefined;
              }
            }
          );
          for (const validator of editor.getValidatorsList()) {
            if (validator.hasDefault()) {
              const def = g(validator.getDefault());
              const defValue = g(def.getValue());
              switch (defValue.getType()) {
                case ValueType.STRING:
                  if (value === undefined) {
                    setConditionalField(editableValue, field, defValue);
                  }
                  break;
              }
            }
          }
          break;
      }
    }

    if (!isCreate) {
      const pendingUpdate = getPendingUpdate(
        props.pendingTransaction,
        deserializeKey(props.match.params.id)
      );
      if (pendingUpdate !== null) {
        const req = new MetaUpdateEntityRequest();
        req.setEntity(editableValue.value);
        pendingUpdate.operation.setUpdaterequest(req);
        props.pendingTransaction.setOperations([
          ...props.pendingTransaction.operations
        ]);
        props.history.push(`/kind/${props.match.params.kind}`);
        return;
      }
    } else {
      const pendingCreate = getPendingCreate(
        props.pendingTransaction,
        props.match.params.idx
      );
      if (pendingCreate !== null) {
        const req = new MetaCreateEntityRequest();
        req.setKindname(props.match.params.kind);
        req.setEntity(editableValue.value);
        pendingCreate.operation.setCreaterequest(req);
        props.pendingTransaction.setOperations([
          ...props.pendingTransaction.operations
        ]);
        props.history.push(`/kind/${props.match.params.kind}`);
        return;
      }
    }

    const operation = new MetaOperation();
    if (isCreate) {
      const req = new MetaCreateEntityRequest();
      req.setKindname(props.match.params.kind);
      req.setEntity(editableValue.value);
      operation.setCreaterequest(req);
    } else {
      const req = new MetaUpdateEntityRequest();
      req.setEntity(editableValue.value);
      operation.setUpdaterequest(req);
    }

    props.pendingTransaction.setOperations([
      ...props.pendingTransaction.operations,
      operation
    ]);
    props.history.push(`/kind/${props.match.params.kind}`);
  };

  let pendingDeleteNotice = null;
  const hasPendingDelete = isPendingDelete(
    props.pendingTransaction,
    editableValue.value
  );
  if (hasPendingDelete) {
    const discardPendingDelete = (e: React.MouseEvent<HTMLAnchorElement>) => {
      e.preventDefault();
      props.pendingTransaction.setOperations([
        ...props.pendingTransaction.operations.filter(
          x =>
            !x.hasDeleterequest() ||
            serializeKey(g(g(x.getDeleterequest()).getKey())) !==
              serializeKey(g(editableValue.value.getKey()))
        )
      ]);
    };
    pendingDeleteNotice = (
      <div className="card border-danger mb-3">
        <div className="card-body text-danger">
          <h5 className="card-title">Pending delete</h5>
          <p className="card-text">
            This entity is pending deletion. To edit it instead, remove the
            deletion operation from the pending changes.
          </p>
          <a onClick={discardPendingDelete} href="#" className="btn btn-dark">
            Discard pending delete
          </a>
        </div>
      </div>
    );
  }

  let keyInput = null;
  if (!isCreate) {
    keyInput = (
      <input
        className="form-control"
        value={prettifyKey(g(editableValue.value.getKey()))}
        readOnly={true}
      />
    );
  } else {
    const key = editableValue.value.getKey();
    let inputValue = "";
    if (key !== undefined) {
      const lastComponent = key.getPathList()[key.getPathList().length - 1];
      inputValue = lastComponent.getName();
    }

    keyInput = (
      <input
        className="form-control"
        value={inputValue}
        placeholder="Auto-generated if left empty."
        onChange={e => {
          const newInputValue = e.target.value;
          const key = new Key();
          const partitionId = new PartitionId();
          partitionId.setNamespace(defaultPartitionIdNamespace);
          key.setPartitionid(partitionId);
          const pathElement = new PathElement();
          pathElement.setKind(props.match.params.kind);
          pathElement.setName(newInputValue);
          key.setPathList([pathElement]);
          editableValue.value.setKey(key);
          setEditableValue({ value: editableValue.value });
        }}
      />
    );
  }

  return (
    <>
      <div className="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-4 border-bottom">
        <h1 className="h2">{header}</h1>
        <div className="btn-toolbar mb-2 mb-md-0" />
      </div>
      {pendingDeleteNotice}
      {errorDisplay}
      <form onSubmit={onSubmit}>
        <div className="form-group">
          <label>ID</label>
          {keyInput}
          <small className="form-text text-muted">
            The ID of the {props.match.params.kind}.{" "}
            {isCreate
              ? "This value can not be changed after the entity is created."
              : ""}
          </small>
        </div>
        {kindSchema.getFieldsList().map(field => {
          const editor = c(field.getEditor(), new SchemaFieldEditorInfo());
          const displayName = c(editor.getDisplayname(), field.getName());
          switch (field.getType()) {
            case ValueType.DOUBLE:
            case ValueType.INT64:
            case ValueType.UINT64: {
              const value = getConditionalField(
                editableValue,
                field,
                undefined,
                value => {
                  switch (field.getType()) {
                    case ValueType.DOUBLE:
                      return value.getDoublevalue();
                    case ValueType.INT64:
                      return value.getInt64value();
                    case ValueType.UINT64:
                      return value.getUint64value();
                    default:
                      return undefined;
                  }
                }
              );
              let fieldErrors = [];
              let placeholder = "";
              for (const validator of editor.getValidatorsList()) {
                if (validator.hasRequired()) {
                  if (value === 0 || value === undefined || isNaN(value)) {
                    fieldErrors.push("A non-zero value is required.");
                  }
                } else if (validator.hasFixedlength()) {
                  // Not applied.
                } else if (validator.hasDefault()) {
                  // Applied at save.
                  const def = g(validator.getDefault());
                  const defValue = g(def.getValue());
                  switch (defValue.getType()) {
                    case ValueType.DOUBLE:
                      placeholder = g(defValue.getDoublevalue()).toString();
                      break;
                    case ValueType.INT64:
                      placeholder = g(defValue.getInt64value()).toString();
                      break;
                    case ValueType.UINT64:
                      placeholder = g(defValue.getUint64value()).toString();
                      break;
                  }
                } else if (validator.hasFormatipaddress()) {
                  // Not applied.
                } else if (validator.hasFormatipaddressport()) {
                  // Not applied.
                }
              }
              if (fieldErrors.length > 0) {
                isValid = false;
              }
              return (
                <div className="form-group" key={field.getId()}>
                  <label>{displayName}</label>
                  <input
                    className={`form-control ${
                      fieldErrors.length > 0 ? "is-invalid" : "is-valid"
                    }`}
                    type="number"
                    value={value}
                    placeholder={placeholder}
                    readOnly={
                      field.getReadonly() || isSaving || hasPendingDelete
                    }
                    onChange={e => {
                      if (editableValue !== undefined) {
                        const value = new Value();
                        switch (field.getType()) {
                          case ValueType.DOUBLE:
                            value.setDoublevalue(parseFloat(e.target.value));
                            break;
                          case ValueType.INT64:
                            value.setInt64value(parseInt(e.target.value));
                            break;
                          case ValueType.UINT64:
                            value.setUint64value(parseInt(e.target.value));
                            break;
                        }
                        setConditionalField(editableValue, field, value);
                        setEditableValue({ value: editableValue.value });
                      }
                    }}
                  />
                  {fieldErrors.map((err, idx) => (
                    <div className="invalid-feedback" key={idx}>
                      {err}
                    </div>
                  ))}
                  <small className="form-text text-muted">
                    {field.getComment()}
                  </small>
                </div>
              );
            }
            case ValueType.STRING: {
              const value = getConditionalField(
                editableValue,
                field,
                "",
                value => value.getStringvalue()
              );
              let fieldErrors = [];
              let placeholder = "";
              for (const validator of editor.getValidatorsList()) {
                if (validator.hasRequired()) {
                  if (value === undefined || value.trim() === "") {
                    fieldErrors.push("A non-empty value is required.");
                  }
                } else if (validator.hasFixedlength()) {
                  if (
                    value.length !== g(validator.getFixedlength()).getLength()
                  ) {
                    fieldErrors.push(
                      `Must be exactly ${g(
                        validator.getFixedlength()
                      ).getLength()} characters long.`
                    );
                  }
                } else if (validator.hasDefault()) {
                  // Applied at save.
                  const def = g(validator.getDefault());
                  const defValue = g(def.getValue());
                  switch (defValue.getType()) {
                    case ValueType.STRING:
                      placeholder = g(defValue.getStringvalue());
                      break;
                  }
                } else if (validator.hasFormatipaddress()) {
                  if (value !== undefined && value !== "") {
                    const address4 = new Address4(value);
                    const address6 = new Address6(value);
                    if (!address4.isValid() && !address6.isValid()) {
                      fieldErrors.push("Must be an IPv4 or IPv6 address.");
                    }
                  }
                } else if (validator.hasFormatipaddressport()) {
                  if (value !== undefined && value !== "") {
                    const portSplit = value.lastIndexOf(":");
                    if (portSplit === -1) {
                      fieldErrors.push(
                        "Must be an IPv4 or IPv6 address, with a port number specified. Use brackets [] around an IPv6 address."
                      );
                    } else {
                      let address = value.substr(0, portSplit);
                      const portNumber = parseInt(value.substr(portSplit + 1));
                      if (
                        address.length > 0 &&
                        address[0] === "[" &&
                        address[address.length - 1] === "]"
                      ) {
                        address = address.substr(1, address.length - 2);
                      }
                      const address4 = new Address4(address);
                      const address6 = new Address6(address);
                      if (!address4.isValid() && !address6.isValid()) {
                        fieldErrors.push(
                          "Must be an IPv4 or IPv6 address, with a port number specified."
                        );
                      }
                      if (isNaN(portNumber)) {
                        fieldErrors.push(
                          "Port number could not be parsed as an integer."
                        );
                      }
                      if (portNumber <= 0 || portNumber >= 65536) {
                        fieldErrors.push(
                          "Port number must be between 1 and 65535."
                        );
                      }
                    }
                  }
                }
              }
              if (fieldErrors.length > 0) {
                isValid = false;
              }
              return (
                <div className="form-group" key={field.getId()}>
                  <label>{displayName}</label>
                  <input
                    className={`form-control ${
                      fieldErrors.length > 0 ? "is-invalid" : "is-valid"
                    }`}
                    value={value}
                    placeholder={placeholder}
                    readOnly={
                      field.getReadonly() || isSaving || hasPendingDelete
                    }
                    onChange={e => {
                      if (editableValue !== undefined) {
                        const value = new Value();
                        value.setStringvalue(e.target.value);
                        setConditionalField(editableValue, field, value);
                        setEditableValue({ value: editableValue.value });
                      }
                    }}
                  />
                  {fieldErrors.map((err, idx) => (
                    <div className="invalid-feedback" key={idx}>
                      {err}
                    </div>
                  ))}
                  <small className="form-text text-muted">
                    {field.getComment()}
                  </small>
                </div>
              );
            }
            case ValueType.BOOLEAN:
              return (
                <div className="form-check" key={field.getId()}>
                  <input
                    id={"checkbox_" + field.getId()}
                    type="checkbox"
                    className="form-check-input"
                    checked={getConditionalField(
                      editableValue,
                      field,
                      false,
                      value => value.getBooleanvalue()
                    )}
                    readOnly={
                      field.getReadonly() || isSaving || hasPendingDelete
                    }
                    onChange={e => {
                      if (editableValue !== undefined) {
                        const value = new Value();
                        value.setBooleanvalue(e.target.checked);
                        setConditionalField(editableValue, field, value);
                        setEditableValue({ value: editableValue.value });
                      }
                    }}
                  />
                  <label
                    htmlFor={"checkbox_" + field.getId()}
                    className="form-check-label"
                  >
                    {displayName}
                  </label>
                  <small className="form-text text-muted">
                    {field.getComment()}
                  </small>
                </div>
              );
            case ValueType.TIMESTAMP:
              const timestampValue = getConditionalField<Timestamp | undefined>(
                editableValue,
                field,
                undefined,
                value => value.getTimestampvalue()
              );
              if (field.getReadonly() || isSaving || hasPendingDelete) {
                return (
                  <div className="form-group" key={field.getId()}>
                    <label>{displayName}</label>
                    <div className="input-group mb-3">
                      <div className="input-group-prepend">
                        <div className="input-group-text">
                          <input
                            type="checkbox"
                            checked={timestampValue !== undefined}
                            disabled={true}
                            onChange={e => {}}
                          />
                        </div>
                      </div>
                      <input
                        className="form-control"
                        defaultValue={
                          timestampValue === undefined
                            ? ""
                            : moment
                                .unix(timestampValue.getSeconds())
                                .toLocaleString()
                        }
                        readOnly={true}
                      />
                    </div>
                    <small className="form-text text-muted">
                      {field.getComment()}
                    </small>
                  </div>
                );
              } else if (timestampValue === undefined) {
                return (
                  <div className="form-group" key={field.getId()}>
                    <label>{displayName}</label>
                    <div className="input-group mb-3">
                      <div className="input-group-prepend">
                        <div className="input-group-text">
                          <input
                            type="checkbox"
                            checked={false}
                            disabled={
                              field.getReadonly() ||
                              isSaving ||
                              hasPendingDelete
                            }
                            onChange={e => {
                              if (editableValue !== undefined) {
                                const timestamp = new Timestamp();
                                timestamp.setSeconds(moment().unix());
                                timestamp.setNanos(0);
                                const value = new Value();
                                value.setTimestampvalue(timestamp);
                                setConditionalField(
                                  editableValue,
                                  field,
                                  value
                                );
                                setEditableValue({
                                  value: editableValue.value
                                });
                              }
                            }}
                          />
                        </div>
                      </div>
                      <input
                        className="form-control"
                        defaultValue={""}
                        readOnly={true}
                      />
                    </div>
                    <small className="form-text text-muted">
                      {field.getComment()}
                    </small>
                  </div>
                );
              } else {
                return (
                  <div className="form-group" key={field.getId()}>
                    <label>{displayName}</label>
                    <div className="input-group mb-3">
                      <div className="input-group-prepend">
                        <div className="input-group-text">
                          <input
                            type="checkbox"
                            checked={true}
                            disabled={
                              field.getReadonly() ||
                              isSaving ||
                              hasPendingDelete
                            }
                            onChange={e => {
                              if (editableValue !== undefined) {
                                const value = new Value();
                                value.clearTimestampvalue();
                                setConditionalField(
                                  editableValue,
                                  field,
                                  value
                                );
                                setEditableValue({
                                  value: editableValue.value
                                });
                              }
                            }}
                          />
                        </div>
                      </div>
                      <Datetime
                        value={moment
                          .unix(timestampValue.getSeconds())
                          .toDate()}
                        onChange={date => {
                          if (editableValue !== undefined) {
                            const timestamp = new Timestamp();
                            timestamp.setSeconds(moment(date).unix());
                            timestamp.setNanos(0);
                            const value = new Value();
                            value.setTimestampvalue(timestamp);
                            setConditionalField(editableValue, field, value);
                            setEditableValue({ value: editableValue.value });
                          }
                        }}
                      />
                    </div>
                    <small className="form-text text-muted">
                      {field.getComment()}
                    </small>
                  </div>
                );
              }
            case ValueType.KEY:
              const keyValue = getConditionalField<Key | undefined>(
                editableValue,
                field,
                undefined,
                value => value.getKeyvalue()
              );
              let fieldErrors = [];
              if (fieldErrors.length === 0) {
                for (const validator of editor.getValidatorsList()) {
                  if (validator.hasRequired()) {
                    if (keyValue === undefined) {
                      fieldErrors.push("You must select a key for this field.");
                    }
                  }
                }
              }
              if (fieldErrors.length > 0) {
                isValid = false;
              }
              if (field.getReadonly() || isSaving || hasPendingDelete) {
                return (
                  <div className="form-group" key={field.getId()}>
                    <label>{displayName}</label>
                    <input
                      className={`form-control ${
                        fieldErrors.length > 0 ? "is-invalid" : "is-valid"
                      }`}
                      defaultValue={
                        keyValue === undefined ? "" : prettifyKey(keyValue)
                      }
                      readOnly={true}
                    />
                    {fieldErrors.map((err, idx) => (
                      <div className="invalid-feedback" key={idx}>
                        {err}
                      </div>
                    ))}
                    <small className="form-text text-muted">
                      {field.getComment()}
                    </small>
                  </div>
                );
              } else {
                return (
                  <div className="form-group key-select" key={field.getId()}>
                    <label>{displayName}</label>
                    <KeySelect
                      className={`${
                        fieldErrors.length > 0 ? "is-invalid" : "is-valid"
                      }`}
                      value={keyValue}
                      onChange={v => {
                        if (editableValue !== undefined) {
                          if (v === undefined) {
                            const value = new Value();
                            value.clearKeyvalue();
                            setConditionalField(editableValue, field, value);
                            setEditableValue({ value: editableValue.value });
                          } else {
                            const value = new Value();
                            value.setKeyvalue(v);
                            setConditionalField(editableValue, field, value);
                            setEditableValue({ value: editableValue.value });
                          }
                        }
                      }}
                      schema={g(props.schema.getSchema())}
                      field={field}
                    />
                    {fieldErrors.map((err, idx) => (
                      <div className="invalid-feedback" key={idx}>
                        {err}
                      </div>
                    ))}
                    <small className="form-text text-muted">
                      {field.getComment()}
                    </small>
                  </div>
                );
              }
            case ValueType.BYTES: {
              const value = getConditionalField(
                editableValue,
                field,
                "",
                value => value.getBytesvalue_asB64()
              );
              let fieldErrors = [];
              try {
                atob(value);
              } catch (err) {
                fieldErrors.push("A valid base64-encoded value is required");
              }
              if (fieldErrors.length === 0) {
                for (const validator of editor.getValidatorsList()) {
                  if (validator.hasRequired()) {
                    const b = atob(value);
                    if (b === "") {
                      fieldErrors.push("A non-empty value is required.");
                    }
                  } else if (validator.hasFixedlength()) {
                    const b = atob(value);
                    if (
                      b.length !== g(validator.getFixedlength()).getLength()
                    ) {
                      fieldErrors.push(
                        `Must be exactly ${g(
                          validator.getFixedlength()
                        ).getLength()} bytes long.`
                      );
                    }
                  } else if (validator.hasDefault()) {
                    // Not applied.
                  } else if (validator.hasFormatipaddress()) {
                    // Not applied.
                  } else if (validator.hasFormatipaddressport()) {
                    // Not applied.
                  }
                }
              }
              if (fieldErrors.length > 0) {
                isValid = false;
              }
              return (
                <div className="form-group" key={field.getId()}>
                  <label>{displayName}</label>
                  <input
                    className={`form-control ${
                      fieldErrors.length > 0 ? "is-invalid" : "is-valid"
                    }`}
                    value={value}
                    placeholder="Copy and paste a base64-encoded value into this field."
                    readOnly={
                      field.getReadonly() || isSaving || hasPendingDelete
                    }
                    onChange={e => {
                      if (editableValue !== undefined) {
                        const value = new Value();
                        let base64Value = "";
                        try {
                          base64Value = atob(e.target.value);
                        } catch {
                          return;
                        }
                        value.setBytesvalue(
                          Uint8Array.from(base64Value, c => c.charCodeAt(0))
                        );
                        setConditionalField(editableValue, field, value);
                        setEditableValue({ value: editableValue.value });
                      }
                    }}
                  />
                  {fieldErrors.map((err, idx) => (
                    <div className="invalid-feedback" key={idx}>
                      {err}
                    </div>
                  ))}
                  <small className="form-text text-muted">
                    {field.getComment()}
                  </small>
                </div>
              );
            }
            default:
              return (
                <div className="form-group" key={field.getId()}>
                  <label>{displayName}</label>
                  <input
                    className="form-control"
                    type="text"
                    value={`This field has type ${field
                      .getType()
                      .toString()} which is not an understood type in the UI.`}
                    readOnly={true}
                  />
                  <small className="form-text text-muted">
                    {field.getComment()}
                  </small>
                </div>
              );
          }
        })}
        <button
          type="submit"
          className="btn btn-primary mb-4"
          disabled={isSaving || hasPendingDelete || !isValid}
        >
          {isSaving ? (
            <>
              <FontAwesomeIcon icon={faSpinner} spin />
              &nbsp;
            </>
          ) : (
            ""
          )}
          Queue Changes
        </button>
        <Link
          className={
            "btn btn-outline-secondary ml-2 mb-4" +
            (hasPendingDelete ? " disabled" : "")
          }
          to={hasPendingDelete ? "#" : `/kind/${props.match.params.kind}`}
        >
          Discard Changes
        </Link>
      </form>
    </>
  );
};
