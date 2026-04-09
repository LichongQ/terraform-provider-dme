# Test: Query VM list with various filters.
# examples/
# ├── main/                # 完整示例：创建 VM + 查询已有 VM
# ├── minimal/             # 最小配置：仅必填字段
# ├── status_running/      # 创建并自动开机
# ├── status_stopped/      # 创建但不开机
# ├── lifecycle/           # 状态转换分步测试
# ├── multi_vm/            # 多 VM 场景
# └── data_source_only/    # 纯数据源查询

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

# List all running VMs
data "dme_vms" "all_running" {
  status = ["running"]
}

# List VMs by name fuzzy match
data "dme_vms" "by_name" {
  name      = "test"
  page_size = 50
}

# List VMs in a specific cluster, sorted by name
data "dme_vms" "by_cluster" {
  cluster_id = "urn:sites:461D0910:clusters:173"
  sort_key   = "name"
  sort_dir   = "asc"
}

# List templates only
data "dme_vms" "templates" {
  is_template = true
}

# List Linux VMs on FusionCompute
data "dme_vms" "linux_fc" {
  os_type = ["Linux"]
  vr_type = "FUSIONCOMPUTE"
  status  = ["running", "stopped"]
}

# Outputs
output "running_vms" {
  value = [for vm in data.dme_vms.all_running.vms : "${vm.name} (${vm.ip_address})"]
}

output "running_vms_total" {
  value = data.dme_vms.all_running.total
}

output "template_names" {
  value = [for vm in data.dme_vms.templates.vms : vm.name]
}

output "linux_fc_vms" {
  value = [for vm in data.dme_vms.linux_fc.vms : { name = vm.name, status = vm.status, cpu = vm.cpu_cores, memory = vm.memory_mb }]
}
