# terraform-provider-dme

## 一. 前提条件
- 已安装go语言环境
- 已安装Terraform

## 二. 获取对接DME环境信息
### step1： 获取对接DME环境的地址，如：https://{ip}:26335，其中ip为DME运维面浮动IP，26335固定为DME北向API访问端口
### step2： 获取对接DME环境的北向用户和密码：
  - 登录DME运维面，登录地址：https://{ip}:31943，其中ip为DME运维面浮动IP，31943固定为DME运维面访问端口
  - 在主菜单选择“设置->安全管理->用户管理”，在左侧菜单栏单击“用户”
  - 单击“创建”
  - 根据提示填写相应的用户基本信息，“类型”选择“三方系统接入”，并单击“下一步”
  - 选择所属角色为“北向用户组”，同时选择对应的业务角色。

## 三. windows环境下provider对接方法

### step1: 下载代码，编译：go build，生成terraform-provider-dme.exe
### step2: 创建目录：C:\DMEProvider，并将terraform-provider-dme.exe拷贝到该目录
### step3：打开文件资源管理器，在地址栏输入“%APPDATA%”并按回车，进入该目录。
然后在该目录新建一个名为“terraform.rc”的文件，并将以下内容粘贴进去
```
provider_installation {
dev_overrides {
"registry.terraform.io/hashicorp/local-storage" = "C:/DMEProvider"
"registry.terraform.io/ict/dme" = "C:/DMEProvider"
}
direct {}
}
```
### step4：切换到examples/main目录，执行terraform plan命令，查看资源执行计划
### step5：执行terraform apply命令，将执行计划应用到真实环境

## 四. provider配置指导
### DME provider配置
```
terraform {
  required_providers {
    dme = {
      source  = "registry.terraform.io/ict/dme"
      version = "0.1.0"
    }
  }
}

# Configure the DME provider.
provider "dme" {
  endpoint  = "https://{ip}:26335"  # DME management IP and port
  user_name = "xxx-north"          # DME northbound user name
  password  = "********"            # DME northbound user password
}
```
### DME datasource示例
```
data "dme_vm" "template" {
  #name = "TPL-UBU22"
  vm_id = "urn:sites:0A1010ED:vms:i-00000003"
}
```
### DME vm resource示例
```
resource "dme_vm" "example" {
  template_id  = data.dme_vm.template.vm_id  # Template to clone from
  name         = "terraform-test-vm01"
  description  = "VM created by Terraform from DME template"
  location     = "urn:sites:0A1010ED:clusters:133"      # Target cluster
  is_binding_host = false
  status       = "stopped"  # Create and auto-boot; change to "stopped" to keep off

  # Compute resources
  cpu_cores = 4
  memory_mb = 8192

  # OS configuration
  os_type       = "Linux"
  os_version    = "10112"
}

output "vm_name" {
  description = "Name of the created VM"
  value       = dme_vm.example.name
}

output "vm_status" {
  description = "Current status of the created VM"
  value       = dme_vm.example.status
}
```
## 五. provider支持场景

### 1. 创建虚拟机 
```
terraform {
  required_providers {
    dme = {
      source  = "registry.terraform.io/ict/dme"
      version = "0.1.0"
    }
  }
}

# Configure the DME provider.
provider "dme" {
  endpoint  = "https://{ip}:26335"  # DME management IP and port
  user_name = "xxx-north"          # DME northbound user name
  password  = "********"            # DME northbound user password
}

data "dme_vm" "template" {
  #name = "TPL-UBU22"
  vm_id = "urn:sites:0A1010ED:vms:i-00000003"
}

resource "dme_vm" "example" {
  template_id  = data.dme_vm.template.vm_id  # Template to clone from
  name         = "terraform-test-vm01"
  description  = "VM created by Terraform from DME template"
  location     = "urn:sites:0A1010ED:clusters:133"      # Target cluster
  is_binding_host = false
  status       = "stopped"  # Create and auto-boot; change to "stopped" to keep off

  # Compute resources
  cpu_cores = 4
  memory_mb = 8192

  # OS configuration
  os_type       = "Linux"
  os_version    = "10112"
}

output "vm_name" {
  description = "Name of the created VM"
  value       = dme_vm.example.name
}

output "vm_status" {
  description = "Current status of the created VM"
  value       = dme_vm.example.status
}
```

### 2. 查询虚拟机
参考创建虚拟机的datasource
### 3. 修改虚拟机基本属性
参考创建虚拟机，修改虚拟机resource的基本属性，如：修改name属性为"terraform-test-vm"
### 4. 启动虚拟机
参考创建虚拟机，修改虚拟机resource的status字段，如修改status属性为"stopped” -> "running"
### 5. 重启虚拟机
参考创建虚拟机，修改虚拟机resource的status字段，如修改status属性为"running” -> "restarted"
### 6. 停止虚拟机
参考创建虚拟机，修改虚拟机resource的status字段，如修改status属性为"restarted” -> "stopped"
### 7. 删除虚拟机
参考创建虚拟机，删除虚拟机resource
