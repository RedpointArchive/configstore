import React from "react";
import { RouteComponentProps, Redirect } from "react-router";
import { GetSchemaResponse, MetaOperation, Key, Schema } from "../api/meta_pb";
import { PendingTransactionContext, PendingTransaction } from "../App";
import { g, serializeKey, getLastKindOfKey, prettifyKey } from "../core";
import { Link } from "react-router-dom";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faCheck } from "@fortawesome/free-solid-svg-icons";
import { KeyView } from "../KeyView";

export interface ReviewRouteProps extends RouteComponentProps<{}> {
  schema: GetSchemaResponse;
}

function getTypeForOperation(operation: MetaOperation) {
  if (operation.hasCreaterequest()) {
    return "Create";
  }
  if (operation.hasDeleterequest()) {
    return "Delete";
  }
  if (operation.hasGetrequest()) {
    return "Get";
  }
  if (operation.hasListrequest()) {
    return "List";
  }
  if (operation.hasUpdaterequest()) {
    return "Update";
  }
  return "(Unknown)";
}

function getEntityLinkForOperation(
  operation: MetaOperation,
  pendingTransaction: PendingTransaction,
  schema: Schema
) {
  let key: Key | null = null;
  if (operation.hasCreaterequest()) {
    return g(operation.getCreaterequest()).getKindname();
  }
  if (operation.hasDeleterequest()) {
    key = g(g(operation.getDeleterequest()).getKey());
  }
  if (operation.hasGetrequest()) {
    return null;
  }
  if (operation.hasListrequest()) {
    return null;
  }
  if (operation.hasUpdaterequest()) {
    key = g(g(g(operation.getUpdaterequest()).getEntity()).getKey());
  }
  if (key !== null) {
    return (
      <KeyView
        pendingTransaction={pendingTransaction}
        schema={schema}
        value={key}
      />
    );
  }
  return null;
}

function displayResult(
  pendingTransaction: PendingTransaction,
  idx: number,
  schema: Schema
) {
  if (pendingTransaction.response === null) {
    return null;
  }
  const results = pendingTransaction.response.getOperationresultsList();
  if (idx >= results.length) {
    return null;
  }
  const result = results[idx];
  if (result.hasError()) {
    return (
      <span className="text-danger">
        {g(result.getError()).getErrormessage()}
      </span>
    );
  } else {
    if (result.hasCreateresponse()) {
      const key = g(g(g(result.getCreateresponse()).getEntity()).getKey());
      return (
        <span>
          <FontAwesomeIcon icon={faCheck} fixedWidth />{" "}
          <KeyView
            pendingTransaction={pendingTransaction}
            schema={schema}
            value={key}
          />
        </span>
      );
    } else {
      return <FontAwesomeIcon icon={faCheck} fixedWidth />;
    }
  }
}

export const ReviewRoute = (props: ReviewRouteProps) => (
  <PendingTransactionContext.Consumer>
    {value => <ReviewRealRoute {...props} pendingTransaction={value} />}
  </PendingTransactionContext.Consumer>
);

const ReviewRealRoute = (
  props: ReviewRouteProps & { pendingTransaction: PendingTransaction }
) => {
  if (props.pendingTransaction.response === null) {
    return <Redirect to="/save" />;
  }

  let messageArea = (
    <div className="alert alert-success mt-4" role="alert">
      Everything was successfully saved.
    </div>
  );
  if (
    g(props.pendingTransaction.response)
      .getOperationresultsList()
      .filter(x => x.hasError()).length > 0
  ) {
    messageArea = (
      <div className="alert alert-warning mt-4" role="alert">
        One or more operations did not complete successfully. Review the output
        below for the error message.
      </div>
    );
  }

  return (
    <>
      <div className="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-0 border-bottom">
        <h1 className="h2">Review Save Result</h1>
      </div>
      {messageArea}
      <div className="table-responsive table-fixed-header">
        <table className="table table-sm table-bt-none table-hover">
          <thead>
            <tr>
              <th>Idx</th>
              <th>Type</th>
              <th>Entity</th>
              <th>Result</th>
            </tr>
          </thead>
          <tbody>
            {props.pendingTransaction.responseOriginalOperations.map(
              (value, idx) => (
                <tr key={idx}>
                  <td>{idx}</td>
                  <td>{getTypeForOperation(value)}</td>
                  <td>
                    {getEntityLinkForOperation(
                      value,
                      props.pendingTransaction,
                      g(props.schema.getSchema())
                    )}
                  </td>
                  <td>
                    {displayResult(
                      props.pendingTransaction,
                      idx,
                      g(props.schema.getSchema())
                    )}
                  </td>
                </tr>
              )
            )}
          </tbody>
        </table>
      </div>
    </>
  );
};
