import React, { useState, useEffect } from "react";
import { KindRouteProps } from "./KindRoute";
import { RouteComponentProps } from "react-router";
import {
  GetSchemaResponse,
  MetaListEntitiesResponse,
  MetaListEntitiesRequest
} from "../api/meta_pb";
import { g, serializeKey } from "../core";
import { grpc } from "@improbable-eng/grpc-web";
import { ConfigstoreMetaService } from "../api/meta_pb_service";
import { UnaryOutput } from "@improbable-eng/grpc-web/dist/typings/unary";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faSpinner, faPencilAlt } from "@fortawesome/free-solid-svg-icons";
import { Link } from "react-router-dom";

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

export const KindListRoute = (props: KindListRouteProps) => {
  const [data, setData] = useState<MetaListEntitiesResponse | null>(null);
  const [selected, setSelected] = useState<SetHolder>({ v: new Set<string>() });

  useEffect(() => {
    setData(null);
    selected.v.clear();
    setSelected({ v: selected.v });
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

  const kindSchema = g(props.schema.getSchema())
    .getKindsList()
    .filter(kind => kind.getName() == props.match.params.kind)[0];

  const kindDisplay =
    selected.v.size == 1
      ? g(kindSchema.getEditor()).getSingular()
      : g(kindSchema.getEditor()).getPlural();

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
                    data !== null
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
                    if (data !== null) {
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
