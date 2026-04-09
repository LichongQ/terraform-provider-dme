# Test: Minimal VM creation with only required fields.

terraform {
  required_providers {
    dme = {
      source  = "local/dme"
      version = "0.1.0"
    }
  }
  required_version = ">= 1.5.0"
}

provider "dme" {
  endpoint  = "https://10.0.0.1:26335"
  user_name = "northbound_user"
  password  = "your_password"
}

# Minimal: only required attributes, VM will be created stopped.
resource "dme_vm" "minimal" {
  template_id = "urn:sites:461D0910:vms:i-000001BB"
  name        = "test-minimal-vm"
  location    = "urn:sites:461D0910:clusters:173"
}

output "vm_id" {
  value = dme_vm.minimal.vm_id
}

output "vm_status" {
  value = dme_vm.minimal.status
}
