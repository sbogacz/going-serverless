package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws/apigateway"
	plambda "github.com/pulumi/pulumi-aws/sdk/v3/go/aws/lambda"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
	"github.com/sbogacz/going-serverless/01_naive/infra/lambda"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		account, err := aws.GetCallerIdentity(ctx)
		if err != nil {
			return err
		}

		region, err := aws.GetRegion(ctx, &aws.GetRegionArgs{})
		if err != nil {
			return err
		}

		function, err := lambda.New(ctx, "./../handler.zip", "naive")
		if err != nil {
			return err
		}

		// Create a new API Gateway.
		gateway, err := apigateway.NewRestApi(ctx, "UpperCaseGateway", &apigateway.RestApiArgs{
			Name:        pulumi.String("UpperCaseGateway"),
			Description: pulumi.String("An API Gateway for the UpperCase function"),
			Policy: pulumi.String(`{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    },
    {
      "Action": "execute-api:Invoke",
      "Resource": "*",
      "Principal": "*",
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}`)})
		if err != nil {
			return err
		}

		// Add a resource to the API Gateway.
		// This makes the API Gateway accept requests on "/{message}".
		apiresource, err := apigateway.NewResource(ctx, "UpperAPI", &apigateway.ResourceArgs{
			RestApi:  gateway.ID(),
			PathPart: pulumi.String("{proxy+}"),
			ParentId: gateway.RootResourceId,
		}, pulumi.DependsOn([]pulumi.Resource{gateway}))
		if err != nil {
			return err
		}

		// Add a method to the API Gateway.
		_, err = apigateway.NewMethod(ctx, "AnyMethod", &apigateway.MethodArgs{
			HttpMethod:    pulumi.String("ANY"),
			Authorization: pulumi.String("NONE"),
			RestApi:       gateway.ID(),
			ResourceId:    apiresource.ID(),
		}, pulumi.DependsOn([]pulumi.Resource{gateway, apiresource}))
		if err != nil {
			return err
		}

		// Add an integration to the API Gateway.
		// This makes communication between the API Gateway and the Lambda function work
		_, err = apigateway.NewIntegration(ctx, "LambdaIntegration", &apigateway.IntegrationArgs{
			HttpMethod:            pulumi.String("ANY"),
			IntegrationHttpMethod: pulumi.String("POST"),
			ResourceId:            apiresource.ID(),
			RestApi:               gateway.ID(),
			Type:                  pulumi.String("AWS_PROXY"),
			Uri:                   function.InvokeArn,
		}, pulumi.DependsOn([]pulumi.Resource{gateway, apiresource, function}))
		if err != nil {
			return err
		}

		// Add a resource based policy to the Lambda function.
		// This is the final step and allows AWS API Gateway to communicate with the AWS Lambda function
		permission, err := plambda.NewPermission(ctx, "APIPermission", &plambda.PermissionArgs{
			Action:    pulumi.String("lambda:InvokeFunction"),
			Function:  function.Name,
			Principal: pulumi.String("apigateway.amazonaws.com"),
			SourceArn: pulumi.Sprintf("arn:aws:execute-api:%s:%s:%s/*/*/*", region.Name, account.AccountId, gateway.ID()),
		}, pulumi.DependsOn([]pulumi.Resource{gateway, apiresource, function}))
		if err != nil {
			return err
		}

		// Create a new deployment
		_, err = apigateway.NewDeployment(ctx, "APIDeployment", &apigateway.DeploymentArgs{
			Description:      pulumi.String("UpperCase API deployment"),
			RestApi:          gateway.ID(),
			StageDescription: pulumi.String("Production"),
			StageName:        pulumi.String("prod"),
		}, pulumi.DependsOn([]pulumi.Resource{gateway, apiresource, function, permission}))
		if err != nil {
			return err
		}

		ctx.Export("invocation URL", pulumi.Sprintf("https://%s.execute-api.%s.amazonaws.com/prod/{message}", gateway.ID(), region.Name))

		return nil
	})
}
