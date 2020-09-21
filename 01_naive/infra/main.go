package main

import (
	"infra/apigw"
	"infra/lambda"

	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		// create the lambda function
		function, err := lambda.New(ctx, "./../handler.zip", "naive")
		if err != nil {
			return err
		}

		// Create a new API Gateway backed by our Lambda Function
		gateway, err := apigw.New(ctx, "01_naive", function)
		if err != nil {
			return err
		}

		// Get region so we can compute our API endpoint
		region, err := aws.GetRegion(ctx, &aws.GetRegionArgs{})
		if err != nil {
			return err
		}
		ctx.Export("invocation URL", pulumi.Sprintf("https://%s.execute-api.%s.amazonaws.com/prod/{message}", gateway.ID(), region.Name))

		return nil
	})
}
