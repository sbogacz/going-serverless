package lambda

import (
	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws/lambda"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

// New takes a pulumi context, path to the zipped binary, and name and it
// returns a lambda.Function on success
func New(ctx *pulumi.Context, path, name string) (*lambda.Function, error) {
	roleDependencies, roleARN, err := createRole(ctx)
	if err != nil {
		return nil, err
	}

	// Set arguments for constructing the function resource.
	args := &lambda.FunctionArgs{
		Handler: pulumi.String("handler"),
		Role:    *roleARN,
		Runtime: pulumi.String("go1.x"),
		Code:    pulumi.NewFileArchive(path),
	}

	// Create the lambda using the args.
	function, err := lambda.NewFunction(
		ctx,
		name,
		args,
		pulumi.DependsOn(roleDependencies),
	)
	if err != nil {
		return nil, err
	}
	return function, nil
}

func createRole(ctx *pulumi.Context) ([]pulumi.Resource, *pulumi.StringOutput, error) {

	// Create an IAM role.
	role, err := iam.NewRole(ctx, "task-exec-role", &iam.RoleArgs{
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
	logPolicy, err := iam.NewRolePolicy(ctx, "lambda-log-policy", &iam.RolePolicyArgs{
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

	return []pulumi.Resource{logPolicy}, &role.Arn, nil
}
