#!/usr/bin/env node
import 'source-map-support/register';
import * as cdk from '@aws-cdk/core';
import * as apigateway from "@aws-cdk/aws-apigateway";
import * as lambda from "@aws-cdk/aws-lambda";
import * as assets from "@aws-cdk/aws-s3-assets";
import * as iam from "@aws-cdk/aws-iam";

export interface Config {
  certArn: string;
  domain: string;
}

import { config } from "./config"

class AuthStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const handler = new lambda.Function(this, "AuthHandler", {
      runtime: lambda.Runtime.GO_1_X,
      code: lambda.Code.asset("assets/function.zip"),
      handler: "main",
    })

    const api = new apigateway.RestApi(this, "auth-api", {
      restApiName: "Auth Service",
      description: "provides auth",
    })

    const apiRole = new iam.Role(this, "auth-api-role", {
      assumedBy: new iam.ServicePrincipal("apigateway.amazonaws.com")
    })

    const authLambdaIntegration = new apigateway.LambdaIntegration(handler)

    api.root.addResource("token").addMethod("POST", authLambdaIntegration)

    const pubkey = new assets.Asset(this, "pubkey", {
      path: "assets/signkey.pub",
    })

    pubkey.grantRead(apiRole)

    const pubkeyIntegration = new apigateway.AwsIntegration({
      service: "s3",
      path: `${pubkey.s3BucketName}/${pubkey.s3ObjectKey}`,
      integrationHttpMethod: "GET",
      options: {
        credentialsRole: apiRole,
        integrationResponses: [
          { statusCode: "200" },
        ],
      }
    })

    api.root.addResource("pubkey").addMethod("GET", pubkeyIntegration, {
      methodResponses: [
        { statusCode: "200" },
      ],
    })
  }
}

const app = new cdk.App();
new AuthStack(app, 'AuthStack');
