# Example Terraform configuration for the eDME provider.
# This file demonstrates how to use the dme_vm resource and data source.

terraform {
  required_providers {
    dme = {
      source  = "registry.terraform.io/ict/dme"
      version = "0.1.0"
    }
  }
  required_version = ">= 1.5.0"
}

# Configure the eDME provider.
# Replace these values with your actual eDME endpoint and credentials.
provider "dme" {
  endpoint  = "https://10.0.0.1:26335"  # eDME management IP and port
  user_name = "northbound_user"          # eDME northbound user name
  password  = "your_password"            # eDME northbound user password
}

# Create a VM from a template with desired running status.
resource "dme_vm" "example" {
  template_id  = "urn:sites:461D0910:vms:i-000001BB"  # Template to clone from
  name         = "terraform-test-vm"
  description  = "VM created by Terraform from eDME template"
  location     = "urn:sites:461D0910:clusters:173"      # Target cluster
  is_binding_host = false
  status       = "running"  # Create and auto-boot; change to "stopped" to keep off

  # Compute resources
  cpu_cores = 4
  memory_mb = 8192

  # OS configuration
  os_type       = "Linux"
  os_version    = "1719"
  guest_os_name = "CentOS 7.9 64bit"
/*
  # Guest OS customization
  hostname = "test-vm-01"
  password = "YourSecurePassword123!"

  # Network configuration
  port_group_id = "urn:sites:45E707F1:dvswitchs:1:portgroups:1"
  ip_address    = "192.168.1.100"
  netmask       = "255.255.255.0"
  gateway       = "192.168.1.1"

  # Disk configuration
  disks {
    sequence_num = 1
    datastore_id = "urn:sites:461D0910:datastores:195"
    capacity_gb  = 100
    pci_type     = "VIRTIO"
    indip_disk   = false
    persistent_disk = true
    is_thin      = true
  }*/
}

/*
# Query an existing VM using the data source.
data "dme_vm" "existing" {
  vm_id = "urn:sites:461D0910:vms:i-000001BB"
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

output "vm_actual_id" {
  description = "Actual eDME VM ID"
  value       = dme_vm.example.vm_id
}

output "vm_ip" {
  description = "Actual IP address of the VM"
  value       = dme_vm.example.ip_address_actual
}

# Output the queried VM details from data source.
output "existing_vm_name" {
  description = "Name of the existing VM"
  value       = data.dme_vm.existing.name
}

output "existing_vm_status" {
  description = "Status of the existing VM"
  value       = data.dme_vm.existing.status
}
*/