import React, { useState, useEffect } from "react";
import { KindRouteProps } from "./KindRoute";
import { RouteComponentProps } from "react-router";
import {
  GetSchemaResponse,
  MetaListEntitiesResponse,
  MetaListEntitiesRequest,
  MetaDeleteEntityRequest,
  ValueType,
  SchemaFieldEditorInfo,
  MetaOperation,
  MetaEntity
} from "../api/meta_pb";
import {
  g,
  serializeKey,
  deserializeKey,
  prettifyKey,
  getLastKindOfKey,
  c
} from "../core";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faSpinner,
  faPencilAlt,
  faCheck
} from "@fortawesome/free-solid-svg-icons";
import { Link } from "react-router-dom";
import { ConfigstoreMetaServicePromiseClient } from "../api/meta_grpc_web_pb";
import { grpcHost } from "../svcHost";
import { useAsync } from "react-async";
import { PendingTransactionContext, PendingTransaction } from "../App";

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
  const svc = new ConfigstoreMetaServicePromiseClient(grpcHost, null, null);
  const req = new MetaListEntitiesRequest();
  req.setKindname(props.kind);
  req.setStart("");
  req.setLimit(0);
  return await svc.metaList(req, {});
};

export const KindListRoute = (props: KindListRouteProps) => (
  <PendingTransactionContext.Consumer>
    {value => <KindListRealRoute {...props} pendingTransaction={value} />}
  </PendingTransactionContext.Consumer>
);

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

const KindListRealRoute = (
  props: KindListRouteProps & { pendingTransaction: PendingTransaction }
) => {
  const [refreshCount, setRefreshCount] = useState<number>(0);
  const { data, error, isLoading } = useAsync<MetaListEntitiesResponse>({
    promiseFn: listKinds,
    watch: refreshCount + "_" + props.match.params.kind,
    kind: props.match.params.kind
  } as any);
  const [selected, setSelected] = useState<SetHolder>({ v: new Set<string>() });

  const kindSchema = g(props.schema.getSchema())
    .getKindsMap()
    .get(props.match.params.kind);
  if (kindSchema === undefined) {
    return <>No such kind.</>;
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
    ];
  } else if (error) {
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
  } else if (data !== undefined && data.getEntitiesList().length == 0) {
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
      const pendingDelete = isPendingDelete(props.pendingTransaction, entity);
      dataset.push(
        <tr
          key={serializeKey(g(entity.getKey()))}
          className={pendingDelete ? "strikethrough" : ""}
        >
          <td className="w-checkbox">
            <input
              type="checkbox"
              checked={
                selected.v.has(serializeKey(g(entity.getKey()))) &&
                !pendingDelete
              }
              disabled={pendingDelete}
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
              {prettifyKey(g(entity.getKey()))}
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
            switch (fieldData.getType()) {
              case ValueType.STRING:
                return (
                  <td key={field.getId()}>{fieldData.getStringvalue()}</td>
                );
              case ValueType.DOUBLE:
                return (
                  <td key={field.getId()}>{fieldData.getDoublevalue()}</td>
                );
              case ValueType.INT64:
                return <td key={field.getId()}>{fieldData.getInt64value()}</td>;
              case ValueType.UINT64:
                return (
                  <td key={field.getId()}>{fieldData.getUint64value()}</td>
                );
              case ValueType.KEY:
                const childKey = fieldData.getKeyvalue();
                if (childKey === undefined) {
                  return <td key={field.getId()}>-</td>;
                } else {
                  return (
                    <td key={field.getId()}>
                      <Link
                        to={`/kind/${getLastKindOfKey(
                          childKey
                        )}/edit/${serializeKey(g(childKey))}`}
                      >
                        {prettifyKey(childKey)}
                      </Link>
                    </td>
                  );
                }
              case ValueType.BOOLEAN:
                return (
                  <td key={field.getId()}>
                    {fieldData.getBooleanvalue() ? (
                      <FontAwesomeIcon icon={faCheck} fixedWidth />
                    ) : (
                      "-"
                    )}
                  </td>
                );
              case ValueType.BYTES:
                return (
                  <td key={field.getId()}>
                    <em>(bytes)</em>
                  </td>
                );
              default:
                return (
                  <td key={field.getId()}>
                    (unknown type {fieldData.getType()})
                  </td>
                );
            }
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

  const doRefresh = () => {
    setRefreshCount(refreshCount + 1);
  };

  const startDelete = async (e: React.MouseEvent<HTMLButtonElement>) => {
    const ops = [];
    const arrayCopy = Array.from(selected.v);
    for (const key of arrayCopy) {
      const operation = new MetaOperation();
      const req = new MetaDeleteEntityRequest();
      req.setKindname(props.match.params.kind);
      req.setKey(deserializeKey(key));
      operation.setDeleterequest(req);
      ops.push(operation);
    }

    props.pendingTransaction.setOperations([
      ...props.pendingTransaction.operations,
      ...ops
    ]);
    setRefreshCount(refreshCount + 1);
  };

  return (
    <>
      <div className="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-0 border-bottom">
        <h1 className="h2">Kind: {props.match.params.kind}</h1>
        <div className="btn-toolbar mb-2 mb-md-0">
          <button
            type="button"
            className={
              "btn btn-sm mr-2 " +
              (selected.v.size == 0 ? "btn-outline-primary" : "btn-primary")
            }
            onClick={doRefresh}
          >
            Refresh {kindDisplay}
          </button>
          <button
            type="button"
            className={
              "btn btn-sm mr-2 " +
              (selected.v.size == 0 ? "btn-outline-danger" : "btn-danger")
            }
            disabled={selected.v.size == 0}
            onClick={startDelete}
          >
            Delete {selected.v.size} {kindDisplay}
          </button>
          <Link
            to={`/kind/${props.match.params.kind}/create`}
            className={"btn btn-sm btn-success"}
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
              {kindSchema.getFieldsList().map(field => {
                const editor = (field.getEditor(), new SchemaFieldEditorInfo());
                const displayName = c(editor.getDisplayname(), field.getName());
                return <th key={field.getId()}>{displayName}</th>;
              })}
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
