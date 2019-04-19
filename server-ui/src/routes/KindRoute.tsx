import React from "react";
import { RouteComponentProps, Switch, Route } from "react-router";
import { KindListRoute, KindListRouteMatch } from "./KindListRoute";
import { GetSchemaResponse } from "../api/meta_pb";
import { KindEditRouteMatch, KindEditRoute } from "./KindEditRoute";

export interface KindRouteMatch {
  kind: string;
}

export interface KindRouteProps extends RouteComponentProps<KindRouteMatch> {
  schema: GetSchemaResponse;
}

export const KindRoute = (outerProps: KindRouteProps) => {
  return (
    <>
      <Switch>
        <Route path="/kind/:kind/create/pending/:idx">
          {(props: RouteComponentProps<KindEditRouteMatch>) => (
            <KindEditRoute {...props} schema={outerProps.schema} />
          )}
        </Route>
        <Route path="/kind/:kind/create">
          {(props: RouteComponentProps<KindEditRouteMatch>) => (
            <KindEditRoute {...props} schema={outerProps.schema} />
          )}
        </Route>
        <Route path="/kind/:kind/edit/:id*">
          {(props: RouteComponentProps<KindEditRouteMatch>) => (
            <KindEditRoute {...props} schema={outerProps.schema} />
          )}
        </Route>
        <Route path="/kind/:kind">
          {(props: RouteComponentProps<KindListRouteMatch>) => (
            <KindListRoute {...props} schema={outerProps.schema} />
          )}
        </Route>
      </Switch>
    </>
  );
};
