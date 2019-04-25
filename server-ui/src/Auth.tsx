import { ApiHostContext } from "./ApiHostProvider";
import * as React from "react";
import auth0 from "auth0-js";
import { Switch, Route, RouteComponentProps, withRouter } from "react-router";
import { g } from "./core";

export interface AuthState {
  accessToken: string;
  idToken: string;
}

export const AuthContext = React.createContext<AuthState | undefined>(
  undefined
);

export type LoginState =
  | {
      isLoggedIn: true;
      accessToken: string;
      idToken: string;
      expiresAt: number;
    }
  | {
      isLoggedIn: false;
      accessToken: null;
      idToken: null;
      expiresAt: null;
    };

interface AuthComponentProps extends RouteComponentProps<{}> {
  auth0ClientConfig: any;
}

const AuthComponent = withRouter(
  class AuthComponent extends React.Component<AuthComponentProps, LoginState> {
    constructor(props: AuthComponentProps) {
      super(props);

      const defaultSerializedState = JSON.stringify({
        isLoggedIn: false,
        accessToken: null,
        idToken: null,
        expiresAt: null
      });
      const serializedState =
        window.localStorage.getItem("loginState") || defaultSerializedState;
      window.localStorage.setItem("loginState", serializedState);

      this.state = JSON.parse(serializedState) as LoginState;
    }

    private getAuth0Instance = () => {
      const auth0ClientConfigSerialized = JSON.stringify(
        this.props.auth0ClientConfig
      );
      if (!auth0InstanceCache.has(auth0ClientConfigSerialized)) {
        console.log("auth0 instance created");
        auth0InstanceCache.set(
          auth0ClientConfigSerialized,
          new auth0.WebAuth(this.props.auth0ClientConfig)
        );
      }
      return auth0InstanceCache.get(auth0ClientConfigSerialized);
    };

    private setLoginState(loginState: LoginState, cb?: () => void) {
      window.localStorage.setItem("loginState", JSON.stringify(loginState));
      this.setState(loginState, cb);
    }

    private logout = () => {
      this.setLoginState(
        {
          isLoggedIn: false,
          accessToken: null,
          idToken: null,
          expiresAt: null
        },
        () => {
          this.props.history.replace("/");
        }
      );
    };

    private authCallbackHandler = (props: RouteComponentProps<{}>) => {
      if (/access_token|id_token|error/.test(props.location.hash)) {
        if (this.state.isLoggedIn) {
          // we are already logged in; this route can be hit twice when React
          // refreshes due to the setLoginState below, so we just wait for the
          // history redirect to apply instead.
          props.history.replace("/");
          return <>Redirecting...</>;
        }

        this.getAuth0Instance().parseHash(async (err: any, authResult: any) => {
          if (authResult && authResult.accessToken && authResult.idToken) {
            const expiresAt =
              authResult.expiresIn * 1000 + new Date().getTime();
            this.setLoginState(
              {
                isLoggedIn: true,
                accessToken: authResult.accessToken,
                idToken: authResult.idToken,
                expiresAt: expiresAt
              },
              () => {
                this.props.history.replace("/");
              }
            );
          } else if (err) {
            console.error("error during hash parse");
            console.error(err);
            this.logout();
          }
        });
      }

      return <>Please wait...</>;
    };

    private authLogoutHandler = (props: RouteComponentProps<{}>) => {
      this.logout();
      return <>Please wait...</>;
    };

    private routeHandler = () => {
      if (
        !this.state.isLoggedIn ||
        this.state.expiresAt <= new Date().getTime()
      ) {
        // ugh, hack!
        if (window.location.hostname.includes("auth0.com")) {
          return <>Please wait...</>;
        }

        // make sure we clear any existing settings (like expiry) before
        // proceeding
        window.localStorage.removeItem("loginState");

        this.getAuth0Instance().authorize();
        return <>Please wait...</>;
      }

      return (
        <AuthContext.Provider
          value={{
            accessToken: this.state.accessToken,
            idToken: this.state.idToken
          }}
        >
          {this.props.children}
        </AuthContext.Provider>
      );
    };

    public render() {
      return (
        <Switch>
          <Route path="/callback" exact component={this.authCallbackHandler} />
          <Route path="/logout" exact component={this.authLogoutHandler} />
          <Route path="*" component={this.routeHandler} />
        </Switch>
      );
    }
  }
);

const auth0InstanceCache = new Map<string, any>();

export const Auth = (props: { children: React.ReactNode }) => {
  return (
    <ApiHostContext.Consumer>
      {value => {
        if (g(value).auth0ClientConfig !== null) {
          return (
            <AuthComponent auth0ClientConfig={g(value).auth0ClientConfig}>
              {props.children}
            </AuthComponent>
          );
        } else {
          return <>{props.children}</>;
        }
      }}
    </ApiHostContext.Consumer>
  );
};
