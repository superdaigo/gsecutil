# gsecutil - Google Secret Manager 实用工具

Google Secret Manager 的简化命令行封装工具，可作为项目级密码管理器使用。通过直观的命令、剪贴板集成、版本控制、团队友好的配置文件和审计日志来存储、检索和管理密钥。

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

从[发布页面](https://github.com/superdaigo/gsecutil/releases)下载适用于您平台的最新二进制文件，或使用 Go 安装：

```bash
go install github.com/superdaigo/gsecutil@latest
```

### 先决条件

- 已安装并认证 [Google Cloud SDK (gcloud)](https://cloud.google.com/sdk/docs/install)
- 已启用 Secret Manager API 的 Google Cloud 项目

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

每个项目通常都有自己的配置文件，用于存储项目 ID、密钥命名约定和元数据属性。

### 1. 创建配置文件

运行交互式设置以生成配置文件。系统将提示您输入 Google Cloud 项目 ID、密钥名称前缀、默认列表属性以及可选的示例凭证。生成的文件默认保存在当前目录中，名为 `gsecutil.conf`（使用 `--home` 可保存到 `~/.config/gsecutil/gsecutil.conf`）。

```bash
gsecutil config init
```

配置文件按以下顺序搜索：
1. `--config` 标志（如果指定）
2. 当前目录：`gsecutil.conf`
3. 主目录：`~/.config/gsecutil/gsecutil.conf`

### 2. 管理密钥

```bash
# 创建密钥
gsecutil create database-password

# 获取最新版本
gsecutil get database-password

# 复制到剪贴板
gsecutil get database-password --clipboard

# 列出所有密钥
gsecutil list

# 更新密钥
gsecutil update database-password

# 删除密钥
gsecutil delete database-password
```

### 配置示例

```yaml
# 项目 ID（如果已通过环境变量或 gcloud 设置，则为可选）
project: "my-project-id"

# 用于团队组织的密钥名称前缀
prefix: "team-shared-"

# list 命令中显示的默认属性
list:
  attributes:
    - title
    - owner
    - environment

# 凭证元数据（名称为裸名 — 前缀自动添加）
credentials:
  - name: "database-password"    # 访问 "team-shared-database-password"
    title: "Production Database Password"
    environment: "production"
    owner: "backend-team"
```

> **前缀是透明的：** 配置前缀后，您在命令、配置和 CSV 文件中始终使用裸名。前缀会自动添加和去除。

有关详细的配置选项，请参阅 [docs/configuration.md](docs/configuration.md)。

## 文档

- **[配置指南](docs/configuration.md)** - 详细的配置选项和示例
- **[命令参考](docs/commands.md)** - 完整的命令文档
- **[审计日志设置](docs/audit-logging.md)** - 启用和使用审计日志
- **[故障排除指南](docs/troubleshooting.md)** - 常见问题和解决方案
- **[构建说明](BUILD.md)** - 从源代码构建
- **[开发指南](WARP.md)** - 使用 WARP AI 开发

## 许可证

本项目根据 MIT 许可证授权 - 详情请参阅 LICENSE 文件。

## 相关

- [Google Cloud SDK](https://cloud.google.com/sdk)
- [Secret Manager 文档](https://cloud.google.com/secret-manager/docs)
