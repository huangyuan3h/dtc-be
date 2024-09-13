import { Api, Stack } from "sst/constructs";

export default (stack: Stack) => {
  const api = new Api(stack, "api", {
    cors: {
      allowOrigins: ["http://localhost:3000", "https://www.it-t.xyz"],
      allowCredentials: true,
      allowHeaders: ["Authorization"],
    },
    routes: {
      "GET /": "api/health-check/main.go",
      "POST /gql": "api/gql/main.go",
      "POST /register": {
        function: {
          handler: "api/auth/register/main.go",
          timeout: 10,
          environment: { EmailToken: process.env.EmailToken ?? "" },
        },
      },
      "POST /verify-email-token": {
        function: {
          handler: "api/auth/verify-email-token/main.go",
          timeout: 10,
        },
      },
      "POST /create-account": {
        function: {
          handler: "api/auth/create-account/main.go",
          timeout: 10,
          environment: {
            AUTH_SECRET: process.env.AUTH_SECRET ?? "",
            JWT_SECRET: process.env.JWT_SECRET ?? "",
          },
        },
      },
      "POST /login": {
        function: {
          handler: "api/auth/login/main.go",
          timeout: 10,
          environment: {
            AUTH_SECRET: process.env.AUTH_SECRET ?? "",
            JWT_SECRET: process.env.JWT_SECRET ?? "",
          },
        },
      },
    },
  });
  return api;
};
