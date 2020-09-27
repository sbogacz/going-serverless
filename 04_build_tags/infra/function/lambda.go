package function

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws/lambda"
	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

// LambdaConfig holds the various bits that we'll need to correctly
// configure an AWS Lambda Function
type LambdaConfig struct {
	Name      string
	Path      string
	TalkPhase string
	Bucket    *s3.Bucket
}

// NewAWS takes a pulumi context, cfg.Path to the zipped binary, and cfg.Name and it
// returns a lambda.Function on success
func NewAWS(ctx *pulumi.Context, cfg LambdaConfig) (*lambda.Function, error) {
	name := fmt.Sprintf("%s-%s", cfg.Name, cfg.TalkPhase)

	roleDependencies, roleARN, err := createRole(ctx, cfg)
	if err != nil {
		return nil, err
	}

	// Set arguments for constructing the function resource.
	args := &lambda.FunctionArgs{
		Handler: pulumi.String("handler"),
		Role:    *roleARN,
		Runtime: pulumi.String("go1.x"),
		Code:    pulumi.NewFileArchive(cfg.Path),
		Timeout: pulumi.Int(5),
		Tags: pulumi.StringMap{
			"project":    pulumi.String("going-serverless-talk"),
			"talk-phase": pulumi.String(cfg.TalkPhase),
		},
		Environment: lambda.FunctionEnvironmentArgs{
			Variables: pulumi.StringMap{
				"BUCKET_NAME": cfg.Bucket.Bucket,
			},
		},
	}

	// Create the lambda using the args.
	function, err := lambda.NewFunction(
		ctx,
		name,
		args,
		pulumi.DependsOn(append(roleDependencies, cfg.Bucket)),
	)
	if err != nil {
		return nil, err
	}
	return function, nil
}

func createRole(ctx *pulumi.Context, cfg LambdaConfig) ([]pulumi.Resource, *pulumi.StringOutput, error) {
	name := fmt.Sprintf("%s-%s", cfg.Name, cfg.TalkPhase)

	// Create an IAM role.
	role, err := iam.NewRole(ctx, fmt.Sprintf("%s-lambda-exec-role", name), &iam.RoleArgs{
		AssumeRolePolicy: pulumi.String(`{
				"Version": "2012-10-17",
				"Statement": [{
					"Sid": "",
					"Effect": "Allow",
					"Principal": {
						"Service": "lambda.amazonaws.com"
					},
					"Action": "sts:AssumeRole"
				}]
			}`),
	})
	if err != nil {
		return nil, nil, err
	}

	// Attach a policy to allow writing logs to CloudWatch
	logPolicy, err := iam.NewRolePolicy(ctx, fmt.Sprintf("%s-lambda-log-policy", name), &iam.RolePolicyArgs{
		Role: role.Name,
		Policy: pulumi.String(`{
                "Version": "2012-10-17",
                "Statement": [{
                    "Effect": "Allow",
                    "Action": [
                        "logs:CreateLogGroup",
                        "logs:CreateLogStream",
                        "logs:PutLogEvents"
                    ],
                    "Resource": "arn:aws:logs:*:*:*"
                }]
            }`),
	})
	if err != nil {
		return nil, nil, err
	}

	s3PolicyFmt := `{
                "Version": "2012-10-17",
                "Statement": [
					{
						"Effect": "Allow",
						"Action": [
							"s3:ListBucket"
						],
						"Resource": "%s"
					},
					{
						"Action": [
							"s3:PutObject",
							"s3:GetObject",
							"s3:DeleteObject"
						],
						"Effect": "Allow",
						"Resource": "%s/*"
					}
				]
            }`
	// attach a policy for the S3 bucket
	lambdaPolicy, err := iam.NewRolePolicy(ctx, fmt.Sprintf("%s-lambda-s3-policy", name), &iam.RolePolicyArgs{
		Role:   role.Name,
		Policy: pulumi.Sprintf(s3PolicyFmt, cfg.Bucket.Arn, cfg.Bucket.Arn),
	})
	if err != nil {
		return nil, nil, err
	}

	return []pulumi.Resource{logPolicy, lambdaPolicy}, &role.Arn, nil
}
