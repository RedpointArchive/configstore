import React, { Component, useState, useEffect } from "react";
import { GetSchemaResponse } from "./api/meta_pb";
import { GetSchemaRequest } from "./api/meta_pb";
import { ConfigstoreMetaServicePromiseClient } from "./api/meta_grpc_web_pb";
import { Switch, Route, RouteComponentProps } from "react-router";
import { KindRoute, KindRouteMatch } from "./routes/KindRoute";
import { BrowserRouter, Link, NavLink } from "react-router-dom";
import { g } from "./core";
import { useAsync } from "react-async";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faSpinner } from "@fortawesome/free-solid-svg-icons";
import { grpcHost } from "./svcHost"

const loadSchema = async () => {
  const svc = new ConfigstoreMetaServicePromiseClient(
    grpcHost,
    null,
    null
  );
  const req = new GetSchemaRequest();
  return await svc.getSchema(req, {});
}

const App = () => {
  const { data, error, isLoading } = useAsync(loadSchema);
  
  let kindsList = null;
  let content: React.ReactNode = null;
  if (isLoading) {
    kindsList = (
      <li className="nav-item">
      <span className="nav-link text-info">
        <FontAwesomeIcon icon={faSpinner} fixedWidth spin /> Waiting for schema...
        </span>
      </li>
    )
  } else if (error) {
    kindsList = (
      <li className="nav-item text-danger">
      <span className="nav-link">
      {JSON.stringify(error)}
      </span>
      </li>
    )
  } else if (data) {
    kindsList = g(data.getSchema()).getKindsMap().getEntryList()
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
      ))
      content = (
        <Switch>
          <Route path="/kind/:kind">
            {(props: RouteComponentProps<KindRouteMatch>) => (
              <KindRoute {...props} schema={data} />
            )}
          </Route>
        </Switch>
      );
  }

  return (
    <BrowserRouter>
      <>
        <nav className="navbar navbar-dark fixed-top bg-dark flex-md-nowrap p-0 shadow">
          <a className="navbar-brand col-sm-3 col-md-2 mr-0" href="#">
            configstore
          </a>
        </nav>

        <div className="container-fluid">
          <div className="row">
            <nav className="col-md-2 d-none d-md-block bg-light sidebar">
              <div className="sidebar-sticky">
      <>
        <ul className="nav flex-column">
          <li className="nav-item">
            <NavLink className="nav-link" activeClassName="active" to="/" exact>
              Dashboard
            </NavLink>
          </li>
        </ul>

        <h6 className="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted">
          <span>Kinds</span>
        </h6>
        <ul className="nav flex-column">
          {kindsList}
        </ul>

        <h6 className="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted">
          <span>SDKs</span>
        </h6>
        <ul className="nav flex-column mb-2">
          <li className="nav-item">
            <a className="nav-link" href="/sdk/client.proto" target="_blank">
              gRPC Protocol Spec
            </a>
          </li>
          <li className="nav-item">
            <a className="nav-link" href="/sdk/client.go" target="_blank">
              gRPC Go Client
            </a>
          </li>
        </ul>
      </></div>
            </nav>

            <main role="main" className="col-sm-12 col-md-10 ml-sm-auto px-4">
              {content}
            </main>
          </div>
        </div>
      </>
    </BrowserRouter>
  );
};

export default App;
