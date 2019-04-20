import * as React from "react";
import logo from "./logo.svg";
import { useState } from "react";

export interface ApiHostConfig {
  svcHost: string;
  auth0ClientConfig: any;
}

export const ApiHostContext = React.createContext<ApiHostConfig | undefined>(
  undefined
);

export const ApiHostProvider = (props: { children: React.ReactNode }) => {
  const [existingConfig, setExistingConfig] = useState(
    window.localStorage.getItem("apiHostConfig")
  );
  const [apiEndpointUrl, setApiEndpointUrl] = useState(
    `${window.location.protocol}//${window.location.hostname}:13390`
  );
  const [auth0ClientConfigJson, setAuth0ClientConfigJson] = useState("");

  if (existingConfig !== null) {
    return (
      <ApiHostContext.Provider value={JSON.parse(existingConfig)}>
        {props.children}
      </ApiHostContext.Provider>
    );
  }

  const saveConfig = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    window.localStorage.setItem(
      "apiHostConfig",
      JSON.stringify({
        svcHost: apiEndpointUrl,
        auth0ClientConfig:
          auth0ClientConfigJson.trim() === ""
            ? null
            : JSON.parse(auth0ClientConfigJson)
      })
    );
    setExistingConfig(window.localStorage.getItem("apiHostConfig"));
  };

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
        <button className="btn btn-lg btn-primary btn-block mt-2" type="submit">
          Connect
        </button>
        <p className="mt-5 mb-3 text-muted">&copy; 2019</p>
      </form>
    </div>
  );
};
