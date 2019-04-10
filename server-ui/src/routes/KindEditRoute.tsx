import React, { useState /* useEffect */ } from "react";
import { RouteComponentProps } from "react-router";
import {
  GetSchemaResponse,
  MetaListEntitiesResponse,
  ValueType,
  MetaGetEntityRequest,
  MetaGetEntityResponse,
  SchemaField,
  Value,
  MetaEntity,
  MetaUpdateEntityRequest,
  MetaCreateEntityRequest
} from "../api/meta_pb";
import { g, deserializeKey, serializeKey, prettifyKey } from "../core";
import { Link } from "react-router-dom";
import { useAsync } from "react-async";
import { ConfigstoreMetaServicePromiseClient } from "../api/meta_grpc_web_pb";
import { grpcHost } from "../svcHost";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faSpinner } from "@fortawesome/free-solid-svg-icons";
/*
import { grpc } from "@improbable-eng/grpc-web";
import { ConfigstoreMetaService } from "../api/meta_pb_service";
import { UnaryOutput } from "@improbable-eng/grpc-web/dist/typings/unary";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faSpinner, faPencilAlt } from "@fortawesome/free-solid-svg-icons";
import { Link } from "react-router-dom";
*/

export interface KindEditRouteMatch {
  kind: string;
  id: string;
}

export interface KindEditRouteProps
  extends RouteComponentProps<KindEditRouteMatch> {
  schema: GetSchemaResponse;
}

const getKind = async (props: any): Promise<MetaGetEntityResponse> => {
  const svc = new ConfigstoreMetaServicePromiseClient(grpcHost, null, null);
  const req = new MetaGetEntityRequest();
  req.setKindname(props.kind);
  req.setKey(deserializeKey(props.key));
  return await svc.metaGet(req, {});
};

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

function setConditionalField<T>(
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

export const KindEditRoute = (props: KindEditRouteProps) => {
  const isCreate = props.location.pathname.startsWith(
    `/kind/${props.match.params.kind}/create`
  );

  const kindSchema = g(props.schema.getSchema())
    .getKindsMap()
    .get(props.match.params.kind);
  if (kindSchema === undefined) {
    return <>No such kind.</>;
  }

  const [editableValue, setEditableValue] = useState<
    { value: MetaEntity } | undefined
  >(isCreate ? { value: new MetaEntity() } : undefined);
  const [isSaving, setIsSaving] = useState<boolean>(false);
  const [saveError, setSaveError] = useState<any | undefined>(undefined);

  let errorDisplay = null;
  if (saveError !== undefined) {
    errorDisplay = (
      <div className="alert alert-danger" role="alert">
        {JSON.stringify(saveError)}
      </div>
    );
  }

  let header = `Create Entity: ${props.match.params.kind}`;
  if (!isCreate) {
    header = `Edit Entity: ${prettifyKey(
      deserializeKey(props.match.params.id)
    )}`;
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

    setIsSaving(true);
    setSaveError(undefined);
    try {
      const svc = new ConfigstoreMetaServicePromiseClient(grpcHost, null, null);
      if (isCreate) {
        const req = new MetaCreateEntityRequest();
        req.setKindname(props.match.params.kind);
        req.setEntity(editableValue.value);
        await svc.metaCreate(req, {});
      } else {
        const req = new MetaUpdateEntityRequest();
        req.setEntity(editableValue.value);
        await svc.metaUpdate(req, {});
      }
      props.history.push(`/kind/${props.match.params.kind}`);
    } catch (e) {
      setSaveError(e);
    } finally {
      setIsSaving(false);
    }
  };

  return (
    <>
      <div className="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-4 border-bottom">
        <h1 className="h2">{header}</h1>
        <div className="btn-toolbar mb-2 mb-md-0" />
      </div>
      {errorDisplay}
      <form onSubmit={onSubmit}>
        <div className="form-group">
          <label>ID</label>
          <input
            className="form-control"
            value={
              !isCreate
                ? prettifyKey(g(editableValue.value.getKey()))
                : "(Automatically generated on save)"
            }
            readOnly={true}
          />
          <small className="form-text text-muted">
            The ID of the {props.match.params.kind}
          </small>
        </div>
        {kindSchema.getFieldsList().map(field => {
          switch (field.getType()) {
            case ValueType.DOUBLE:
            case ValueType.INT64:
            case ValueType.UINT64: {
              const value = getConditionalField(
                editableValue,
                field,
                0,
                value => {
                  switch (field.getType()) {
                    case ValueType.DOUBLE:
                      return value.getDoublevalue();
                    case ValueType.INT64:
                      return value.getInt64value();
                    case ValueType.UINT64:
                      return value.getUint64value();
                    default:
                      return 0;
                  }
                }
              );
              return (
                <div className="form-group" key={field.getId()}>
                  <label>{g(field.getEditor()).getDisplayname()}</label>
                  <input
                    className="form-control"
                    type="number"
                    value={value}
                    readOnly={field.getReadonly() || isSaving}
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
              return (
                <div className="form-group" key={field.getId()}>
                  <label>{g(field.getEditor()).getDisplayname()}</label>
                  <input
                    className="form-control"
                    value={value}
                    readOnly={field.getReadonly() || isSaving}
                    onChange={e => {
                      if (editableValue !== undefined) {
                        const value = new Value();
                        value.setStringvalue(e.target.value);
                        setConditionalField(editableValue, field, value);
                        setEditableValue({ value: editableValue.value });
                      }
                    }}
                  />
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
                    readOnly={field.getReadonly() || isSaving}
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
                    {g(field.getEditor()).getDisplayname()}
                  </label>
                  <small className="form-text text-muted">
                    {field.getComment()}
                  </small>
                </div>
              );
            case ValueType.TIMESTAMP:
              return (
                <div className="form-group" key={field.getId()}>
                  <label>{g(field.getEditor()).getDisplayname()}</label>
                  <input
                    className="form-control"
                    type="datetime-local"
                    value={""}
                    readOnly={field.getReadonly() || isSaving}
                  />
                  <small className="form-text text-muted">
                    {field.getComment()}
                  </small>
                </div>
              );
            case ValueType.BYTES: {
              const value = getConditionalField(
                editableValue,
                field,
                "",
                value => value.getBytesvalue_asB64()
              );
              return (
                <div className="form-group" key={field.getId()}>
                  <label>{g(field.getEditor()).getDisplayname()}</label>
                  <input
                    className="form-control"
                    value={value}
                    readOnly={field.getReadonly() || isSaving}
                    onChange={e => {
                      if (editableValue !== undefined) {
                        const value = new Value();
                        value.setBytesvalue(e.target.value);
                        setConditionalField(editableValue, field, value);
                        setEditableValue({ value: editableValue.value });
                      }
                    }}
                  />
                  <small className="form-text text-muted">
                    {field.getComment()}
                  </small>
                </div>
              );
            }
            default:
              return (
                <div className="form-group" key={field.getId()}>
                  <label>{g(field.getEditor()).getDisplayname()}</label>
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
        <button type="submit" className="btn btn-primary" disabled={isSaving}>
          {isSaving ? (
            <>
              <FontAwesomeIcon icon={faSpinner} spin />
              &nbsp;
            </>
          ) : (
            ""
          )}
          Save Changes
        </button>
        <Link
          className="btn btn-outline-secondary ml-2"
          to={`/kind/${props.match.params.kind}`}
        >
          Discard Changes
        </Link>
      </form>
    </>
  );
};
