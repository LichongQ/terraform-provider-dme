# terraform-provider-dme

## 一. 前提条件
- 已安装go语言环境
- 已安装Terraform

## 二. 获取对接DME环境信息
### step1： 获取对接DME环境的地址，如：https://192.168.1.2:26335，其中192.168.1.2为DME运维面浮动IP，26335固定为DME北向API访问端口
### step2： 获取对接DME环境的北向用户和密码：
  - 登录DME运维面，登录地址：https://192.168.1.2:31943，其中192.168.1.2为DME运维面浮动IP，31943固定为DME运维面访问端口
  - 在主菜单选择“设置->安全管理->用户管理”，在左侧菜单栏单击“用户”
  - 单击“创建”
  - 根据提示填写相应的用户基本信息，“类型”选择“三方系统接入”，并单击“下一步”
  - 选择所属角色为“北向用户组”，同时选择对应的业务角色。

## 三. windows环境下provider对接方法

### step1: 下载代码，编译：go build，生成terraform-provider-dme.exe
### step2: 创建目录：C:\DMEProvider，并将terraform-provider-dme.exe拷贝到该目录
### step3：打开文件资源管理器，在地址栏输入“%APPDATA%”并按回车，进入该目录。
然后在该目录新建一个名为“terraform.rc”的文件，并将以下内容粘贴进去
``` provider_installation {
dev_overrides {
"registry.terraform.io/hashicorp/local-storage" = "C:/DMEProvider"
"registry.terraform.io/ict/dme" = "C:/DMEProvider"
}
direct {}
}
```
### step4：切换到examples目录，执行terraform plan命令
### step5：执行terraform apply命令
### step6：执行terraform destroy命令
