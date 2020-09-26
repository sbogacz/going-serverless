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
	phase = "03"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// set up GCP
		// create code bucket
		gcpCodeBucket, err := storage.NewGCP(ctx, "function-code", phase, false)
		if err != nil {
			return err
		}

		gcpBlobBucket, err := storage.NewGCP(ctx, name, phase, true)
		if err != nil {
			return err
		}

		// create the GCP cloud function
		gcpFunctionCfg := function.GCPFunctionConfig{
			Name:       name,
			TalkPhase:  phase,
			Path:       "./../aws-handler.zip",
			CodeBucket: gcpCodeBucket,
			BlobBucket: gcpBlobBucket,
		}
		gcpFunction, err := function.NewGCP(ctx, gcpFunctionCfg)
		if err != nil {
			return err
		}

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
		ctx.Export("GCP invocation URL", gcpFunction.HttpsTriggerUrl)
		// Export function
		return nil
	})
}

// TODO: Try and add Azure back in
//		// Create an Azure Resource Group
//		resourceGroup, err := core.NewResourceGroup(ctx, "resourceGroup", &core.ResourceGroupArgs{
//			Location: pulumi.String("WestUS"),
//		})
//		if err != nil {
//			return err
//		}
//
//		// Set up storage in Azure
//		account, err := storage.NewAzure(ctx, resourceGroup.Name, name, phase)
//		if err != nil {
//			return err
//		}
//
//		// set up a function in azure
//		azFunctionCfg := function.AzFunctionConfig{
//			Name:              name,
//			TalkPhase:         phase,
//			Path:              "./../az-handler.zip",
//			StorageAccount:    account,
//			ResourceGroupName: resourceGroup.Name,
//		}
//		_, err = function.NewAzure(ctx, azFunctionCfg)
//		if err != nil {
//			return err
//		}
//
//     ctx.Export("connectionString", account.PrimaryConnectionString)
