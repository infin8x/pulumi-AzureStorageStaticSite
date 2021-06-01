# Azure Storage Static Site Pulumi Package

An example `StaticPage` [component resource](https://www.pulumi.com/docs/intro/concepts/resources/#components) is available in `provider/pkg/provider/staticPage.go`. This component creates a static web page hosted in an AWS S3 Bucket. There is nothing special about `StaticPage` -- it is a typical component resource written in Go.

The component provider makes component resources available to other languages. The implementation is in `provider/pkg/provider/provider.go`. Each component resource in the provider must have an implementation in the `Construct` function to create an instance of the requested component resource and return its `URN` and state (outputs). There is an initial implementation that demonstrates an implementation of `Construct` for the example `StaticPage` component.

A code generator is available which generates SDKs in TypeScript, Python, Go and .NET which are also checked in to the `sdk` folder. The SDKs are generated from a schema in `schema.json`. This file should be kept aligned with the component resources supported by the component provider implementation.

An example of using the `StaticPage` component in TypeScript is in `examples/simple`.

Note that the generated provider plugin (`pulumi-resource-AzureStorageStaticSite`) must be on your `PATH` to be used by Pulumi deployments. If creating a provider for distribution to other users, you should ensure they install this plugin to their `PATH`.

## Prerequisites

- Go 1.15
- Pulumi CLI
- Node.js (to build the Node.js SDK)
- Yarn (to build the Node.js SDK)
- Python 3.6+ (to build the Python SDK)
- .NET Core SDK (to build the .NET SDK)

## Build and Test

```bash
# Build and install the provider (plugin copied to $GOPATH/bin)
make install_provider

# Regenerate SDKs
make generate

# Test Node.js SDK
$ make install_nodejs_sdk
$ cd examples/simple
$ yarn install
$ yarn link @pulumi/AzureStorageStaticSite
$ pulumi stack init test
$ pulumi config set aws:region us-east-1
$ pulumi up
```
