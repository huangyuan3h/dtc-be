import { Api, Stack } from "sst/constructs";
import { Certificate } from "aws-cdk-lib/aws-certificatemanager";

const certArn =
  "arn:aws:acm:us-east-1:319653899185:certificate/dde24c52-09e4-4058-b1d0-2a7769f24e3a";

export default (stack: Stack) => {
  const customDomain = {
    domainName: "api.it-t.xyz",
    isExternalDomain: true,
    cdk: {
      certificate: Certificate.fromCertificateArn(stack, "MyCert", certArn),
    },
  };

  const api = new Api(stack, "api", {
    cors: {
      allowOrigins: ["http://localhost:3000", "https://www.it-t.xyz"],
      allowCredentials: true,
      allowHeaders: ["Authorization"],
    },
    customDomain: stack.stage === "production" ? customDomain : undefined,
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
