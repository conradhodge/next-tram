import { Construct, Duration, CfnOutput } from "@aws-cdk/core";
import { PolicyDocument, PolicyStatement, Role, ServicePrincipal } from "@aws-cdk/aws-iam";
import { Code, Function, Runtime } from "@aws-cdk/aws-lambda";

/**
 * The props that will be required to instantiate a GetNextTramLambda construct
 */
export interface GetNextTramLambdaProps {
  /**
   * The AWS Account the lambda is running in
   */
  readonly account: string;

  /**
   * The AWS Region the lambda is running in
   */
  readonly region: string;

  /**
   * The username to use in the Traveline API
   */
  readonly apiUsername: string;

  /**
   * The password to use in the Traveline API
   */
  readonly apiPassword: string;

  /**
   * NaPTAN code of the tram stop for the next tram times
   */
  readonly naptanCode: string;

  /**
   * The amount of memory, in MB, that is allocated to the Lambda function.
   */
  readonly memorySize: number

  /**
   * The function execution time (in seconds) after which Lambda terminates
   * the function.
   */
  readonly timeout: number;
}

/**
 * Create a lambda that uses the Traveline API to fetch the time of the next tram from a specific tram stop as
 * defined by a NaPTAN code.
 */
export class GetNextTramLambda extends Construct {
  /**
   * Constructor for the get next tram lambda
   *
   * @param scope
   * @param id
   * @param props
   */
  constructor(scope: Construct, id: string, props: GetNextTramLambdaProps) {
    super(scope, id);

    const name = "get-next-tram";

    // Define a policy statement to access the lambda log group
    const logsPolicyDocument = new PolicyDocument({
      statements: [
        new PolicyStatement({
          actions: ["logs:CreateLogGroup", "logs:CreateLogStream", "logs:PutLogEvents"],
          resources: ["arn:aws:logs:" + props.region + ":" + props.account + ":log-group:/aws/lambda/*:*"]
        })
      ]
    });

    // Define a role to execute the lambda
    const role = new Role(this, "GetNextTramLambdaRole", {
      roleName: name + "-role",
      assumedBy: new ServicePrincipal("lambda.amazonaws.com"),
      inlinePolicies: { logsPolicyDocument }
    });

    // Define a lambda function that will get the next tram
    const lambdaFunction = new Function(this, "GetNextTramLambda", {
      functionName: name + "-lambda",
      description: "Lambda function that will get the next tram",
      code: Code.fromAsset("src/lambda/get-next-tram/handler.zip"),
      handler: "main",
      runtime: Runtime.GO_1_X,
      memorySize: props.memorySize,
      timeout: Duration.seconds(props.timeout),
      environment: {
        TRAVELINE_API_USERNAME: props.apiUsername,
        TRAVELINE_API_PASSWORD: props.apiPassword,
        NAPTAN_CODE: props.naptanCode
      },
      role: role
    });

    new CfnOutput(this, "LambdaARN", {
      description: "Get next tram Lambda ARN",
      value: lambdaFunction.functionArn,
      exportName: name + "-arn"
    });
  }
}
