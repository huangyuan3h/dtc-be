import { SSTConfig } from "sst";
import { getDynamodb } from "./cdk/dynamodb";
import Api from "./cdk/api";
import { Bucket } from "sst/constructs";

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

      const bucket = new Bucket(stack, "avatar");

      const api = Api(stack);
      api.attachPermissions([auth, user, token]);
      api.bind([bucket]);

      stack.addOutputs({
        ApiEndpoint: api.url,
      });
    });
  },
} satisfies SSTConfig;
