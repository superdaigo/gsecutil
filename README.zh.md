# gsecutil - Google Secret Manager 实用工具

🚀 一个简化的 Google Secret Manager 命令行包装器，支持配置文件和团队友好功能。

## 🌍 语言版本

- **English** - [README.md](README.md)
- **日本語** - [README.ja.md](README.ja.md)
- **中文** - [README.zh.md](README.zh.md)（当前）
- **Español** - [README.es.md](README.es.md)
- **हिंदी** - [README.hi.md](README.hi.md)
- **Português** - [README.pt.md](README.pt.md)

> **注意**：所有非英语版本均为机器翻译。有关最准确的信息，请参阅英文版本。

## 快速开始

### 安装

从[发布页面](https://github.com/superdaigo/gsecutil/releases)下载适合您平台的最新二进制文件：

```bash
# macOS Apple Silicon
curl -L https://github.com/superdaigo/gsecutil/releases/latest/download/gsecutil-darwin-arm64 -o gsecutil
chmod +x gsecutil
sudo mv gsecutil /usr/local/bin/

# macOS Intel
curl -L https://github.com/superdaigo/gsecutil/releases/latest/download/gsecutil-darwin-amd64 -o gsecutil
chmod +x gsecutil
sudo mv gsecutil /usr/local/bin/

# Linux
curl -L https://github.com/superdaigo/gsecutil/releases/latest/download/gsecutil-linux-amd64 -o gsecutil
chmod +x gsecutil
sudo mv gsecutil /usr/local/bin/

# Windows (PowerShell)
Invoke-WebRequest -Uri "https://github.com/superdaigo/gsecutil/releases/latest/download/gsecutil-windows-amd64.exe" -OutFile "gsecutil.exe"
# Move to a directory in your PATH, e.g., C:\Windows\System32
Move-Item gsecutil.exe C:\Windows\System32\gsecutil.exe
```

或使用 Go 安装：
```bash
go install github.com/superdaigo/gsecutil@latest
```

### 先决条件

- 已安装并认证 [Google Cloud SDK (gcloud)](https://cloud.google.com/sdk/docs/install)
- 启用了 Secret Manager API 的 Google Cloud 项目

### 认证

```bash
# 使用 gcloud 进行认证
gcloud auth login

# 设置默认项目
gcloud config set project YOUR_PROJECT_ID

# 或设置环境变量
export GSECUTIL_PROJECT=YOUR_PROJECT_ID
```

## 基本用法

### 创建密钥
```bash
# 交互式输入
gsecutil create database-password

# 从命令行
gsecutil create api-key -d "sk-1234567890"

# 从文件
gsecutil create config --data-file ./config.json
```

### 获取密钥
```bash
# 获取最新版本
gsecutil get database-password

# 复制到剪贴板
gsecutil get api-key --clipboard

# 获取特定版本
gsecutil get api-key --version 2
```

### 列出密钥
```bash
# 列出所有密钥
gsecutil list

# 按标签过滤
gsecutil list --filter "labels.env=prod"
```

### 更新密钥
```bash
# 交互式输入
gsecutil update database-password

# 从命令行
gsecutil update api-key -d "new-secret-value"
```

### 删除密钥
```bash
gsecutil delete old-secret
```

## 配置

在 `~/.config/gsecutil/gsecutil.conf` 创建配置文件：

```yaml
# 项目 ID（如果通过环境变量或 gcloud 设置则为可选）
project: "my-project-id"

# 用于团队组织的密钥名称前缀
prefix: "team-shared-"

# list 命令中显示的默认属性
list:
  attributes:
    - title
    - owner
    - environment

# 凭据元数据
credentials:
  - name: "database-password"
    title: "Production Database Password"
    environment: "production"
    owner: "backend-team"
```

交互式生成配置：
```bash
gsecutil config init
```

有关详细的配置选项，请参阅 [docs/configuration.md](docs/configuration.md)。

## 主要功能

- ✅ **简单的 CRUD 操作** - 用于管理密钥的直观命令
- ✅ **剪贴板集成** - 直接将密钥复制到剪贴板
- ✅ **版本管理** - 访问特定版本并管理版本生命周期
- ✅ **配置文件支持** - 团队友好的元数据和组织
- ✅ **访问管理** - 基本的 IAM 策略管理
- ✅ **审计日志** - 查看谁在何时访问了密钥
- ✅ **多种输入方法** - 交互式、内联或基于文件
- ✅ **跨平台** - Linux、macOS、Windows（amd64/arm64）

## 文档

- **[配置指南](docs/configuration.md)** - 详细的配置选项和示例
- **[命令参考](docs/commands.md)** - 完整的命令文档
- **[审计日志设置](docs/audit-logging.md)** - 启用和使用审计日志
- **[故障排除指南](docs/troubleshooting.md)** - 常见问题和解决方案
- **[构建说明](BUILD.md)** - 从源代码构建
- **[开发指南](WARP.md)** - 使用 WARP AI 进行开发

## 常用命令

```bash
# 显示密钥详情
gsecutil describe my-secret

# 显示版本历史
gsecutil describe my-secret --show-versions

# 查看审计日志
gsecutil auditlog my-secret

# 管理访问
gsecutil access list my-secret
gsecutil access grant my-secret --principal user:alice@example.com

# 验证配置
gsecutil config validate

# 显示配置
gsecutil config show
```

## 许可证

本项目根据 MIT 许可证授权 - 有关详细信息，请参阅 LICENSE 文件。

## 相关链接

- [Google Cloud SDK](https://cloud.google.com/sdk)
- [Secret Manager 文档](https://cloud.google.com/secret-manager/docs)
