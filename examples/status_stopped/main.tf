# Test: Create a VM with status=stopped (is_auto_boot=false).
# The VM will NOT be powered on after creation.
# To start it later, change status to "running" and run terraform apply.

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

resource "dme_vm" "stopped" {
  template_id = "urn:sites:461D0910:vms:i-000001BB"
  name        = "test-stopped-vm"
  location    = "urn:sites:461D0910:clusters:173"
  # status not set -> defaults to stopped (is_auto_boot = false)

  cpu_cores = 2
  memory_mb = 4096
}

# ---- Phase 2: uncomment below to power on the VM ----
# resource "dme_vm" "stopped" {
#   template_id = "urn:sites:461D0910:vms:i-000001BB"
#   name        = "test-stopped-vm"
#   location    = "urn:sites:461D0910:clusters:173"
#   status      = "running"  # <-- triggers batch-power-on (force mode)
#
#   cpu_cores = 2
#   memory_mb = 4096
# }

output "vm_status" {
  value = dme_vm.stopped.status
}

output "vm_id" {
  value = dme_vm.stopped.vm_id
}
