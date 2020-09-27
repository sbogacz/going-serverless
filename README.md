# going-serverless

Exploring patterns for serverless deployments in Go... or to be more specific, serveless "functions" (e.g. 
AWS Lambda, Google Cloud Functions, and Azure Function Apps). 

## Toy API

In order to explore these patterns we'll be building a somewhat simple API. It will be a blob storage API
that can accept a payload, store it in some backing block storage system under a generated UUID, and can
then access and delete that blob by using that ID.

This repo is structured such that each of the folders represents a holistic, and distinct, approach to
writing Go for such deployments. The "base case" for all the examples is AWS' Lambda, as it is the most
common (as of writing) of the three, and has had native Go support the longest. 

Before we get into what each phase of the code tries to achieve, some notes on building and deploying. 

### Building

Each phase contains a Makefile that can be used to build the code and cleanup code artifacts. The commands
are typically `make build` and `make clean` respectively. This project does use `go mod` so you will need 
a Go installation that's `>=1.11`

### Deploying

For this project I chose to use [Pulumi](https://www.pulumi.com/) rather than [Terraform](https://www.terraform.io/), 
largely out of curiosity. In order to deploy the code in each phase, you'll need a Pulumi account, an AWS account, 
and a GCP account. Once you've configured your local environment with the proper credentials, in order to deploy,
simply navigate to `infra` directory, initialize the Pulumi [stack](https://www.pulumi.com/docs/intro/concepts/stack/)
(just use the default stack name by hitting enter),
and run a Pulumi [update](https://www.pulumi.com/docs/reference/cli/pulumi_up/). 

e.g.

```sh
pwd
/Users/stevenbogacz/code/go/src/github.com/sbogacz/going-serverless/04_build_tags
cd infra
pulumi stack init
...
pulumi up
```

## Phases

### 01 - Naïve

This phase takes a straightforward approach to writing a Lambda backed API, sitting behind an API Gateway in order to
allow us to trigger the Lambda with an HTTP call. The code all lives at the root of the directory, and can only be 
deployed to AWS.

```sh
tree  
.
├── Makefile
├── config.go
├── errors.go
├── errors_test.go
├── go.mod
├── go.sum
├── handler
├── handler.zip
├── infra
│   ├── Pulumi.dev.yaml
│   ├── Pulumi.yaml
│   ├── apigw
│   │   └── apigw.go
│   ├── go.mod
│   ├── go.sum
│   ├── lambda
│   │   └── lambda.go
│   ├── main.go
│   └── s3
│       └── s3.go
├── main.go
└── main_test.go
```

### 02 - Separate Binaries

Here we take a more "traditional" Go package structure, splitting out our `package main` under a `cmd` directory
and our business logic in `internal`.



