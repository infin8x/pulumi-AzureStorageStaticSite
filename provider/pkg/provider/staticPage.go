// Copyright 2016-2021, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package provider

import (
	"fmt"
	"strings"

	"github.com/pulumi/pulumi-azure-native/sdk/go/azure/cdn"
	"github.com/pulumi/pulumi-azure-native/sdk/go/azure/resources"
	"github.com/pulumi/pulumi-azure-native/sdk/go/azure/storage"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// The set of arguments for creating a StaticPage component resource.
type StaticPageArgs struct {
	// The HTML content for index.html.
	IndexContent string
}

// The StaticPage component resource.
type StaticPage struct {
	pulumi.ResourceState

	StorageAccount *storage.StorageAccount
	WebsiteUrl     pulumi.StringOutput
	CdnUrl         pulumi.StringOutput
}

// NewStaticPage creates a new StaticPage component resource.
func NewStaticPage(ctx *pulumi.Context,
	name string, args *StaticPageArgs, opts ...pulumi.ResourceOption) (*StaticPage, error) {
	if args == nil {
		args = &StaticPageArgs{}
	}

	component := &StaticPage{}

	var err error

	// Create a bucket and expose a website index document.
	resourceGroup, err := resources.NewResourceGroup(ctx, "website-rg", nil)
	if err != nil {
		return nil, err
	}

	profile, err := cdn.NewProfile(ctx, "profile", &cdn.ProfileArgs{
		ResourceGroupName: resourceGroup.Name,
		Sku: &cdn.SkuArgs{
			Name: cdn.SkuName_Standard_Microsoft,
		},
	})
	if err != nil {
		return nil, err
	}

	storageAccount, err := storage.NewStorageAccount(ctx, "sa", &storage.StorageAccountArgs{
		ResourceGroupName: resourceGroup.Name,
		Kind:              storage.KindStorageV2,
		Sku: &storage.SkuArgs{
			Name: storage.SkuName_Standard_LRS,
		},
	})
	if err != nil {
		return nil, err
	}

	endpointOrigin := storageAccount.PrimaryEndpoints.Web().ApplyT(func(endpoint string) string {
		endpoint = strings.ReplaceAll(endpoint, "https://", "")
		endpoint = strings.ReplaceAll(endpoint, "/", "")
		return endpoint
	}).(pulumi.StringOutput)

	queryStringCachingBehaviorNotSet := cdn.QueryStringCachingBehaviorNotSet
	endpoint, err := cdn.NewEndpoint(ctx, "endpoint", &cdn.EndpointArgs{
		IsHttpAllowed:    pulumi.Bool(false),
		IsHttpsAllowed:   pulumi.Bool(true),
		OriginHostHeader: endpointOrigin,
		Origins: cdn.DeepCreatedOriginArray{
			&cdn.DeepCreatedOriginArgs{
				HostName:  endpointOrigin,
				HttpsPort: pulumi.Int(443),
				Name:      pulumi.String("origin-storage-account"),
			},
		},
		ProfileName:                profile.Name,
		QueryStringCachingBehavior: &queryStringCachingBehaviorNotSet,
		ResourceGroupName:          resourceGroup.Name,
	})
	if err != nil {
		return nil, err
	}

	// Enable static website support
	staticWebsite, err := storage.NewStorageAccountStaticWebsite(ctx, "staticWebsite", &storage.StorageAccountStaticWebsiteArgs{
		AccountName:       storageAccount.Name,
		ResourceGroupName: resourceGroup.Name,
		IndexDocument:     pulumi.String("index.html"),
		Error404Document:  pulumi.String("404.html"),
	})
	if err != nil {
		return nil, err
	}

	// Upload the files
	_, err = storage.NewBlob(ctx, "index.html", &storage.BlobArgs{
		ResourceGroupName: resourceGroup.Name,
		AccountName:       storageAccount.Name,
		ContainerName:     staticWebsite.ContainerName,
		Source:            pulumi.NewStringAsset(args.IndexContent),
		ContentType:       pulumi.String("text/html"),
	})
	if err != nil {
		return nil, err
	}

	// _, err = storage.NewBlob(ctx, "404.html", &storage.BlobArgs{
	// 	ResourceGroupName: resourceGroup.Name,
	// 	AccountName:       storageAccount.Name,
	// 	ContainerName:     staticWebsite.ContainerName,
	// 	Source:            pulumi.NewFileAsset("./wwwroot/404.html"),
	// 	ContentType:       pulumi.String("text/html"),
	// })
	// if err != nil {
	// 	return nil, err
	// }

	err = ctx.RegisterComponentResource("AzureStorageStaticSite:index:StaticPage", name, component, opts...)
	if err != nil {
		return nil, err
	}

	cdnUrl := endpoint.HostName.ApplyT(func(hostName string) string {
		return fmt.Sprintf("%v%v", "https://", hostName)
	})

	component.StorageAccount = storageAccount
	component.WebsiteUrl = storageAccount.PrimaryEndpoints.Web()
	// component.CdnUrl = cdnUrl

	if err := ctx.RegisterResourceOutputs(component, pulumi.Map{
		"storageAccount": storageAccount,
		"websiteUrl":     storageAccount.PrimaryEndpoints.Web(),
		"cdnUrl":         cdnUrl,
	}); err != nil {
		return nil, err
	}

	return component, nil
}
