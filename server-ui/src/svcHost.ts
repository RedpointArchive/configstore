import { ApiHostConfig } from "./ApiHostProvider";
import { LoginState } from "./Auth";

export function createGrpcPromiseClient<T>(
  constructor: new (
    hostname: string,
    credentials: null | { [index: string]: string },
    options: null | { [index: string]: string }
  ) => T
) {
  return {
    svc: new constructor(
      (JSON.parse(window.localStorage.getItem(
        "apiHostConfig"
      ) as string) as ApiHostConfig).svcHost,
      null,
      null
    ),
    meta:
      window.localStorage.getItem("loginState") === null
        ? {}
        : ({
            Authorization: (JSON.parse(window.localStorage.getItem(
              "loginState"
            ) as string) as LoginState).accessToken
          } as { [s: string]: string })
  };
}
