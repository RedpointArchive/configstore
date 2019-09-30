import {
  Key,
  Schema,
  MetaEntity,
  MetaGetEntityRequest,
  MetaGetEntityResponse,
  SchemaKind,
  ValueType
} from "./api/meta_pb";
import { Link } from "react-router-dom";
import { getLastKindOfKey, serializeKey, g, prettifyKey } from "./core";
import * as React from "react";
import { PendingTransaction } from "./App";
import { createGrpcPromiseClient } from "./svcHost";
import { ConfigstoreMetaServicePromiseClient } from "./api/meta_grpc_web_pb";
import { useAsync } from "react-async";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faSpinner } from "@fortawesome/free-solid-svg-icons";

function getPendingUpdate(
  pendingTransaction: PendingTransaction,
  entityKey: Key
) {
  let idx = 0;
  for (const operation of pendingTransaction.operations) {
    if (operation.hasUpdaterequest()) {
      const updateRequest = g(operation.getUpdaterequest());
      if (
        serializeKey(g(g(updateRequest.getEntity()).getKey())) ==
        serializeKey(entityKey)
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

const loadEntity = async (props: {
  value: Key;
}): Promise<MetaGetEntityResponse> => {
  const client = createGrpcPromiseClient(ConfigstoreMetaServicePromiseClient);
  const request = new MetaGetEntityRequest();
  request.setKey(props.value);
  request.setKindname(getLastKindOfKey(props.value));
  const response = await client.svc.metaGet(request, client.meta);
  return response;
};

const EntityValueRenderer = (props: {
  entity: MetaEntity;
  kind: SchemaKind;
}) => {
  const kindForLookupEditor = props.kind.getEditor();
  if (kindForLookupEditor !== undefined) {
    const values = props.entity.getValuesList();
    const fields = props.kind
      .getFieldsList()
      .filter(
        x =>
          x.getName() === kindForLookupEditor.getRendereditordropdownwithfield()
      );
    const lookupField = fields.length === 0 ? undefined : fields[0];
    if (lookupField !== undefined) {
      const value = values.filter(x => x.getId() === lookupField.getId())[0];
      switch (lookupField.getType()) {
        case ValueType.STRING:
          return <>{value.getStringvalue()}</>;
      }
    }
  }
  return <>{prettifyKey(g(props.entity.getKey()))}</>;
};

const KeyAsyncView = (props: { value: Key; kind: SchemaKind }) => {
  const { data, error, isLoading } = useAsync<MetaGetEntityResponse>({
    promiseFn: loadEntity,
    value: props.value
  } as any);
  if (isLoading || data === undefined) {
    return <FontAwesomeIcon icon={faSpinner} spin fixedWidth />;
  }

  let content: React.ReactNode = prettifyKey(g(props.value));

  const r = data.getEntity();
  if (r !== undefined) {
    content = <EntityValueRenderer entity={r} kind={props.kind} />;
  }

  return (
    <Link
      to={`/kind/${getLastKindOfKey(props.value)}/edit/${serializeKey(
        g(props.value)
      )}`}
    >
      {content}
    </Link>
  );
};

export const KeyView = (props: {
  pendingTransaction: PendingTransaction;
  schema: Schema;
  value: Key | undefined;
}) => {
  if (props.value === undefined) {
    return <>-</>;
  }

  const defaultStyle = (
    <Link
      to={`/kind/${getLastKindOfKey(props.value)}/edit/${serializeKey(
        g(props.value)
      )}`}
    >
      {prettifyKey(props.value)}
    </Link>
  );

  const kindsMap = props.schema.getKindsMap();
  const kind = kindsMap.get(getLastKindOfKey(props.value));
  if (kind === undefined) {
    console.log("no kind: " + getLastKindOfKey(props.value));
    return defaultStyle;
  }

  const kindEditor = kind.getEditor();
  if (kindEditor === undefined) {
    console.log("no editor override");
    return defaultStyle;
  }

  if (kindEditor.getRendereditordropdownwithfield() === "") {
    console.log("using default rendering for that type");
    return defaultStyle;
  }

  // Otherwise, use the dropdown field to render the link.
  let content: React.ReactNode;
  const pendingUpdate = getPendingUpdate(props.pendingTransaction, props.value);
  if (pendingUpdate !== null) {
    content = <EntityValueRenderer entity={pendingUpdate.entity} kind={kind} />;
  } else {
    // We need to look it up.
    content = <KeyAsyncView value={props.value} kind={kind} />;
  }

  return (
    <Link
      to={`/kind/${getLastKindOfKey(props.value)}/edit/${serializeKey(
        g(props.value)
      )}`}
    >
      {content}
    </Link>
  );
};
