import React, { useState /* useEffect */ } from "react";
import { RouteComponentProps } from "react-router";
import {
  GetSchemaResponse,
  MetaListEntitiesResponse,
  ValueType
} from "../api/meta_pb";
import { g } from "../core";
import { Link } from "react-router-dom";
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

export const KindEditRoute = (props: KindEditRouteProps) => {
  const [data, setData] = useState<MetaListEntitiesResponse | null>(null);

  const isCreate = props.location.pathname.startsWith(
    `/kind/${props.match.params.kind}/create`
  );
  const createVerb = isCreate ? "Create" : "Edit";

  /*
  useEffect(() => {
    setData(null);
    const req = new MetaListEntitiesRequest();
    req.setKindname(props.match.params.kind);
    req.setStart("");
    req.setLimit(0);
    grpc.unary(ConfigstoreMetaService.MetaList, {
      request: req,
      host: "http://localhost:13390",
      onEnd: (res: UnaryOutput<MetaListEntitiesResponse>) => {
        const { status, statusMessage, headers, message, trailers } = res;
        if (status === grpc.Code.OK && message) {
          setData(message);
        }
      }
    });
  }, [props.match.params.kind]);
  */

 const kindSchema = g(props.schema.getSchema())
 .getKindsMap().get(props.match.params.kind);
 if (kindSchema === undefined) {
   return (<>No such kind.</>);
 }

  /*
  const kindDisplay =
    selected.v.size == 1
      ? g(kindSchema.getEditor()).getSingular()
      : g(kindSchema.getEditor()).getPlural();
  */

  /*
  let dataset = [
    <tr key="loading">
      <td
        colSpan={3 + kindSchema.getFieldsList().length}
        style={{
          textAlign: "center"
        }}
      >
        <FontAwesomeIcon icon={faSpinner} spin /> Loading data...
      </td>
    </tr>
  ];
  if (data !== null && data.getEntitiesList().length == 0) {
    dataset = [
      <tr key="loading">
        <td
          colSpan={3 + kindSchema.getFieldsList().length}
          style={{
            textAlign: "center"
          }}
        >
          There are no entities of this kind.
        </td>
      </tr>
    ];
  } else if (data !== null) {
    dataset = [];
    for (const entity of data.getEntitiesList()) {
      dataset.push(
        <tr key={entity.getId()}>
          <td className="w-checkbox">
            <input
              type="checkbox"
              checked={selected.v.has(entity.getId())}
              onChange={e => {
                if (e.target.checked) {
                  selected.v.add(entity.getId());
                } else {
                  selected.v.delete(entity.getId());
                }
                setSelected({ v: selected.v });
              }}
            />
          </td>
          <td>
            <Link
              to={`/kind/${props.match.params.kind}/edit/${entity.getId()}`}
            >
              {entity.getId()}
            </Link>
          </td>
          {kindSchema.getFieldsList().map(field => {
            const fieldData = entity
              .getValuesList()
              .filter(fieldData => fieldData.getId() == field.getId())[0];
            if (fieldData == undefined) {
              return (
                <td key={field.getId()}>
                  <em className="text-muted">-</em>
                </td>
              );
            }
            return <td key={field.getId()}>{fieldData.getStringvalue()}</td>;
          })}
          <td className="w-checkbox">
            <Link
              to={`/kind/${props.match.params.kind}/edit/${entity.getId()}`}
            >
              <FontAwesomeIcon icon={faPencilAlt} />
            </Link>
          </td>
        </tr>
      );
    }
  }
  */

  return (
    <>
      <div className="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-4 border-bottom">
        <h1 className="h2">
          {createVerb} Kind: {props.match.params.kind}
        </h1>
        <div className="btn-toolbar mb-2 mb-md-0" />
      </div>
      <form>
        <div className="form-group">
          <label>ID</label>
          <input className="form-control" value={""} readOnly={true} />
          <small className="form-text text-muted">
            The ID of the {props.match.params.kind}
          </small>
        </div>
        {kindSchema.getFieldsList().map(field => {
          switch (field.getType()) {
            case ValueType.DOUBLE:
            case ValueType.INT64:
            case ValueType.UINT64:
              return (
                <div className="form-group">
                  <label>{g(field.getEditor()).getDisplayname()}</label>
                  <input
                    className="form-control"
                    type="number"
                    value={0}
                    readOnly={field.getReadonly()}
                  />
                  <small className="form-text text-muted">
                    {field.getComment()}
                  </small>
                </div>
              );
            case ValueType.STRING:
              return (
                <div className="form-group">
                  <label>{g(field.getEditor()).getDisplayname()}</label>
                  <input
                    className="form-control"
                    value={""}
                    readOnly={field.getReadonly()}
                  />
                  <small className="form-text text-muted">
                    {field.getComment()}
                  </small>
                </div>
              );
            case ValueType.BOOLEAN:
              return (
                <div className="form-check">
                  <input
                    type="checkbox"
                    className="form-check-input"
                    readOnly={field.getReadonly()}
                  />
                  <label className="form-check-label">
                    {g(field.getEditor()).getDisplayname()}
                  </label>
                  <small className="form-text text-muted">
                    {field.getComment()}
                  </small>
                </div>
              );
            case ValueType.TIMESTAMP:
              return (
                <div className="form-group">
                  <label>{g(field.getEditor()).getDisplayname()}</label>
                  <input
                    className="form-control"
                    type="datetime-local"
                    value={""}
                    readOnly={field.getReadonly()}
                  />
                  <small className="form-text text-muted">
                    {field.getComment()}
                  </small>
                </div>
              );
          }
        })}
        <button type="submit" className="btn btn-primary">
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
