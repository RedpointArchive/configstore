import { RouteComponentProps } from "react-router";
import { MetaTransactionBatch, WatchTransactionsRequest } from "../api/meta_pb";
import { useState, useEffect } from "react";
import { ConfigstoreMetaServicePromiseClient } from "../api/meta_grpc_web_pb";
import { grpcHost } from "../svcHost";
import { g, serializeKey, prettifyKey, getLastKindOfKey } from "../core";
import * as React from "react";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faCircle, faStopCircle } from "@fortawesome/free-solid-svg-icons";
import { Error } from "grpc-web";
import { Link } from "react-router-dom";

export interface DashboardRouteProps extends RouteComponentProps<{}> {}

interface DashboardRouteState {
  transactions: MetaTransactionBatch[];
  gotInitialState: boolean;
  connected: boolean;
  error: Error | null;
}

export class DashboardRoute extends React.Component<
  DashboardRouteProps,
  DashboardRouteState
> {
  private onUnmount: (() => void) | null = null;

  constructor(props: DashboardRouteProps) {
    super(props);

    this.state = {
      transactions: [],
      gotInitialState: false,
      connected: false,
      error: null
    };
  }

  private connect = () => {
    const svc = new ConfigstoreMetaServicePromiseClient(grpcHost, null, null);
    const req = new WatchTransactionsRequest();
    const stream = svc.watchTransactions(req, {});
    this.setState({
      connected: true,
      error: null
    });
    stream.on("error", err => {
      this.setState({
        error: err,
        connected: false
      });
      stream.cancel();
      // try to reconnect
      this.connect();
    });
    stream.on("data", resp => {
      if (resp.hasInitialstate()) {
        this.setState({
          gotInitialState: true
        });
      } else if (resp.hasBatch()) {
        let trs = [g(resp.getBatch()), ...this.state.transactions];
        if (trs.length > 50) {
          trs = trs.slice(0, 50);
        }
        this.setState({
          transactions: trs
        });
      }
    });
    stream.on("end", () => {
      this.setState({
        connected: false
      });
      // try to reconnect
      this.connect();
    });
    this.onUnmount = () => stream.cancel();
  };

  public componentDidMount() {
    this.connect();
  }

  public componentWillUnmount() {
    if (this.onUnmount !== null) {
      this.onUnmount();
      this.onUnmount = null;
    }
  }

  public render() {
    return (
      <>
        <div className="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-0 border-bottom">
          <h1 className="h2">Dashboard: Recent Transactions</h1>
          <div
            className="btn-toolbar mb-2 mb-md-0"
            style={{ lineHeight: "100%" }}
          >
            {this.state.connected ? (
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
                <th>Entities Mutated</th>
                <th>Entities Deleted</th>
              </tr>
            </thead>
            <tbody>
              {this.state.transactions.length === 0 ? (
                <>
                  <tr>
                    <td colSpan={4} className="text-muted">
                      Transactions will be shown here as they arrive. Historical
                      transactions are not shown.
                    </td>
                  </tr>
                </>
              ) : null}
              {this.state.error !== null ? (
                <>
                  <tr>
                    <td colSpan={4} className="text-danger">
                      {JSON.stringify(this.state.error)}
                    </td>
                  </tr>
                </>
              ) : null}
              {this.state.transactions.map(transaction => (
                <tr key={transaction.getId()}>
                  <td>{transaction.getId()}</td>
                  <td>{transaction.getDescription()}</td>
                  <td>
                    {transaction.getMutatedentitiesList().map(entity => (
                      <Link
                        key={serializeKey(g(entity.getKey()))}
                        style={{
                          display: "block"
                        }}
                        to={`/kind/${getLastKindOfKey(
                          g(entity.getKey())
                        )}/edit/${serializeKey(g(entity.getKey()))}`}
                      >
                        {prettifyKey(g(entity.getKey()))}
                      </Link>
                    ))}
                  </td>
                  <td>
                    {transaction.getDeletedkeysList().map(key => (
                      <div key={serializeKey(key)}>{prettifyKey(key)}</div>
                    ))}
                  </td>
                </tr>
              ))}
              {this.state.gotInitialState ? (
                <>
                  <tr>
                    <td colSpan={4} className="text-muted">
                      Successfully received the initial database state.
                    </td>
                  </tr>
                </>
              ) : null}
            </tbody>
          </table>
        </div>
      </>
    );
  }
}
