import { SSTConfig } from "sst";
import { getDynamodb } from "./cdk/dynamodb";
import Api from "./cdk/api";

export default {
  config(_input) {
    return {
      name: "dtc-common-be",
      region: "us-east-1",
    };
  },
  stacks(app) {
    app.setDefaultFunctionProps({
      runtime: "go",
    });
    app.stack(function Stack({ stack }) {
      const { auth, user, token } = getDynamodb(stack);

      const api = Api(stack);
      api.attachPermissions([auth, user, token]);

      stack.addOutputs({
        ApiEndpoint: api.url,
      });
    });
  },
} satisfies SSTConfig;
