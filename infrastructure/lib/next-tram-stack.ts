import { Construct, Stack, StackProps } from "@aws-cdk/core";
import { GetNextTramLambda } from "./constructs/get-next-tram-lambda";

export class NextTramStack extends Stack {
  constructor(scope: Construct, id: string, props?: StackProps) {
    super(scope, id, props);

    // Lambda to go get us the next tram time
    new GetNextTramLambda(this, "GetNextTramLambda", {
      account: this.account,
      region: this.region,
      apiUsername: this.node.tryGetContext("travelineApiUsername"),
      apiPassword: this.node.tryGetContext("travelineApiPassword"),
      naptanCode: this.node.tryGetContext("naptanCode"),
      memorySize: this.node.tryGetContext("getNextTramLambdaMemorySize"),
      timeout: this.node.tryGetContext("getNextTramLambdaTimeout"),
    });
  }
}
