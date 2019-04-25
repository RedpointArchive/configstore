import React from "react";
import ReactDOM from "react-dom";
import "./index.css";
import { ApiHostProvider } from "./ApiHostProvider";
import * as serviceWorker from "./serviceWorker";
import App from "./App";
import { Auth } from "./Auth";
import { BrowserRouter } from "react-router-dom";

import "bootstrap/dist/css/bootstrap.css";
import "./App.css";

ReactDOM.render(
  <BrowserRouter>
    <ApiHostProvider>
      <Auth>
        <App />
      </Auth>
    </ApiHostProvider>
  </BrowserRouter>,
  document.getElementById("root")
);

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: http://bit.ly/CRA-PWA
serviceWorker.unregister();
