import { Table, Stack } from "sst/constructs";

export const getDynamodb = (stack: Stack) => {
  const auth = new Table(stack, "auth", {
    fields: {
      email: "string",
      password: "string",
      status: "string", // sendEmail, actived, deactivated
    },
    primaryIndex: { partitionKey: "email" },
  });

  const user = new Table(stack, "user", {
    fields: {
      email: "string",
      avatar: "string",
      userName: "string",
      bio: "string",
    },
    primaryIndex: { partitionKey: "email" },
  });

  const token = new Table(stack, "token", {
    fields: {
      tokenId: "string",
      expireAt: "string",
      isConsumed: "boolean",
      consumedBy: "string", // email
    },
    primaryIndex: { partitionKey: "tokenId" },
    timeToLiveAttribute: "expireAt",
  });

  return { auth, user, token };
};
