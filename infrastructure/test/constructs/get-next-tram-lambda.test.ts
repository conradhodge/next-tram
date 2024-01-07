import { Capture, Match, Template } from "aws-cdk-lib/assertions";
import * as cdk from "aws-cdk-lib";
import { GetNextTramLambda } from "../../lib/constructs/get-next-tram-lambda";

/**
 * This is an example of a fine-grained test. This will check a single
 * aspect of the construct/stack.
 *
 * @see https://docs.aws.amazon.com/cdk/latest/guide/testing.html
 * @see https://aws.amazon.com/blogs/developer/testing-infrastructure-with-the-aws-cloud-development-kit-cdk/
 */
test("Lambda is created with parameters given", () => {
  const stack = new cdk.Stack();

  new GetNextTramLambda(stack, "TestInstance", {
    account: "999999999999",
    region: "eu-west-2",
    apiUsername: "api-username",
    apiPassword: "api-password",
    naptanCode: "111222333",
    memorySize: 128,
    timeout: 10,
    localTimezone: "Europe/London",
  });

  const template = Template.fromStack(stack);

  template.resourceCountIs("AWS::Lambda::Function", 1);
  template.hasResourceProperties("AWS::Lambda::Function", {
    FunctionName: "get-next-tram-lambda",
    Description: "Lambda function that will get the next tram",
    Handler: "main",
    Runtime: "provided.al2",
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

/**
 * This is an example of a fine-grained test. This will check a single
 * aspect of the construct/stack.
 *
 * @see https://docs.aws.amazon.com/cdk/latest/guide/testing.html
 * @see https://aws.amazon.com/blogs/developer/testing-infrastructure-with-the-aws-cloud-development-kit-cdk/
 */
test("Role is created to execute lambda", () => {
  const stack = new cdk.Stack();

  new GetNextTramLambda(stack, "TestInstance", {
    account: "999999999999",
    region: "eu-west-2",
    apiUsername: "api-username",
    apiPassword: "api-password",
    naptanCode: "111222333",
    memorySize: 128,
    timeout: 10,
    localTimezone: "Europe/London",
  });

  const template = Template.fromStack(stack);

  template.resourceCountIs("AWS::IAM::Role", 1);
  template.hasResourceProperties("AWS::IAM::Role", {
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
  });
});
