# Terraform Provider Pulumi

This is the `transcend-io/pulumi` provider available on the Terraform registry.

It's goal is to allow terraform projects to consume pulumi outputs from the Pulumi Cloud via data source lookups.

Check out the `examples` directory for usage.

## Building

Run the following command to build the provider

```shell
go build -o terraform-provider-pulumi
```

## Test sample configuration

First, build and install the provider.

```shell
make install
```

Then, run the following command to initialize the workspace and apply the sample configuration.

```shell
terraform init && terraform apply
```