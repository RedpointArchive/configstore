import * as React from "react";
import logo from "./logo.svg";
import { useState } from "react";
import {
  withRouter,
  RouteComponentProps,
  Redirect,
  Switch,
  Route
} from "react-router";

export interface ApiHostConfig {
  svcHost: string;
  auth0ClientConfig: any;
  hideTitlebar: boolean;
}

export const ApiHostContext = React.createContext<ApiHostConfig | undefined>(
  undefined
);

export const ApiHostProvider = withRouter(
  (props: { children: React.ReactNode } & RouteComponentProps) => {
    const [existingConfig, setExistingConfig] = useState(
      window.localStorage.getItem("apiHostConfig")
    );
    let [apiEndpointUrl, setApiEndpointUrl] = useState(
      `${window.location.protocol}//${window.location.hostname}:13390`
    );
    let [hideTitlebar, setHideTitlebar] = useState(false);
    let [auth0ClientConfigJson, setAuth0ClientConfigJson] = useState("");

    const qs = new URLSearchParams(props.location.search);

    if (existingConfig !== null) {
      if (
        qs.has("api_endpoint_url") &&
        qs.has("auth0_client_config_json") &&
        qs.has("hide_titlebar")
      ) {
        return <Redirect to="/" />;
      }

      return (
        <ApiHostContext.Provider value={JSON.parse(existingConfig)}>
          <Switch>
            <Route
              path="/reset_config"
              exact
              render={() => {
                window.localStorage.removeItem("apiHostConfig");
                setExistingConfig(null);
                return <Redirect to="/" />;
              }}
            />
            <Route path="*">{props.children}</Route>
          </Switch>
        </ApiHostContext.Provider>
      );
    }

    const setConfig = () => {
      window.localStorage.setItem(
        "apiHostConfig",
        JSON.stringify({
          svcHost: apiEndpointUrl,
          auth0ClientConfig: parseAuth0ClientConfigJson(auth0ClientConfigJson),
          hideTitlebar: hideTitlebar
        })
      );
      setExistingConfig(window.localStorage.getItem("apiHostConfig"));
    };

    if (
      qs.has("api_endpoint_url") &&
      qs.has("auth0_client_config_json") &&
      qs.has("hide_titlebar")
    ) {
      apiEndpointUrl = atob(qs.get("api_endpoint_url") as string);
      auth0ClientConfigJson = atob(qs.get(
        "auth0_client_config_json"
      ) as string);
      hideTitlebar = qs.get("hide_titlebar") === "true";
      setConfig();
      return <Redirect to="/" />;
    }

    const saveConfig = (e: React.FormEvent<HTMLFormElement>) => {
      e.preventDefault();
      setConfig();
    };

    const quickUrl = `${window.location.protocol}//${
      window.location.host
    }/?api_endpoint_url=${encodeURIComponent(
      btoa(apiEndpointUrl)
    )}&auth0_client_config_json=${encodeURIComponent(
      btoa(auth0ClientConfigJson)
    )}&hide_titlebar=${hideTitlebar ? "true" : "false"}`;

    return (
      <div className="form-holder">
        <form onSubmit={saveConfig} className="form-signin">
          <img className="mb-4" src={logo} alt="" width="72" height="72" />
          <h1 className="h3 mb-3 font-weight-normal">configstore</h1>
          <label htmlFor="inputApiUrl" className="sr-only">
            API endpoint URL
          </label>
          <input
            type="text"
            id="inputApiUrl"
            className="form-control"
            placeholder="API endpoint URL"
            required
            autoFocus
            value={apiEndpointUrl}
            onChange={e => setApiEndpointUrl(e.target.value)}
          />
          <label htmlFor="inputPassword" className="sr-only">
            Auth0 client configuration JSON (leave blank for no authentication)
          </label>
          <textarea
            rows={8}
            id="inputAuth0Config"
            className="form-control"
            placeholder="Auth0 client configuration (JSON; leave blank for no authentication)"
            value={auth0ClientConfigJson}
            onChange={e => setAuth0ClientConfigJson(e.target.value)}
          />
          <div className="checkbox mb-3">
            <label>
              <input
                type="checkbox"
                value="true"
                checked={hideTitlebar}
                onChange={e => setHideTitlebar(e.target.checked)}
              />{" "}
              Hide titlebar
            </label>
          </div>
          <button
            className="btn btn-lg btn-primary btn-block mt-2"
            type="submit"
          >
            Connect
          </button>
          <p className="mt-5 mb-3 text-muted">
            Automatic config: <a href={quickUrl}>(Copy this link)</a>
          </p>
          <p className="mt-5 mb-3 text-muted">&copy; 2019</p>
        </form>
      </div>
    );
  }
);

function parseAuth0ClientConfigJson(auth0ClientConfigJson: string) {
  let auth0ClientConfig =
    auth0ClientConfigJson.trim() === ""
      ? null
      : JSON.parse(auth0ClientConfigJson);
  if (auth0ClientConfig != null) {
    if (auth0ClientConfig.redirectUri[0] === "/") {
      // Prepend the protocol and host if it's a relative URL.
      auth0ClientConfig.redirectUri = `${window.location.protocol}//${
        window.location.host
      }${auth0ClientConfig.redirectUri}`;
    }
  }
  return auth0ClientConfig;
}
