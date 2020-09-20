module github.com/sbogacz/going-serverless/01_naive

go 1.14

require (
	github.com/aws/aws-lambda-go v1.19.1
	github.com/aws/aws-sdk-go-v2 v0.24.0
	github.com/gofrs/uuid v3.3.0+incompatible
	github.com/kr/pretty v0.2.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.6.0
	github.com/stretchr/testify v1.6.1
	golang.org/x/sys v0.0.0-20200317113312-5766fd39f98d // indirect
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
)

replace infra v0.0.0 => ./infra
