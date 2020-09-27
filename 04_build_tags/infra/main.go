package main

import (
	"infra/apigw"
	"infra/function"
	"infra/storage"

	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

var (
	name  = "going-serverless"
	phase = "04"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// set up AWS
		s3Bucket, err := storage.NewAWS(ctx, name, phase)
		if err != nil {
			return err
		}

		// create the lambda function
		lambdaCfg := function.LambdaConfig{
			Name:      name,
			TalkPhase: phase,
			Path:      "./../aws-handler.zip",
			Bucket:    s3Bucket,
		}
		lambda, err := function.NewAWS(ctx, lambdaCfg)
		if err != nil {
			return err
		}

		// Create a new API Gateway backed by our Lambda Function
		gateway, err := apigw.New(ctx, name, phase, lambda)
		if err != nil {
			return err
		}

		// Get region so we can compute our API endpoint
		region, err := aws.GetRegion(ctx, &aws.GetRegionArgs{})
		if err != nil {
			return err
		}
		ctx.Export("AWS invocation URL", pulumi.Sprintf("https://%s.execute-api.%s.amazonaws.com/prod/{message}", gateway.ID(), region.Name))
		// Export function
		return nil
	})
}
