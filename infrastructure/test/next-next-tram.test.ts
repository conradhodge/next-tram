import { haveResource, expect as exp, haveType } from "@aws-cdk/assert";
import "@aws-cdk/assert/jest";
import { App } from "@aws-cdk/core";
import { NextTramStack } from "../lib/next-tram-stack"

/**
 * This is an example of a fine-grained test. This will check a single
 * aspect of the construct/stack.
 *
 * @see https://docs.aws.amazon.com/cdk/latest/guide/testing.html
 * @see https://aws.amazon.com/blogs/developer/testing-infrastructure-with-the-aws-cloud-development-kit-cdk/
 */
test("Stack is created with lambda", () => {
  const context = require("./cdk.json");
  const app = new App(context);
  const stack = new NextTramStack(app, "MyTestStack", { env: { account: "999999999999", region: "eu-west-2" } });
  exp(stack).to(
    haveResource("AWS::Lambda::Function", {
      FunctionName: "get-next-tram-lambda",
      Description: "Lambda function that will get the next tram",
      Handler: "main",
      Runtime: "go1.x",
      Environment: {
        Variables: {
          TRAVELINE_API_USERNAME: "api-username",
          TRAVELINE_API_PASSWORD: "api-password",
          NAPTAN_CODE: "111222333",
        }
      },
      MemorySize: 128,
      Timeout: 10
    })
  );
});
