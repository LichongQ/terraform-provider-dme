# Test: Create a VM with status=running (is_auto_boot=true).
# The VM will be powered on immediately after creation.

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

resource "dme_vm" "auto_start" {
  template_id = "urn:sites:461D0910:vms:i-000001BB"
  name        = "test-auto-start-vm"
  location    = "urn:sites:461D0910:clusters:173"
  status      = "running"  # is_auto_boot = true

  cpu_cores = 2
  memory_mb = 4096

  os_type       = "Linux"
  os_version    = "1719"
  guest_os_name = "CentOS 7.9 64bit"

  hostname = "auto-start-vm"
  password = "SecurePass123!"

  port_group_id = "urn:sites:45E707F1:dvswitchs:1:portgroups:1"

  disks {
    sequence_num = 1
    datastore_id = "urn:sites:461D0910:datastores:195"
    capacity_gb  = 50
    pci_type     = "VIRTIO"
  }
}

output "vm_status" {
  value = dme_vm.auto_start.status
}

output "vm_ip" {
  value = dme_vm.auto_start.ip_address_actual
}
