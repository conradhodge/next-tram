import { haveResource, expect as exp } from "@aws-cdk/assert";
import "@aws-cdk/assert/jest";
import { Stack } from "@aws-cdk/core";
import { GetNextTramLambda } from "../../lib/constructs/get-next-tram-lambda";

/**
 * This is an example of a fine-grained test. This will check a single
 * aspect of the construct/stack.
 *
 * @see https://docs.aws.amazon.com/cdk/latest/guide/testing.html
 * @see https://aws.amazon.com/blogs/developer/testing-infrastructure-with-the-aws-cloud-development-kit-cdk/
 */
test("Lambda is created with parameters given", () => {
  const stack = new Stack(undefined, undefined, {
    env: { account: "999999999999", region: "eu-west-2" },
  });
  new GetNextTramLambda(stack, "TestInstance", {
    account: "999999999999",
    region: "eu-west-2",
    apiUsername: "api-username",
    apiPassword: "api-password",
    naptanCode: "111222333",
    memorySize: 128,
    timeout: 10,
  });
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
        },
      },
      MemorySize: 128,
      Timeout: 10,
    })
  );
});

/**
 * This is an example of a fine-grained test. This will check a single
 * aspect of the construct/stack.
 *
 * @see https://docs.aws.amazon.com/cdk/latest/guide/testing.html
 * @see https://aws.amazon.com/blogs/developer/testing-infrastructure-with-the-aws-cloud-development-kit-cdk/
 */
test("Role is created to execute lambda", () => {
  const stack = new Stack(undefined, undefined, {
    env: { account: "999999999999", region: "eu-west-2" },
  });
  new GetNextTramLambda(stack, "TestInstance", {
    account: "999999999999",
    region: "eu-west-2",
    apiUsername: "api-username",
    apiPassword: "api-password",
    naptanCode: "111222333",
    memorySize: 128,
    timeout: 10,
  });
  exp(stack).to(
    haveResource("AWS::IAM::Role", {
      RoleName: "get-next-tram-role",
      AssumeRolePolicyDocument: {
        Statement: [
          {
            Action: "sts:AssumeRole",
            Effect: "Allow",
            Principal: {
              Service: "lambda.amazonaws.com",
            },
          },
        ],
        Version: "2012-10-17",
      },
    })
  );
});
