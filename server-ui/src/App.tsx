import React, { useState } from "react";
import { MetaOperation, MetaTransactionResult } from "./api/meta_pb";
import { GetSchemaRequest } from "./api/meta_pb";
import { ConfigstoreMetaServicePromiseClient } from "./api/meta_grpc_web_pb";
import { Switch, Route, RouteComponentProps } from "react-router";
import { KindRoute, KindRouteMatch } from "./routes/KindRoute";
import { BrowserRouter, NavLink } from "react-router-dom";
import { g } from "./core";
import { useAsync } from "react-async";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faSpinner } from "@fortawesome/free-solid-svg-icons";
import { DashboardRoute } from "./routes/DashboardRoute";
import { SaveRoute } from "./routes/SaveRoute";
import { encode, decode } from "base64-arraybuffer";
import { ReviewRoute } from "./routes/ReviewRoute";
import { createGrpcPromiseClient } from "./svcHost";
import { ApiHostContext, ApiHostConfig } from "./ApiHostProvider";

interface PendingTransactionInternal {
  operations: MetaOperation[];
}

export interface PendingTransaction {
  operations: MetaOperation[];
  setOperations(operations: MetaOperation[]): void;
  responseOriginalOperations: MetaOperation[];
  setResponseOriginalOperations(operations: MetaOperation[]): void;
  response: MetaTransactionResult | null;
  setResponse(response: MetaTransactionResult | null): void;
}

export const PendingTransactionContext = React.createContext<
  PendingTransaction
>({
  operations: [],
  setOperations: () => {},
  responseOriginalOperations: [],
  setResponseOriginalOperations: () => {},
  response: null,
  setResponse: () => {}
});

const loadSchema = async () => {
  const client = createGrpcPromiseClient(ConfigstoreMetaServicePromiseClient);
  const req = new GetSchemaRequest();
  return await client.svc.getSchema(req, client.meta);
};

function loadOperationsFromLocalStorage() {
  const v = window.localStorage.getItem("pendingOperations");
  if (v === null) {
    return [];
  }
  const m = JSON.parse(v) as string[];
  return m.map(x => MetaOperation.deserializeBinary(new Uint8Array(decode(x))));
}

function saveOperationsToLocalStorage(ops: MetaOperation[]) {
  const m = ops.map(x => {
    return encode(x.serializeBinary());
  });
  window.localStorage.setItem("pendingOperations", JSON.stringify(m));
}

const App = () => {
  const { data, error, isLoading } = useAsync(loadSchema);
  const [transaction, setTransaction] = useState<PendingTransactionInternal>({
    operations: loadOperationsFromLocalStorage()
  });
  const [responseOriginalOperations, setResponseOriginalOperations] = useState<
    PendingTransactionInternal
  >({
    operations: []
  });
  const [response, setResponse] = useState<MetaTransactionResult | null>(null);
  const pendingTransaction = {
    response: response,
    setResponse: setResponse,
    responseOriginalOperations: responseOriginalOperations.operations,
    setResponseOriginalOperations: (value: MetaOperation[]) => {
      setResponseOriginalOperations({
        operations: value
      });
    },
    operations: transaction.operations,
    setOperations: (value: MetaOperation[]) => {
      saveOperationsToLocalStorage(value);
      setTransaction({
        operations: value
      });
    }
  };

  let kindsList: React.ReactNode = null;
  let content: React.ReactNode = null;
  if (isLoading) {
    kindsList = (
      <li className="nav-item">
        <span className="nav-link text-info">
          <FontAwesomeIcon icon={faSpinner} fixedWidth spin /> Waiting for
          schema...
        </span>
      </li>
    );
  } else if (error) {
    kindsList = (
      <li className="nav-item text-danger">
        <span className="nav-link">{JSON.stringify(error)}</span>
      </li>
    );
  } else if (data) {
    kindsList = g(data.getSchema())
      .getKindsMap()
      .getEntryList()
      .map(kind => (
        <li className="nav-item" key={kind[0]}>
          <NavLink
            className="nav-link"
            activeClassName="active"
            to={`/kind/${kind[0]}`}
          >
            {kind[0]}
          </NavLink>
        </li>
      ));
    content = (
      <Switch>
        <Route path="/" exact component={DashboardRoute} />
        <Route path="/save" exact>
          {(props: RouteComponentProps<KindRouteMatch>) => (
            <SaveRoute {...props} schema={data} />
          )}
        </Route>
        <Route path="/review" exact>
          {(props: RouteComponentProps<KindRouteMatch>) => (
            <ReviewRoute {...props} schema={data} />
          )}
        </Route>
        <Route path="/kind/:kind">
          {(props: RouteComponentProps<KindRouteMatch>) => (
            <KindRoute {...props} schema={data} />
          )}
        </Route>
      </Switch>
    );
  }

  const saveLink =
    transaction.operations.length === 0 ? (
      <NavLink className="btn btn-success btn-block" to="/save">
        No pending changes.
      </NavLink>
    ) : (
      <NavLink className="btn btn-warning btn-block" to="/save">
        You have {pendingTransaction.operations.length} pending change
        {pendingTransaction.operations.length === 1 ? "" : "s"}. Click here to
        save.
      </NavLink>
    );

  const getSvcHost = () =>
    (JSON.parse(window.localStorage.getItem(
      "apiHostConfig"
    ) as string) as ApiHostConfig).svcHost;

  return (
    <PendingTransactionContext.Provider value={pendingTransaction}>
      <ApiHostContext.Consumer>
        {value => (
          <>
            {g(value).hideTitlebar ? null : (
              <nav className="navbar navbar-dark fixed-top bg-dark flex-md-nowrap p-0 shadow">
                <a className="navbar-brand col-sm-3 col-md-2 mr-0" href="#">
                  configstore
                </a>
              </nav>
            )}
            <div
              className={
                "container-fluid" +
                (g(value).hideTitlebar ? " no-titlebar" : "")
              }
            >
              <div className="row">
                <nav className="col-md-2 d-none d-md-block bg-light sidebar">
                  <div className="sidebar-sticky">
                    <>
                      <ul className="nav flex-column">
                        <li className="nav-item">
                          <NavLink
                            className="nav-link"
                            activeClassName="active"
                            to="/"
                            exact
                          >
                            Dashboard
                          </NavLink>
                        </li>
                        <li className="nav-item">
                          <span className="nav-link">{saveLink}</span>
                        </li>
                      </ul>

                      <h6 className="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted">
                        <span>Kinds</span>
                      </h6>
                      <ul className="nav flex-column">{kindsList}</ul>

                      <h6 className="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted">
                        <span>SDKs</span>
                      </h6>
                      <ul className="nav flex-column mb-2">
                        <li className="nav-item">
                          <a
                            className="nav-link"
                            href={`${getSvcHost()}/sdk/client.proto`}
                            target="_blank"
                          >
                            gRPC Protocol Spec
                          </a>
                        </li>
                        <li className="nav-item">
                          <a
                            className="nav-link"
                            href={`${getSvcHost()}/sdk/client.go`}
                            target="_blank"
                          >
                            gRPC Go Client
                          </a>
                        </li>
                      </ul>
                    </>
                  </div>
                </nav>

                <main
                  role="main"
                  className="col-sm-12 col-md-10 ml-sm-auto px-4"
                >
                  {content}
                </main>
              </div>
            </div>
          </>
        )}
      </ApiHostContext.Consumer>
    </PendingTransactionContext.Provider>
  );
};

export default App;
