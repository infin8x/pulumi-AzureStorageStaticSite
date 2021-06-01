module github.com/infin8x/pulumi-AzureStorageStaticSite

go 1.16

require (
	github.com/pkg/errors v0.9.1
	github.com/pulumi/pulumi-azure-native/sdk v1.9.0
	github.com/pulumi/pulumi/pkg/v3 v3.0.1-0.20210419234039-6a33b4b7ee41
	github.com/pulumi/pulumi/sdk/v3 v3.3.1
)

replace github.com/pulumi/pulumi/pkg/v3 => ../../../pulumi/pulumi/pkg
