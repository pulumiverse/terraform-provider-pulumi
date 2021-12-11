# Terraform Provider Pulumi

This is the `transcend-io/pulumi` provider [available on the Terraform registry](https://registry.terraform.io/providers/transcend-io/pulumi/latest).

It's goal is to allow terraform projects to consume pulumi outputs from the Pulumi Cloud via data source lookups.

## Introduction

When you use [Infrastructure as Code](https://en.wikipedia.org/wiki/Infrastructure_as_code) tools like Terraform or Pulumi, you’ll doubtlessly need to have dependencies between your stacks/modules.

This can be illustrated by a small example:

<img width="711" alt="An example backend architecture" src="https://user-images.githubusercontent.com/8922077/145655873-ba6e67e7-7c34-4006-9600-63177379f717.png">

In this example, we have a piece of code that creates a Virtual Private Cloud (VPC). The code that creates our backend application and database needs to depend on the VPCs private subnet identifiers, so that we can place our backend and database inside the VPC. Likewise, our backend code will need to depend on our database code so that it can create a connection and talk to the database.

Each tool handles this in a slightly different way: Terraform handles module dependencies via [data_source lookups](https://www.terraform.io/docs/language/state/remote-state-data.html) or [Terragrunt dependencies](https://terragrunt.gruntwork.io/docs/reference/config-blocks-and-attributes/#dependency), while Pulumi handles stack dependencies via [StackReferences](https://www.pulumi.com/docs/intro/concepts/stack/#stackreferences).

Both tools makes it easy to have multiple stacks that all reference each other, and Pulumi has [a native way to lookup Terraform outputs](https://www.pulumi.com/blog/using-terraform-remote-state-with-pulumi/), but there is no built-in Terraform support for consuming Pulumi stack outputs. That's where this tool fits in, providing a Terraform data source to consume Pulumi stack outputs.

This opens up two potential outcomes for your company:
- Some teams at your company can use Terraform while others use Pulumi.
- It becomes possible to migrate incrementally from one tool to the other, so you can have intermediary steps where some modules/stacks from one tool depend on module/stack outputs from the other.

## Usage

Usage will often follow the following flow

```terraform
terraform {
  required_providers {
    pulumi = {
      version = "0.0.2"
      source  = "hashicorp.com/transcend-io/pulumi"
    }
  }
}

provider "pulumi" {}

data "pulumi_stack_outputs" "stack_outputs" {
  organization = "transcend-io"
  project      = "airgap-telemetry-backend"
  stack        = "dev"
}

output "stack_outputs" {
  value = data.pulumi_stack_outputs.stack_outputs.stack_outputs
}
```

This code block tells Terraform to use the `0.0.2` version of the `pulumi` Terraform provider from `transcend-io`. It then declares the `provider`, where you can either supply a Pulumi cloud token directly, or via the `PULUMI_ACCESS_TOKEN` environment variable. Lastly, it looks up the `airgap-telemetry-backend/dev` stack under the `transcend-io` organization and allows Terraform to consume that stack’s outputs. You’ll want to use your own organization, project, and stack names here, of course.

And with that, your Terraform modules can now depend on your Pulumi Cloud stacks 😄

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
