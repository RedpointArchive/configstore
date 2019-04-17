import { RouteComponentProps } from "react-router";
import { MetaTransactionBatch, WatchTransactionsRequest } from "../api/meta_pb";
import { useState, useEffect } from "react";
import { ConfigstoreMetaServicePromiseClient } from "../api/meta_grpc_web_pb";
import { grpcHost } from "../svcHost";
import { g } from "../core";
import * as React from "react";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faCircle, faStopCircle } from "@fortawesome/free-solid-svg-icons";
import { Error } from "grpc-web";

export interface DashboardRouteProps extends RouteComponentProps<{}> {}

export const DashboardRoute = (props: DashboardRouteProps) => {
  const [transactions, setTransactions] = useState<{
    transactions: MetaTransactionBatch[];
  }>({ transactions: [] });
  const [connected, setConnected] = useState<boolean>(false);
  const [error, setError] = useState<Error | null>(null);
  useEffect(() => {
    const svc = new ConfigstoreMetaServicePromiseClient(grpcHost, null, null);
    const req = new WatchTransactionsRequest();
    const stream = svc.watchTransactions(req, {});
    setConnected(true);
    stream.on("error", err => {
      setError(err);
      setConnected(false);
      stream.cancel();
    });
    stream.on("data", resp => {
      console.log(resp.getBatch());
      transactions.transactions.splice(0, 0, g(resp.getBatch()));
      setTransactions({
        transactions: transactions.transactions
      });
    });
    stream.on("end", () => {
      setConnected(false);
    });
    return () => stream.cancel();
  }, []);

  return (
    <>
      <div className="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-0 border-bottom">
        <h1 className="h2">Dashboard: Recent Transactions</h1>
        <div
          className="btn-toolbar mb-2 mb-md-0"
          style={{ lineHeight: "100%" }}
        >
          {connected ? (
            <span className="text-success">
              <FontAwesomeIcon icon={faCircle} /> Connected
            </span>
          ) : (
            <span className="text-danger">
              <FontAwesomeIcon icon={faStopCircle} /> Disconnected
            </span>
          )}
        </div>
      </div>
      <div className="table-responsive table-fixed-header">
        <table className="table table-sm table-bt-none table-hover">
          <thead>
            <tr>
              <th>ID</th>
              <th>Description</th>
            </tr>
          </thead>
          <tbody>
            {transactions.transactions.length === 0 ? (
              <>
                <tr>
                  <td colSpan={2} className="text-muted">
                    Transactions will be shown here as they arrive. Historical
                    transactions are not shown.
                  </td>
                </tr>
              </>
            ) : null}
            {error !== null ? (
              <>
                <tr>
                  <td colSpan={2} className="text-danger">
                    {JSON.stringify(error)}
                  </td>
                </tr>
              </>
            ) : null}
            {transactions.transactions.map(transaction => (
              <tr key={transaction.getId()}>
                <td>{transaction.getId()}</td>
                <td>{transaction.getDescription()}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </>
  );
};
