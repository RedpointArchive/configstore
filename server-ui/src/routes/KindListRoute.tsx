import React, { useState, useEffect } from "react";
import { KindRouteProps } from "./KindRoute";
import { RouteComponentProps } from "react-router";
import {
  GetSchemaResponse,
  MetaListEntitiesResponse,
  MetaListEntitiesRequest
} from "../api/meta_pb";
import { g, serializeKey } from "../core";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faSpinner, faPencilAlt } from "@fortawesome/free-solid-svg-icons";
import { Link } from "react-router-dom";
import { ConfigstoreMetaServicePromiseClient } from "../api/meta_grpc_web_pb";
import { grpcHost } from "../svcHost"
import { useAsync } from "react-async";

export interface KindListRouteMatch {
  kind: string;
}

export interface KindListRouteProps
  extends RouteComponentProps<KindListRouteMatch> {
  schema: GetSchemaResponse;
}

interface SetHolder {
  v: Set<string>;
}

const listKinds = async (props: any) => {
  const svc = new ConfigstoreMetaServicePromiseClient(
    grpcHost,
    null,
    null
  );
  const req = new MetaListEntitiesRequest();
  req.setKindname(props.kind); 
  req.setStart("");
  req.setLimit(0);
  return await svc.metaList(req, {});
}

export const KindListRoute = (props: KindListRouteProps) => {
  const { data, error, isLoading } = useAsync<MetaListEntitiesResponse>({
    promiseFn: listKinds,
    watch: props.match.params.kind, 
    kind: props.match.params.kind,
  } as any);
  const [selected, setSelected] = useState<SetHolder>({ v: new Set<string>() });

  const kindSchema = g(props.schema.getSchema())
    .getKindsMap().get(props.match.params.kind);
    if (kindSchema === undefined) {
      return (<>No such kind.</>);
    }

  const kindDisplay =
    selected.v.size == 1
      ? g(kindSchema.getEditor()).getSingular()
      : g(kindSchema.getEditor()).getPlural();

  let dataset: React.ReactNode[] = [];
  if (isLoading) {
    dataset = [
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
    ]
  }
  else if (error) {
    dataset = [
      <tr key="loading">
        <td
          colSpan={3 + kindSchema.getFieldsList().length}
          style={{
            textAlign: "center"
          }}
        >
          {JSON.stringify(error)}
        </td>
      </tr>
    ];
  }
  else if (data !== undefined && data.getEntitiesList().length == 0) {
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
  } else if (data !== undefined) {
    dataset = [];
    for (const entity of data.getEntitiesList()) {
      dataset.push(
        <tr key={serializeKey(g(entity.getKey()))}>
          <td className="w-checkbox">
            <input
              type="checkbox"
              checked={selected.v.has(serializeKey(g(entity.getKey())))}
              onChange={e => {
                if (e.target.checked) {
                  selected.v.add(serializeKey(g(entity.getKey())));
                } else {
                  selected.v.delete(serializeKey(g(entity.getKey())));
                }
                setSelected({ v: selected.v });
              }}
            />
          </td>
          <td>
            <Link
              to={`/kind/${props.match.params.kind}/edit/${serializeKey(
                g(entity.getKey())
              )}`}
            >
              {serializeKey(g(entity.getKey()))}
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
              to={`/kind/${props.match.params.kind}/edit/${serializeKey(
                g(entity.getKey())
              )}`}
            >
              <FontAwesomeIcon icon={faPencilAlt} />
            </Link>
          </td>
        </tr>
      );
    }
  }

  return (
    <>
      <div className="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-0 border-bottom">
        <h1 className="h2">Kind: {props.match.params.kind}</h1>
        <div className="btn-toolbar mb-2 mb-md-0">
          <button
            type="button"
            className={
              "btn btn-sm mr-2 " +
              (selected.v.size == 0 ? "btn-outline-danger" : "btn-danger")
            }
            disabled={selected.v.size == 0}
          >
            Delete {selected.v.size} {kindDisplay}
          </button>
          <Link
            to={`/kind/${props.match.params.kind}/create`}
            className="btn btn-sm btn-success"
          >
            Create {props.match.params.kind}
          </Link>
        </div>
      </div>
      <div className="table-responsive table-fixed-header">
        <table className="table table-sm table-bt-none table-hover">
          <thead>
            <tr>
              <th className="w-checkbox">
                <input
                  type="checkbox"
                  checked={
                    data !== undefined
                      ? data
                          .getEntitiesList()
                          .filter(
                            value =>
                              !selected.v.has(serializeKey(g(value.getKey())))
                          ).length === 0
                        ? true
                        : false
                      : false
                  }
                  onChange={e => {
                    if (data !== undefined) {
                      if (e.target.checked) {
                        selected.v.clear();
                        for (const entity of data.getEntitiesList()) {
                          selected.v.add(serializeKey(g(entity.getKey())));
                        }
                      } else {
                        selected.v.clear();
                      }
                      setSelected({ v: selected.v });
                    }
                  }}
                />
              </th>
              <th>ID</th>
              {kindSchema.getFieldsList().map(field => (
                <th key={field.getId()}>
                  {g(field.getEditor()).getDisplayname()}
                </th>
              ))}
              <th className="w-checkbox">
                <FontAwesomeIcon icon={faPencilAlt} />
              </th>
            </tr>
          </thead>
          <tbody>{dataset}</tbody>
        </table>
      </div>
    </>
  );
};
