terraform {
  required_providers {
    pulumi = {
      version = "0.1"
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

output "version" {
  value = data.pulumi_stack_outputs.stack_outputs.version
}

output "stack_outputs" {
  value = data.pulumi_stack_outputs.stack_outputs.stack_outputs
}