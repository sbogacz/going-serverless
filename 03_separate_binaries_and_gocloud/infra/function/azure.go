package function

import (
	"path/filepath"

	"github.com/pulumi/pulumi-azure/sdk/v3/go/azure/appservice"
	"github.com/pulumi/pulumi-azure/sdk/v3/go/azure/storage"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

// AzFunctionConfig holds the various bits that we'll need to correctly
// configure an Azure FunctionApp
type AzFunctionConfig struct {
	Name              string
	Path              string
	TalkPhase         string
	StorageAccount    *storage.Account
	ResourceGroupName pulumi.StringOutput
}

// NewAzure takes a pulumi context, cfg.Path to the zipped binary, and cfg.Name and it
// returns a lambda.Function on success
func NewAzure(ctx *pulumi.Context, cfg AzFunctionConfig) (*appservice.FunctionApp, error) {
	// first upload the code
	if err := uploadCode(ctx, cfg); err != nil {
		return nil, err
	}

	return nil, nil
}

func uploadCode(ctx *pulumi.Context, cfg AzFunctionConfig) error {
	// first create a storage container for the function releases
	container, err := storage.NewContainer(ctx, "going-serverless-releases", &storage.ContainerArgs{
		StorageAccountName:  cfg.StorageAccount.Name,
		ContainerAccessType: pulumi.String("private"),
	})
	if err != nil {
		return err
	}

	// upload function definition zip file
	_, err = storage.NewBlob(ctx, filepath.Base(cfg.Path), &storage.BlobArgs{
		StorageAccountName:   cfg.StorageAccount.Name,
		StorageContainerName: container.Name,
		Type:                 pulumi.String("Block"),
		Source:               pulumi.NewFileAsset(cfg.Path),
	})
	if err != nil {
		return err
	}
	return nil
}

//func setupServiceAccount(ctx *pulumi.Context, cfg AzFunctionConfig) ()
