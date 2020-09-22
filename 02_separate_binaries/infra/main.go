package main

import (
	"infra/apigw"
	"infra/lambda"
	"infra/s3"

	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

var (
	name  = "going-serverless"
	phase = "02"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		s3Bucket, err := s3.New(ctx, name, phase)
		if err != nil {
			return err
		}

		// create the lambda function
		lambdaCfg := lambda.Config{
			Name:      name,
			TalkPhase: phase,
			Path:      "./../handler.zip",
			Bucket:    s3Bucket,
		}
		function, err := lambda.New(ctx, lambdaCfg)
		if err != nil {
			return err
		}

		// Create a new API Gateway backed by our Lambda Function
		gateway, err := apigw.New(ctx, name, phase, function)
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
