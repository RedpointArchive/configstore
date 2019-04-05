import React, { useState /* useEffect */ } from "react";
import { RouteComponentProps } from "react-router";
import {
  GetSchemaResponse,
  MetaListEntitiesResponse,
  ValueType,
  MetaGetEntityRequest
} from "../api/meta_pb";
import { g, deserializeKey } from "../core";
import { Link } from "react-router-dom";
import { useAsync } from "react-async";
import { ConfigstoreMetaServicePromiseClient } from "../api/meta_grpc_web_pb";
import { grpcHost } from "../svcHost";
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

const getKind = async (props: any) => {
  const svc = new ConfigstoreMetaServicePromiseClient(
    grpcHost,
    null,
    null
  );
  const req = new MetaGetEntityRequest();
  req.setKindname(props.kind); 
  req.setKey(deserializeKey(props.key))
  return await svc.metaGet(req, {});
}

export const KindEditRoute = (props: KindEditRouteProps) => {
  const isCreate = props.location.pathname.startsWith(
    `/kind/${props.match.params.kind}/create`
  );
  const createVerb = isCreate ? "Create" : "Edit";

 const kindSchema = g(props.schema.getSchema())
 .getKindsMap().get(props.match.params.kind);
 if (kindSchema === undefined) {
   return (<>No such kind.</>);
 }

 if (isCreate) {

 } else {
    const { data, error, isLoading } = useAsync<MetaListEntitiesResponse>({
      promiseFn: getKind,
      watch: props.match.params.kind, 
      kind: props.match.params.kind,
      key: props.match.params.id,
    } as any);
  }

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
