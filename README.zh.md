# gsecutil - Google Secret Manager 实用工具

> **翻译说明**: 这个README文件是机器翻译的。最准确和最新的信息，请参考英语版[README.md](README.md)。
>
> **🆕 新功能**: v1.1.1添加了自动版本管理功能，以保持在Google Cloud免费层（6个活动版本）内。详情请参阅英语 README。

🚀 **v1.1.0** - 带有配置文件支持的Google Secret Manager简化命令行包装器。`gsecutil`为常见的密钥操作提供便利命令，使小团队能够更轻松地使用Google Cloud的Secret Manager管理密码和凭据，而无需专门的密码管理工具。

**v1.1.0新功能**: YAML配置文件支持、前缀功能、带有团队自定义元数据的增强列表和描述命令。

## ✨ 功能特性

### 🔐 **完整的密钥管理**
- **CRUD操作**: 使用简化命令创建、读取、更新、删除密钥
- **版本管理**: 访问任何版本，查看版本历史和元数据
- **跨平台**支持（Linux、macOS、Windows，支持ARM64）
- **剪贴板集成** - 直接将密钥值复制到剪贴板
- **交互式和文件输入** - 安全提示或基于文件的密钥加载

### 🛡️ **高级访问管理**
*(在v1.0.0中引入)*
- **完整的IAM策略分析** - 查看谁可以在任何级别访问密钥
- **多级权限检查** - 密钥级别和项目级别的访问分析
- **IAM条件感知** - 完全支持带CEL表达式的条件访问策略
- **主体管理** - 为用户、组和服务账户授予/撤销访问权限
- **项目范围分析** - 识别提供Secret Manager访问权限的编辑者/所有者角色

### 📊 **审计与合规**
- **综合审计日志** - 跟踪谁访问了密钥、何时访问以及执行了什么操作
- **基于主体的过滤** - 查看特定用户/组可访问的所有密钥
- **灵活过滤** - 按密钥、主体、操作类型、时间范围过滤
- **条件评估** - 了解条件访问何时适用

### 🎯 **生产就绪**
- **一致的API** - 所有命令统一的参数命名
- **企业功能** - IAM条件、项目级别分析、合规审计
- **健壮的错误处理** - 优雅处理权限不足和网络问题
- **灵活的输出** - JSON、YAML、表格格式，支持丰富的格式化

## 前提条件

- 已安装并认证的[Google Cloud SDK (gcloud)](https://cloud.google.com/sdk/docs/install)
- 启用了Secret Manager API的Google Cloud项目
- 适当的Secret Manager操作IAM权限

## 安装

### 预构建二进制文件

从[发布页面](https://github.com/superdaigo/gsecutil/releases)下载适合您平台的最新版本：

| 平台 | 架构 | 下载 |
|----------|--------------|----------|
| Linux | x64 | `gsecutil-linux-amd64-v{version}` |
| Linux | ARM64 | `gsecutil-linux-arm64-v{version}` |
| macOS | Intel | `gsecutil-darwin-amd64-v{version}` |
| macOS | Apple Silicon | `gsecutil-darwin-arm64-v{version}` |
| Windows | x64 | `gsecutil-windows-amd64-v{version}.exe` |

**下载后:** 重命名二进制文件以便一致使用：

```bash
# Linux/macOS 示例:
mv gsecutil-linux-amd64-v1.1.0 gsecutil
chmod +x gsecutil

# Windows 示例 (PowerShell/Command Prompt):
ren gsecutil-windows-amd64-v1.1.0.exe gsecutil.exe
```

这样您就可以不受版本影响地一致使用`gsecutil`。

### 使用Go安装

```bash
go install github.com/superdaigo/gsecutil@latest
```

### 从源码构建

有关综合构建说明，请参见[BUILD.md](BUILD.md)。

**快速构建:**
```bash
git clone https://github.com/superdaigo/gsecutil.git
cd gsecutil

# 为当前平台构建
make build
# 或
./build.sh          # Linux/macOS
.\\build.ps1         # Windows

# 为所有平台构建
make build-all
# 或
./build.sh all      # Linux/macOS
.\\build.ps1 all     # Windows
```

## 使用方法

### 全局选项

- `-p, --project`: Google Cloud项目ID（也可通过`GOOGLE_CLOUD_PROJECT`环境变量设置）

### 命令

#### Get Secret（获取密钥）

从Google Secret Manager获取密钥值。默认情况下获取最新版本，但您可以指定任何版本：

```bash
# 获取密钥的最新版本
gsecutil get my-secret

# 获取特定版本（对回滚、调试或访问历史值很有用）
gsecutil get my-secret --version 1
gsecutil get my-secret -v 3

# 获取密钥并复制到剪贴板
gsecutil get my-secret --clipboard

# 使用剪贴板获取特定版本
gsecutil get my-secret --version 2 --clipboard

# 获取带版本元数据的密钥（版本、创建时间、状态）
gsecutil get my-secret --show-metadata
gsecutil get my-secret -v 1 --show-metadata    # 带元数据的旧版本

# 从特定项目获取密钥
gsecutil get my-secret --project my-gcp-project
```

**版本支持:**
- 🔄 **最新版本**: 未指定`--version`时的默认行为
- 📅 **历史版本**: 按编号访问任何以前的版本（例如：`--version 1`、`--version 2`）
- 🔍 **版本元数据**: 使用`--show-metadata`查看版本详细信息（创建时间、状态、ETag）
- 📋 **剪贴板支持**: 使用`--clipboard`处理任何版本

## 配置

### 环境变量

- `GOOGLE_CLOUD_PROJECT`: 默认项目ID（被`--project`标志覆盖）

### 认证

`gsecutil`使用与`gcloud`相同的认证。确保您已认证：

```bash
# 使用gcloud认证
gcloud auth login

# 设置默认项目
gcloud config set project YOUR_PROJECT_ID

# 对于服务账户（在CI/CD中）
gcloud auth activate-service-account --key-file=service-account.json
```

### Shell 自动补全

`gsecutil`支持bash、zsh、fish和PowerShell的shell自动补全。这使得命令、标志和选项的tab补全成为可能，使CLI更加用户友好。

#### 设置说明

**Bash:**
```bash
# 临时（仅当前会话）
source <(gsecutil completion bash)

# 永久安装（需要bash-completion包）
# 系统全局（需要sudo）
sudo gsecutil completion bash > /etc/bash_completion.d/gsecutil

# 用户本地安装
gsecutil completion bash > ~/.local/share/bash-completion/completions/gsecutil

# 或添加到~/.bashrc以自动加载
echo 'source <(gsecutil completion bash)' >> ~/.bashrc
```

**Zsh:**
```bash
# 临时（仅当前会话）
source <(gsecutil completion zsh)

# 永久安装
gsecutil completion zsh > "${fpath[1]}/_gsecutil"

# 或添加到~/.zshrc以自动加载
echo 'source <(gsecutil completion zsh)' >> ~/.zshrc
```

**Fish:**
```bash
# 临时（仅当前会话）
gsecutil completion fish | source

# 永久安装
gsecutil completion fish > ~/.config/fish/completions/gsecutil.fish
```

**PowerShell:**
```powershell
# 添加到PowerShell配置文件
gsecutil completion powershell | Out-String | Invoke-Expression

# 或保存到配置文件以自动加载
gsecutil completion powershell >> $PROFILE
```

#### 功能特性

安装后，shell补全提供：
- **命令补全**: Tab补全`gsecutil`子命令（`get`、`create`、`list`等）
- **标志补全**: Tab补全标志，如`--project`、`--version`、`--clipboard`
- **智能建议**: 基于当前命令的上下文相关补全
- **帮助文本**: 命令和标志的简要说明（在支持的情况下）

#### 使用示例

```bash
# 输入并按Tab查看可用命令
gsecutil <Tab>
# 显示: access, auditlog, completion, create, delete, describe, get, help, list, update

# 输入部分命令并按Tab补全
gsecutil des<Tab>
# 补全为: gsecutil describe

# Tab补全也适用于标志
gsecutil get my-secret --<Tab>
# 显示: --clipboard, --project, --show-metadata, --version
```

**注意**: 您可能需要重启您的shell或源加载您的shell配置文件，以使补全生效。

## 安全性和最佳实践

### 安全特性

- **无持久存储**: 密钥值永远不会被`gsecutil`记录或存储
- **安全输入**: 交互式提示使用隐藏密码输入
- **OS原生剪贴板**: 剪贴板操作使用安全的OS原生API
- **gcloud委托**: 所有操作委托给已认证的`gcloud` CLI

### 最佳实践

- **谨慎使用`--force`**: 在自动化环境中使用`--force`之前始终检查
- **环境变量**: 设置`GOOGLE_CLOUD_PROJECT`以避免重复的`--project`标志
- **版本控制**: 在生产环境中使用特定的密钥版本（`--version N`）
- **定期审计**: 使用`gsecutil auditlog secret-name`监控密钥访问
- **密钥轮换**: 使用`gsecutil update`进行定期密钥轮换

## 故障排除

### 常见问题

1. **"gcloud command not found"**
   - 确保已安装Google Cloud SDK且`gcloud`在您的PATH中

2. **认证错误**
   - 运行`gcloud auth login`进行认证
   - 验证项目访问：`gcloud config get-value project`

3. **权限被拒绝错误**
   - 确保您的账户具有必要的IAM角色：
     - `roles/secretmanager.admin`（用于所有操作）
     - `roles/secretmanager.secretAccessor`（用于读取操作）
     - `roles/secretmanager.secretVersionManager`（用于创建/更新操作）

4. **剪贴板不工作**
   - 确保您有图形环境（对于Linux）
   - 在无头服务器上，剪贴板操作可能会优雅地失败

### 调试模式

通过设置以下内容为gcloud命令添加详细输出：

```bash
export CLOUDSDK_CORE_VERBOSITY=debug
```

## 文档

- **[BUILD.md](BUILD.md)** - 所有平台的综合构建说明
- **[CONTRIBUTING.md](CONTRIBUTING.md)** - 贡献指南和开发工作流程
- **[WARP.md](WARP.md)** - WARP AI终端集成的开发指导
- **README.md** - 此文件，使用方法和概述

## 贡献

欢迎贡献！有关如何为此项目做出贡献的详细指南，包括开发环境和预提交钩子的设置说明，请参见[CONTRIBUTING.md](CONTRIBUTING.md)。

## 许可证

此项目在MIT许可证下授权 - 详细信息请参见LICENSE文件。

## 相关项目

- [Google Cloud SDK](https://cloud.google.com/sdk)
- [Secret Manager Documentation](https://cloud.google.com/secret-manager/docs)
- [Cobra CLI Framework](https://github.com/spf13/cobra)
