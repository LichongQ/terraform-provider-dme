# Test: Data source only - query existing VMs without creating new ones.

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

data "dme_vm" "vm1" {
  vm_id = "urn:sites:461D0910:vms:i-000001AA"
}

data "dme_vm" "vm2" {
  vm_id = "urn:sites:461D0910:vms:i-000002BB"
}

output "vm1_name" {
  value = data.dme_vm.vm1.name
}

output "vm1_status" {
  value = data.dme_vm.vm1.status
}

output "vm1_ip" {
  value = data.dme_vm.vm1.ip_address
}

output "vm1_cpu" {
  value = data.dme_vm.vm1.cpu_cores
}

output "vm1_memory" {
  value = data.dme_vm.vm1.memory_mb
}

output "vm2_name" {
  value = data.dme_vm.vm2.name
}

output "vm2_status" {
  value = data.dme_vm.vm2.status
}
