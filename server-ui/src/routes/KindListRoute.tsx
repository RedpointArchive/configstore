import React, { useState, useEffect } from "react";
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
import { createGrpcPromiseClient } from "../svcHost";
import { useAsync } from "react-async";
import { PendingTransactionContext, PendingTransaction } from "../App";
import moment from "moment";
import { nibblinsToDollarString } from "../FinancialInput";
import BigInt from "big-integer";
import { KeyView } from "../KeyView";

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
  const client = createGrpcPromiseClient(ConfigstoreMetaServicePromiseClient);
  const req = new MetaListEntitiesRequest();
  req.setKindname(props.kind);
  req.setStart("");
  req.setLimit(0);
  return await client.svc.metaList(req, client.meta);
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

function getPendingUpdate(
  pendingTransaction: PendingTransaction,
  entity: MetaEntity
) {
  let idx = 0;
  for (const operation of pendingTransaction.operations) {
    if (operation.hasUpdaterequest()) {
      const updateRequest = g(operation.getUpdaterequest());
      if (
        serializeKey(g(g(updateRequest.getEntity()).getKey())) ==
        serializeKey(g(entity.getKey()))
      ) {
        return {
          id: `${idx}`,
          entity: g(updateRequest.getEntity()),
          operation: operation
        };
      }
    }
    idx++;
  }
  return null;
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
  } else if (data !== undefined) {
    dataset = [];
    const entities: {
      id: string;
      selectId: string;
      entity: MetaEntity;
      operation: MetaOperation | null;
    }[] = [];
    let transactionIdx = 0;
    for (const operation of props.pendingTransaction.operations) {
      if (operation.hasCreaterequest()) {
        const createRequest = g(operation.getCreaterequest());
        if (createRequest.getKindname() === props.match.params.kind) {
          entities.push({
            id: `${transactionIdx}`,
            selectId: `pendingop_${transactionIdx}`,
            entity: g(createRequest.getEntity()),
            operation: g(operation)
          });
        }
      }
      transactionIdx++;
    }
    entities.push(
      ...data.getEntitiesList().map(x => ({
        id: serializeKey(g(x.getKey())),
        selectId: serializeKey(g(x.getKey())),
        entity: g(x),
        operation: null
      }))
    );
    if (entities.length === 0) {
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
    } else {
      for (const entity of entities) {
        const pendingDelete = isPendingDelete(
          props.pendingTransaction,
          entity.entity
        );
        const pendingUpdate = getPendingUpdate(
          props.pendingTransaction,
          entity.entity
        );
        const effectiveEntity =
          pendingUpdate === null ? entity.entity : pendingUpdate.entity;
        dataset.push(
          <tr key={entity.id} className={pendingDelete ? "strikethrough" : ""}>
            <td className="w-checkbox">
              <input
                type="checkbox"
                checked={selected.v.has(entity.selectId) && !pendingDelete}
                disabled={pendingDelete}
                onChange={e => {
                  if (e.target.checked) {
                    selected.v.add(entity.selectId);
                  } else {
                    selected.v.delete(entity.selectId);
                  }
                  setSelected({ v: selected.v });
                }}
              />
            </td>
            <td>
              {entity.operation !== null ? (
                <Link
                  to={`/kind/${props.match.params.kind}/create/pending/${
                    entity.id
                  }`}
                >
                  Pending{" "}
                  {effectiveEntity.getKey() === undefined
                    ? props.match.params.kind
                    : prettifyKey(g(effectiveEntity.getKey()))}
                </Link>
              ) : (
                <Link
                  to={`/kind/${props.match.params.kind}/edit/${serializeKey(
                    g(effectiveEntity.getKey())
                  )}`}
                >
                  {prettifyKey(g(effectiveEntity.getKey()))}
                </Link>
              )}
            </td>
            {kindSchema.getFieldsList().map(field => {
              const fieldData = effectiveEntity
                .getValuesList()
                .filter(fieldData => fieldData.getId() == field.getId())[0];
              const editor = c(field.getEditor(), new SchemaFieldEditorInfo());
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
                  return (
                    <td key={field.getId()}>
                      {editor.getUsefinancialvaluetonibblinsconversion()
                        ? nibblinsToDollarString(
                            BigInt(fieldData.getInt64value())
                          )
                        : fieldData.getInt64value()}
                    </td>
                  );
                case ValueType.UINT64:
                  return (
                    <td key={field.getId()}>
                      {editor.getUsefinancialvaluetonibblinsconversion()
                        ? nibblinsToDollarString(
                            BigInt(fieldData.getUint64value())
                          )
                        : fieldData.getUint64value()}
                    </td>
                  );
                case ValueType.KEY:
                  const childKey = fieldData.getKeyvalue();
                  return (
                    <td key={field.getId()}>
                      <KeyView
                        pendingTransaction={props.pendingTransaction}
                        schema={g(props.schema.getSchema())}
                        value={childKey}
                      />
                    </td>
                  );
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
                case ValueType.TIMESTAMP:
                  const timestamp = fieldData.getTimestampvalue();
                  if (timestamp === undefined) {
                    return <td key={field.getId()}>-</td>;
                  } else {
                    return (
                      <td key={field.getId()}>
                        {moment.unix(timestamp.getSeconds()).toLocaleString()}
                      </td>
                    );
                  }
                default:
                  return (
                    <td key={field.getId()}>
                      (unknown type {fieldData.getType()})
                    </td>
                  );
              }
            })}
            <td className="w-checkbox">
              {entity.operation !== null ? (
                <Link
                  to={`/kind/${props.match.params.kind}/create/pending/${
                    entity.id
                  }`}
                >
                  <FontAwesomeIcon icon={faPencilAlt} />
                </Link>
              ) : (
                <Link
                  to={`/kind/${props.match.params.kind}/edit/${serializeKey(
                    g(effectiveEntity.getKey())
                  )}`}
                >
                  <FontAwesomeIcon icon={faPencilAlt} />
                </Link>
              )}
            </td>
          </tr>
        );
      }
    }
  }

  const doRefresh = () => {
    setRefreshCount(refreshCount + 1);
  };

  const startDelete = async (e: React.MouseEvent<HTMLButtonElement>) => {
    const ops = [];
    const pendingOpsToRemove = [];
    const oldOps = [...props.pendingTransaction.operations];
    const arrayCopy = Array.from(selected.v);
    for (const key of arrayCopy) {
      if (key.startsWith("pendingop_")) {
        pendingOpsToRemove.push(parseInt(key.substr("pendingop_".length)));
      }
    }
    pendingOpsToRemove.sort((a, b) => b - a);
    console.log(pendingOpsToRemove);
    for (const opId of pendingOpsToRemove) {
      oldOps.splice(opId, 1);
    }
    for (const key of arrayCopy) {
      if (!key.startsWith("pendingop_")) {
        const operation = new MetaOperation();
        const req = new MetaDeleteEntityRequest();
        req.setKindname(props.match.params.kind);
        req.setKey(deserializeKey(key));
        operation.setDeleterequest(req);
        ops.push(operation);
      }
    }
    props.pendingTransaction.setOperations([...oldOps, ...ops]);
    setRefreshCount(refreshCount + 1);
    setSelected({ v: new Set<string>() });
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
