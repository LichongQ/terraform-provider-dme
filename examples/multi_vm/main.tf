# Test: Multiple VMs with different configurations in one run.

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

# VM 1: Web server - running, with full network config
resource "dme_vm" "web" {
  template_id = "urn:sites:461D0910:vms:i-000001BB"
  name        = "test-web-server"
  location    = "urn:sites:461D0910:clusters:173"
  status      = "running"

  cpu_cores = 4
  memory_mb = 8192

  os_type       = "Linux"
  os_version    = "1719"
  guest_os_name = "CentOS 7.9 64bit"

  hostname = "web-01"
  password = "WebPass123!"

  port_group_id = "urn:sites:45E707F1:dvswitchs:1:portgroups:1"
  ip_address    = "192.168.1.10"
  netmask       = "255.255.255.0"
  gateway       = "192.168.1.1"

  disks {
    sequence_num = 1
    datastore_id = "urn:sites:461D0910:datastores:195"
    capacity_gb  = 100
    pci_type     = "VIRTIO"
    is_thin      = true
  }
  disks {
    sequence_num = 2
    datastore_id = "urn:sites:461D0910:datastores:195"
    capacity_gb  = 200
    pci_type     = "VIRTIO"
    indip_disk   = true
    persistent_disk = true
  }
}

# VM 2: DB server - running, more memory, separate datastore
resource "dme_vm" "db" {
  template_id = "urn:sites:461D0910:vms:i-000001BB"
  name        = "test-db-server"
  location    = "urn:sites:461D0910:clusters:173"
  status      = "running"

  cpu_cores = 8
  memory_mb = 32768

  os_type       = "Linux"
  os_version    = "1719"
  guest_os_name = "CentOS 7.9 64bit"

  hostname = "db-01"
  password = "DbPass123!"

  port_group_id = "urn:sites:45E707F1:dvswitchs:1:portgroups:2"
  ip_address    = "192.168.2.10"
  netmask       = "255.255.255.0"
  gateway       = "192.168.2.1"

  disks {
    sequence_num = 1
    datastore_id = "urn:sites:461D0910:datastores:196"
    capacity_gb  = 500
    pci_type     = "SCSI"
    indip_disk   = true
    persistent_disk = true
    is_thin      = false
  }
}

# VM 3: Test node - stopped, minimal config
resource "dme_vm" "test_node" {
  template_id = "urn:sites:461D0910:vms:i-000001BB"
  name        = "test-minimal-node"
  location    = "urn:sites:461D0910:clusters:173"
  # no status -> stopped

  cpu_cores = 1
  memory_mb = 1024
}

# Query an existing VM via data source
data "dme_vm" "existing" {
  vm_id = "urn:sites:461D0910:vms:i-000002CC"
}

# Outputs
output "web_vm_id" {
  value = dme_vm.web.vm_id
}

output "web_ip" {
  value = dme_vm.web.ip_address_actual
}

output "db_vm_id" {
  value = dme_vm.db.vm_id
}

output "db_ip" {
  value = dme_vm.db.ip_address_actual
}

output "test_node_status" {
  value = dme_vm.test_node.status
}

output "existing_vm_name" {
  value = data.dme_vm.existing.name
}
