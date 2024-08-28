import { Api, Stack } from "sst/constructs";

export default (stack: Stack) => {
  const api = new Api(stack, "api", {
    // cors: {
    //   allowOrigins: ["http://localhost:3000"],
    //   allowCredentials: true,
    //   allowHeaders: ["Authorization"],
    // },
    routes: {
      "GET /": "api/health-check/main.go",
      // "POST /gql": "api/gql/main.go",
      // "POST /register": "api/auth/register.go",
    },
  });
  return api;
};
