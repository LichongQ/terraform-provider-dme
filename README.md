# terraform-provider-dme

## 一. 前提条件
- 已安装go语言环境
- 已安装Terraform

## 二. provider对接方法

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
