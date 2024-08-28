import { Api, Stack } from "sst/constructs";

export default (stack: Stack) => {
  const api = new Api(stack, "api", {
    cors: {
      allowOrigins: ["http://localhost:3000"],
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
    },
  });
  return api;
};
