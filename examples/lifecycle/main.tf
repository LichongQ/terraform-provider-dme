# Test: VM status lifecycle - demonstrates state transitions.
#
# Usage:
#   1. terraform apply  -> creates VM (stopped)
#   2. Uncomment "Phase 2", change status to "running", apply -> powers on (force)
#   3. Uncomment "Phase 3", change status to "restarted", apply -> reboots (force)
#   4. Uncomment "Phase 4", change status to "stopped", apply -> powers off (force)
#   5. terraform destroy -> deletes VM

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

# ---- Phase 1: Create VM, stopped ----
resource "dme_vm" "lifecycle" {
  template_id = "urn:sites:461D0910:vms:i-000001BB"
  name        = "test-lifecycle-vm"
  location    = "urn:sites:461D0910:clusters:173"
  # status not set -> is_auto_boot = false

  cpu_cores = 2
  memory_mb = 4096

  os_type       = "Linux"
  os_version    = "1719"
  guest_os_name = "CentOS 7.9 64bit"

  hostname = "lifecycle-vm"
  password = "LifecyclePass123!"

  port_group_id = "urn:sites:45E707F1:dvswitchs:1:portgroups:1"
}

# ---- Phase 2: Power on ----
# resource "dme_vm" "lifecycle" {
#   template_id = "urn:sites:461D0910:vms:i-000001BB"
#   name        = "test-lifecycle-vm"
#   location    = "urn:sites:461D0910:clusters:173"
#   status      = "running"   # triggers batch-power-on with force mode
#
#   cpu_cores = 2
#   memory_mb = 4096
#
#   os_type       = "Linux"
#   os_version    = "1719"
#   guest_os_name = "CentOS 7.9 64bit"
#
#   hostname = "lifecycle-vm"
#   password = "LifecyclePass123!"
#
#   port_group_id = "urn:sites:45E707F1:dvswitchs:1:portgroups:1"
# }

# ---- Phase 3: Reboot ----
# resource "dme_vm" "lifecycle" {
#   template_id = "urn:sites:461D0910:vms:i-000001BB"
#   name        = "test-lifecycle-vm"
#   location    = "urn:sites:461D0910:clusters:173"
#   status      = "restarted"  # triggers batch-reboot with force mode
#
#   cpu_cores = 2
#   memory_mb = 4096
#
#   os_type       = "Linux"
#   os_version    = "1719"
#   guest_os_name = "CentOS 7.9 64bit"
#
#   hostname = "lifecycle-vm"
#   password = "LifecyclePass123!"
#
#   port_group_id = "urn:sites:45E707F1:dvswitchs:1:portgroups:1"
# }

# ---- Phase 4: Power off ----
# resource "dme_vm" "lifecycle" {
#   template_id = "urn:sites:461D0910:vms:i-000001BB"
#   name        = "test-lifecycle-vm"
#   location    = "urn:sites:461D0910:clusters:173"
#   status      = "stopped"   # triggers batch-power-off with force mode
#
#   cpu_cores = 2
#   memory_mb = 4096
#
#   os_type       = "Linux"
#   os_version    = "1719"
#   guest_os_name = "CentOS 7.9 64bit"
#
#   hostname = "lifecycle-vm"
#   password = "LifecyclePass123!"
#
#   port_group_id = "urn:sites:45E707F1:dvswitchs:1:portgroups:1"
# }

output "vm_id" {
  value = dme_vm.lifecycle.vm_id
}

output "vm_status" {
  value = dme_vm.lifecycle.status
}

output "vm_ip" {
  value = dme_vm.lifecycle.ip_address_actual
}

output "vm_host" {
  value = dme_vm.lifecycle.host_name
}

output "vm_cluster" {
  value = dme_vm.lifecycle.cluster_name
}
