module infra

go 1.14

require (
	github.com/pulumi/pulumi-aws/sdk/v3 v3.3.0
	github.com/pulumi/pulumi/sdk/v2 v2.10.1
	github.com/sbogacz/going-serverless/01_naive/infra/lambda v0.0.0-00010101000000-000000000000
)

replace github.com/sbogacz/going-serverless/01_naive/infra/lambda => ./lambda
