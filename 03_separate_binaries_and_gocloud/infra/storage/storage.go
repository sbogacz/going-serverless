package storage

import (
	"strings"

	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws/s3"
	"github.com/pulumi/pulumi-azure/sdk/v3/go/azure/storage"
	gstorage "github.com/pulumi/pulumi-gcp/sdk/v3/go/gcp/storage"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

// NewAWS takes a pulumi context and returns a private s3 bucket configured to be the
// backing store of the toy project
func NewAWS(ctx *pulumi.Context, name, talkPhase string) (*s3.Bucket, error) {
	// Create a bucket and expose a website index document
	bucket, err := s3.NewBucket(ctx, name, &s3.BucketArgs{
		Tags: pulumi.StringMap{
			"project":    pulumi.String("going-serverless-talk"),
			"talk-phase": pulumi.String(talkPhase),
		},
		Acl: pulumi.String("private"),
		LifecycleRules: &s3.BucketLifecycleRuleArray{
			&s3.BucketLifecycleRuleArgs{
				Enabled: pulumi.Bool(true),
				Expiration: &s3.BucketLifecycleRuleExpirationArgs{
					Days: pulumi.Int(1),
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return bucket, nil
}

// NewAzure takes a pulumi context and returns an Azure storage account configured to be the
// backing store of the toy project
func NewAzure(ctx *pulumi.Context, resourceGroupName pulumi.StringOutput, name, talkPhase string) (*storage.Account, error) {

	// strip dashes from the name as they're not allowed
	name = strings.ReplaceAll(name, "-", "")
	// Create an Azure resource (Storage Account)
	account, err := storage.NewAccount(ctx, name, &storage.AccountArgs{
		ResourceGroupName:      resourceGroupName,
		AccountTier:            pulumi.String("Standard"),
		AccountReplicationType: pulumi.String("LRS"),
		Tags: pulumi.StringMap{
			"project":    pulumi.String("going-serverless-talk"),
			"talk-phase": pulumi.String(talkPhase),
		},
	})
	if err != nil {
		return nil, err
	}

	return account, nil
}

// NewGCP takes a pulumi context and returns a GCP storage bucket configured to be the
// backing store of the toy project
func NewGCP(ctx *pulumi.Context, name, talkPhase string, lifecycleEnabled bool) (*gstorage.Bucket, error) {
	bucketArgs := &gstorage.BucketArgs{
		StorageClass: pulumi.StringPtr("STANDARD"),
		Labels: pulumi.StringMap{
			"project":    pulumi.String("going-serverless-talk"),
			"talk-phase": pulumi.String(talkPhase),
		},
	}
	if lifecycleEnabled {
		bucketArgs.LifecycleRules = gstorage.BucketLifecycleRuleArray{
			gstorage.BucketLifecycleRuleArgs{
				Action: gstorage.BucketLifecycleRuleActionArgs{
					Type: pulumi.String("Delete"),
				},
				Condition: gstorage.BucketLifecycleRuleConditionArgs{
					Age: pulumi.IntPtr(1),
				},
			},
		}
	}

	bucket, err := gstorage.NewBucket(ctx, name, bucketArgs)
	if err != nil {
		return nil, err
	}
	return bucket, nil
}
