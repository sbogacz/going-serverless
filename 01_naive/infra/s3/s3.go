package s3

import (
	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

// New takes a pulumi context and returns a private s3 bucket configured to be the
// backing store of the toy project
func New(ctx *pulumi.Context, name, talkPhase string) (*s3.Bucket, error) {
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
