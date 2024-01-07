import { Capture, Match, Template } from "aws-cdk-lib/assertions";
import * as cdk from "aws-cdk-lib";
import { NextTramStack } from "../lib/next-tram-stack";

/**
 * This is an example of a fine-grained test. This will check a single
 * aspect of the construct/stack.
 *
 * @see https://docs.aws.amazon.com/cdk/latest/guide/testing.html
 * @see https://aws.amazon.com/blogs/developer/testing-infrastructure-with-the-aws-cloud-development-kit-cdk/
 */
test("Stack is created with lambda", () => {
  const context = require("./cdk.json");
  const app = new cdk.App(context);
  const stack = new NextTramStack(app, "TestStack", {
    env: { account: "999999999999", region: "eu-west-2" },
  });

  const template = Template.fromStack(stack);

  template.resourceCountIs("AWS::Lambda::Function", 1);
  template.hasResourceProperties("AWS::Lambda::Function", {
    FunctionName: "get-next-tram-lambda",
    Description: "Lambda function that will get the next tram",
    Handler: "bootstrap",
    Runtime: "provided.al2023",
    Environment: {
      Variables: {
        TRAVELINE_API_USERNAME: "api-username",
        TRAVELINE_API_PASSWORD: "api-password",
        NAPTAN_CODE: "111222333",
        LOCAL_TIMEZONE: "Europe/London",
      },
    },
    MemorySize: 128,
    Timeout: 10,
  });
});
