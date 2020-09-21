package apigw

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws/apigateway"
	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws/lambda"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

// New takes a pulumi Context, the name of the API, a name for the Gateway, and a lambda function
// to configure as the backend of the gateway, and returns a configured API GW
func New(ctx *pulumi.Context, name string, function *lambda.Function) (*apigateway.RestApi, error) {
	account, err := aws.GetCallerIdentity(ctx)
	if err != nil {
		return nil, err
	}

	region, err := aws.GetRegion(ctx, &aws.GetRegionArgs{})
	if err != nil {
		return nil, err
	}

	// Create a new API Gateway.
	gateway, err := apigateway.NewRestApi(ctx, name, &apigateway.RestApiArgs{
		Name:        pulumi.String(name),
		Description: pulumi.String(fmt.Sprintf("An API Gateway for the %s function", name)),
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
		return nil, err
	}

	// Add a resource to the API Gateway.
	// This makes the API Gateway accept requests on "/{message}".
	apiresource, err := apigateway.NewResource(ctx, "API", &apigateway.ResourceArgs{
		RestApi:  gateway.ID(),
		PathPart: pulumi.String("{proxy+}"),
		ParentId: gateway.RootResourceId,
	}, pulumi.DependsOn([]pulumi.Resource{gateway}))
	if err != nil {
		return nil, err
	}

	// Add a method to the API Gateway.
	_, err = apigateway.NewMethod(ctx, "AnyMethod", &apigateway.MethodArgs{
		HttpMethod:    pulumi.String("ANY"),
		Authorization: pulumi.String("NONE"),
		RestApi:       gateway.ID(),
		ResourceId:    apiresource.ID(),
	}, pulumi.DependsOn([]pulumi.Resource{gateway, apiresource}))
	if err != nil {
		return nil, err
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
		return nil, err
	}

	// Add a resource based policy to the Lambda function.
	// This is the final step and allows AWS API Gateway to communicate with the AWS Lambda function
	permission, err := lambda.NewPermission(ctx, "APIPermission", &lambda.PermissionArgs{
		Action:    pulumi.String("lambda:InvokeFunction"),
		Function:  function.Name,
		Principal: pulumi.String("apigateway.amazonaws.com"),
		SourceArn: pulumi.Sprintf("arn:aws:execute-api:%s:%s:%s/*/*/*", region.Name, account.AccountId, gateway.ID()),
	}, pulumi.DependsOn([]pulumi.Resource{gateway, apiresource, function}))
	if err != nil {
		return nil, err
	}

	// Create a new deployment
	_, err = apigateway.NewDeployment(ctx, "APIDeployment", &apigateway.DeploymentArgs{
		Description:      pulumi.String(fmt.Sprintf("%s API deployment", name)),
		RestApi:          gateway.ID(),
		StageDescription: pulumi.String("Production"),
		StageName:        pulumi.String("prod"),
	}, pulumi.DependsOn([]pulumi.Resource{gateway, apiresource, function, permission}))
	if err != nil {
		return nil, err
	}
	return gateway, nil
}
