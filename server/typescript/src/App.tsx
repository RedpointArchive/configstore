import React, { Component, useState, useEffect } from "react";
import { GetSchemaResponse } from "./api/meta_pb";
import { GetSchemaRequest } from "./api/meta_pb";
import { grpc } from "@improbable-eng/grpc-web";
import { ConfigstoreMetaService } from "./api/meta_pb_service";
import { UnaryOutput } from "@improbable-eng/grpc-web/dist/typings/unary";

const App = () => {
  const [data, setData] = useState<GetSchemaResponse | null>(null);

  useEffect(() => {
    if (data === null) {
      const req = new GetSchemaRequest();
      grpc.unary(ConfigstoreMetaService.GetSchema, {
        request: req,
        host: "http://localhost:13390",
        onEnd: (res: UnaryOutput<GetSchemaResponse>) => {
          const { status, statusMessage, headers, message, trailers } = res;
          if (status === grpc.Code.OK && message) {
            setData(message);
          }
        }
      });
    }
  });

  function t<T>(v: T | undefined): T {
    return v as T;
  }

  let nav = null;
  if (data != null) {
    nav = (
      <>
        <ul className="nav flex-column">
          <li className="nav-item">
            <a className="nav-link active" href="#">
              Dashboard
            </a>
          </li>
        </ul>

        <h6 className="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted">
          <span>Kinds</span>
        </h6>
        <ul className="nav flex-column">
          {t(data.getSchema())
            .getKindsList()
            .map(kind => (
              <li className="nav-item">
                <a className="nav-link" href="#">
                  {kind.getName()}
                </a>
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
  }

  return (
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
            <div className="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-0 border-bottom">
              <h1 className="h2">Kind: User</h1>
              <div className="btn-toolbar mb-2 mb-md-0">
                <div className="btn-group mr-2">
                  <button
                    type="button"
                    className="btn btn-sm btn-outline-secondary"
                  >
                    Page 1
                  </button>
                  <button
                    type="button"
                    className="btn btn-sm btn-outline-secondary"
                  >
                    2
                  </button>
                </div>
                <button
                  type="button"
                  className="btn btn-sm btn-success"
                >
                  Create User
                </button>
              </div>
            </div>
            <div className="table-responsive">
              <table className="table table-striped table-sm">
                <thead>
                  <tr>
                    <th style={{borderTop: 'none'}}>#</th>
                    <th style={{borderTop: 'none'}}>Header</th>
                    <th style={{borderTop: 'none'}}>Header</th>
                    <th style={{borderTop: 'none'}}>Header</th>
                    <th style={{borderTop: 'none'}}>Header</th>
                  </tr>
                </thead>
                <tbody>
                  <tr>
                    <td>1,001</td>
                    <td>Lorem</td>
                    <td>ipsum</td>
                    <td>dolor</td>
                    <td>sit</td>
                  </tr>
                  <tr>
                    <td>1,002</td>
                    <td>amet</td>
                    <td>consectetur</td>
                    <td>adipiscing</td>
                    <td>elit</td>
                  </tr>
                  <tr>
                    <td>1,003</td>
                    <td>Integer</td>
                    <td>nec</td>
                    <td>odio</td>
                    <td>Praesent</td>
                  </tr>
                  <tr>
                    <td>1,003</td>
                    <td>libero</td>
                    <td>Sed</td>
                    <td>cursus</td>
                    <td>ante</td>
                  </tr>
                  <tr>
                    <td>1,004</td>
                    <td>dapibus</td>
                    <td>diam</td>
                    <td>Sed</td>
                    <td>nisi</td>
                  </tr>
                  <tr>
                    <td>1,005</td>
                    <td>Nulla</td>
                    <td>quis</td>
                    <td>sem</td>
                    <td>at</td>
                  </tr>
                  <tr>
                    <td>1,006</td>
                    <td>nibh</td>
                    <td>elementum</td>
                    <td>imperdiet</td>
                    <td>Duis</td>
                  </tr>
                  <tr>
                    <td>1,007</td>
                    <td>sagittis</td>
                    <td>ipsum</td>
                    <td>Praesent</td>
                    <td>mauris</td>
                  </tr>
                  <tr>
                    <td>1,008</td>
                    <td>Fusce</td>
                    <td>nec</td>
                    <td>tellus</td>
                    <td>sed</td>
                  </tr>
                  <tr>
                    <td>1,009</td>
                    <td>augue</td>
                    <td>semper</td>
                    <td>porta</td>
                    <td>Mauris</td>
                  </tr>
                  <tr>
                    <td>1,010</td>
                    <td>massa</td>
                    <td>Vestibulum</td>
                    <td>lacinia</td>
                    <td>arcu</td>
                  </tr>
                  <tr>
                    <td>1,011</td>
                    <td>eget</td>
                    <td>nulla</td>
                    <td>Class</td>
                    <td>aptent</td>
                  </tr>
                  <tr>
                    <td>1,012</td>
                    <td>taciti</td>
                    <td>sociosqu</td>
                    <td>ad</td>
                    <td>litora</td>
                  </tr>
                  <tr>
                    <td>1,013</td>
                    <td>torquent</td>
                    <td>per</td>
                    <td>conubia</td>
                    <td>nostra</td>
                  </tr>
                  <tr>
                    <td>1,014</td>
                    <td>per</td>
                    <td>inceptos</td>
                    <td>himenaeos</td>
                    <td>Curabitur</td>
                  </tr>
                  <tr>
                    <td>1,015</td>
                    <td>sodales</td>
                    <td>ligula</td>
                    <td>in</td>
                    <td>libero</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </main>
        </div>
      </div>
    </>
  );
};

export default App;
