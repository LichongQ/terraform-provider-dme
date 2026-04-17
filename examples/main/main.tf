# Example Terraform configuration for the eDME provider.
# This file demonstrates how to use the dme_vm resource and data source.

terraform {
  required_providers {
    dme = {
      source  = "registry.terraform.io/ict/dme"
      version = "0.1.0"
    }
  }
}

# Configure the eDME provider.
# Replace these values with your actual eDME endpoint and credentials.
provider "dme" {
  endpoint  = "https://10.0.0.1:26335"  # eDME management IP and port
  user_name = "northbound_user"          # eDME northbound user name
  password  = "your_password"            # eDME northbound user password
}

data "dme_vm" "template" {
  vm_id = "urn:sites:0A1010ED:vms:i-00000003"
}

# Create a VM from a template with desired running status.
resource "dme_vm" "example" {
  template_id  = data.dme_vm.template.vm_id  # Template to clone from
  name         = "terraform-test-vm"
  description  = "VM created by Terraform from eDME template"
  location     = "urn:sites:461D0910:clusters:173"      # Target cluster
  status       = "stopped"  # Create and auto-boot; change to "stopped" to keep off

  # Compute resources
  cpu_cores = 4
  memory_mb = 8192

  # OS configuration
  os_type       = data.dme_vm.template.os_type
  os_version    = data.dme_vm.template.os_version
}

# Output the created VM details.
output "vm_name" {
  description = "Name of the created VM"
  value       = dme_vm.example.name
}

output "vm_status" {
  description = "Current status of the created VM"
  value       = dme_vm.example.status
}



