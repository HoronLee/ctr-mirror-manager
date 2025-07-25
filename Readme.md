# Containerd Mirror Manager

> Containerd的镜像源管理器

# 快速开始

通过`go install`快速安装

```shell
go install github.com/HoronLee/ctr-mirror-manager@latest
```

新建一个配置文件用来声明镜像配置，下面是示例配置

```toml
# 可选：指定证书目录路径（相对路径或绝对路径）
# 如果不指定，默认使用 /etc/containerd/certs.d
certs_dir = "/etc/containerd/certs.d"

[[mirror]]
name = "docker.io"
server = "https://docker.io"

  [[mirror.host]]
  url = "https://874ca5b8e7cb6d8cabb.d.1ms.run"
  capabilities = ["pull", "resolve"]
  username = "user"
  password = "pass"

  [[mirror.host]]
  url = "https://test.d.1ms.run"
  capabilities = ["pull", "resolve"]

[[mirror]]
name = "harbor.local.com"
server = "https://harbor.local.com"

  [[mirror.host]]
  url = "https://harbor.local.com"
  capabilities = ["pull", "push", "resolve"]
```

# 使用方法

ctr-mirror-manager 是一个命令行工具，提供多个子命令来管理 Containerd 的镜像配置。

## 应用配置

使用 `apply` 子命令应用镜像配置：

```shell
$ ctr-mirror-manager apply -c mirror.toml

目录 /etc/containerd/certs.d 不存在，已自动创建
已更新: /etc/containerd/certs.d/docker.io/hosts.toml
已更新: /etc/containerd/certs.d/harbor.local.com/hosts.toml
操作成功，镜像源已更新。
```

参数说明：
- `-c, --config`: 指定配置文件路径（必需）

## 检查配置

使用 `check` 子命令检查当前镜像配置状态：

```shell
$ ctr-mirror-manager check -c mirror.toml

正在检查 /etc/containerd/certs.d 目录的配置...
备份目录 /etc/containerd/certs.d.bak 不存在
发现镜像配置：docker.io
发现镜像配置：harbor.local.com
共发现 2 个镜像配置
```

参数说明：
- `-c, --config`: 指定配置文件路径（可选，用于读取 `certs_dir`）

## 恢复备份

如果配置出现问题，可以使用 `restore` 子命令恢复之前的备份：

```shell
$ ctr-mirror-manager restore -c mirror.toml

正在从 /etc/containerd/certs.d.bak 恢复备份...
备份恢复成功
```

参数说明：
- `-c, --config`: 指定配置文件路径（可选，用于读取 `certs_dir`）