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

Directory structure:
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
and our business logic in `internal`. In `internal/toy` we have a deployment agnostic toy API server, and in `internal/store` we define
both an S3 backed store as well as an in-memory implementation. 

The two targets under `cmd`, `cmd/http` and `cmd/lambda`, define two builds: the first enables us to more quickly iterate on code by allowing us
to deploy and test locally, while the second exposes a Lambda-specific entrypoint. 


Directory structure:
```sh
tree -L 2   
.
├── Makefile
├── cmd
│   ├── http
│   └── lambda
├── go.mod
├── go.sum
├── infra
│   ├── Pulumi.dev.yaml
│   ├── Pulumi.yaml
│   ├── apigw
│   ├── go.mod
│   ├── go.sum
│   ├── lambda
│   ├── main.go
│   └── s3
└── internal
    ├── httperrs
    ├── s3store
    └── toy
```

### 03 - Separate Binaries Using Go Cloud SDK

In the third phase we iterate on our separate binaries approach by also taking advantage of the [Go CDK libraries](https://gocloud.dev/)
in order to replace our old `internal/store` package. Not only did that let us delete a good amount of our code, but it also allowed us
to easily create a deployment compatible with GCP as well as AWS. 

Some deployment requirements means that we need to rely on some work from our build process to create the proper file structure for a 
[Google Cloud Function](https://cloud.google.com/functions/docs/writing#structuring_source_code). 

Directory structure:
```sh
tree -L 2  
.
├── Makefile
├── aws-handler.zip
├── cmd
│   ├── azure
│   ├── gcp
│   ├── http
│   └── lambda
├── gcp-handler.zip
├── go.mod
├── go.sum
├── handler
├── infra
│   ├── Pulumi.dev.yaml
│   ├── Pulumi.yaml
│   ├── README.md
│   ├── apigw
│   ├── function
│   ├── go.mod
│   ├── go.sum
│   ├── main.go
│   └── storage
└── internal
    ├── httperrs
    └── toy
```

### 04 - Using Build Tags

Finaly we take a different approach, and take a leaf out of the Go cross-compilation playbook: [build tags](https://golang.org/pkg/go/build/#hdr-Build_Constraints) 
(there's nothing cross-compilation specific about them, anecdotally that's just where they're most commonly used). 

This allows us to write a single `package main` and lets us write less duplicate code, by deferring the decision of which implementation for our 
`serve` and `getStore` functions to use until compile time, which we toggle by using `go build -tags aws`. 

Directory structure:
```sh
tree -L 2   
.
├── Makefile
├── aws-handler.zip
├── cmd
│   └── toy
├── go.mod
├── go.sum
├── handler
├── infra
│   ├── Pulumi.dev.yaml
│   ├── Pulumi.yaml
│   ├── README.md
│   ├── apigw
│   ├── function
│   ├── go.mod
│   ├── go.sum
│   ├── main.go
│   └── storage
└── internal
    ├── httperrs
    └── toy
```
