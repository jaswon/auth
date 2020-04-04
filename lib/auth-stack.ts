import * as cdk from '@aws-cdk/core';
import * as apigateway from "@aws-cdk/aws-apigateway";
import * as lambda from "@aws-cdk/aws-lambda";

export class AuthStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const handler = new lambda.Function(this, "AuthHandler", {
      runtime: lambda.Runtime.GO_1_X,
      code: lambda.Code.asset("function/bin"),
      handler: "main",
    })

    const api = new apigateway.RestApi(this, "auth-api", {
      restApiName: "Auth Service",
      description: "provides auth",
    })

    const authLambdaIntegration = new apigateway.LambdaIntegration(handler)

    api.root.addMethod("POST", authLambdaIntegration)
    api.root.addMethod("GET", authLambdaIntegration)
  }
}
