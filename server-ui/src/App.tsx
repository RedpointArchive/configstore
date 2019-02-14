import React, { Component, useState, useEffect } from "react";
import { GetSchemaResponse } from "./api/meta_pb";
import { GetSchemaRequest } from "./api/meta_pb";
import { grpc } from "@improbable-eng/grpc-web";
import { ConfigstoreMetaService } from "./api/meta_pb_service";
import { UnaryOutput } from "@improbable-eng/grpc-web/dist/typings/unary";
import { Switch, Route, RouteComponentProps } from "react-router";
import { KindRoute, KindRouteMatch } from "./routes/KindRoute";
import { BrowserRouter, Link, NavLink } from "react-router-dom";
import { g } from "./core";

const App = () => {
  const [schema, setSchema] = useState<GetSchemaResponse | null>(null);

  useEffect(() => {
    if (schema === null) {
      const req = new GetSchemaRequest();
      grpc.unary(ConfigstoreMetaService.GetSchema, {
        request: req,
        host: "http://localhost:13390",
        onEnd: (res: UnaryOutput<GetSchemaResponse>) => {
          const { status, statusMessage, headers, message, trailers } = res;
          if (status === grpc.Code.OK && message) {
            setSchema(message);
          }
        }
      });
    }
  });

  let nav = null;
  let content = null;
  if (schema != null) {
    nav = (
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
          {g(schema.getSchema())
            .getKindsList()
            .map(kind => (
              <li className="nav-item">
                <NavLink
                  className="nav-link"
                  activeClassName="active"
                  to={`/kind/${kind.getName()}`}
                >
                  {kind.getName()}
                </NavLink>
              </li>
            ))}
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
      </>
    );
    content = (
      <Switch>
        <Route path="/kind/:kind">
          {(props: RouteComponentProps<KindRouteMatch>) => (
            <KindRoute {...props} schema={schema} />
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
              <div className="sidebar-sticky">{nav}</div>
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
