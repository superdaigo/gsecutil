# gsecutil - Google Secret Manager ユーティリティ

🚀 設定ファイルサポートとチームフレンドリーな機能を備えた、Google Secret Manager のシンプルなコマンドラインラッパー。

## 🌍 言語バージョン

- **English** - [README.md](README.md)
- **日本語** - [README.ja.md](README.ja.md) (現在)
- **中文** - [README.zh.md](README.zh.md)
- **Español** - [README.es.md](README.es.md)
- **हिंदी** - [README.hi.md](README.hi.md)
- **Português** - [README.pt.md](README.pt.md)

> **注意**: 英語以外のすべてのバージョンは機械翻訳されています。最も正確な情報については、英語版を参照してください。

## クイックスタート

### インストール

[リリースページ](https://github.com/superdaigo/gsecutil/releases)から、お使いのプラットフォーム用の最新バイナリをダウンロードしてください：

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

または Go でインストール：
```bash
go install github.com/superdaigo/gsecutil@latest
```

### 前提条件

- [Google Cloud SDK (gcloud)](https://cloud.google.com/sdk/docs/install) がインストールされ、認証されていること
- Secret Manager API が有効化された Google Cloud プロジェクト

### 認証

```bash
# gcloud で認証
gcloud auth login

# デフォルトプロジェクトを設定
gcloud config set project YOUR_PROJECT_ID

# または環境変数を設定
export GSECUTIL_PROJECT=YOUR_PROJECT_ID
```

## 基本的な使い方

### シークレットの作成
```bash
# 対話的な入力
gsecutil create database-password

# コマンドラインから
gsecutil create api-key -d "sk-1234567890"

# ファイルから
gsecutil create config --data-file ./config.json
```

### シークレットの取得
```bash
# 最新バージョンを取得
gsecutil get database-password

# クリップボードにコピー
gsecutil get api-key --clipboard

# 特定のバージョンを取得
gsecutil get api-key --version 2
```

### シークレットの一覧表示
```bash
# すべてのシークレットを表示
gsecutil list

# ラベルでフィルター
gsecutil list --filter "labels.env=prod"
```

### シークレットの更新
```bash
# 対話的な入力
gsecutil update database-password

# コマンドラインから
gsecutil update api-key -d "new-secret-value"
```

### シークレットの削除
```bash
gsecutil delete old-secret
```

## 設定

gsecutil はプロジェクト固有の設定のための設定ファイルをサポートしています。設定ファイルは次の順序で検索されます：

1. `--config` フラグ（指定されている場合）
2. カレントディレクトリ: `gsecutil.conf` または `.gsecutil.conf`
3. ホームディレクトリ: `~/.config/gsecutil/gsecutil.conf`

### 設定例

```yaml
# プロジェクト ID（環境変数または gcloud で設定されている場合はオプション）
project: "my-project-id"

# チーム組織のためのシークレット名プレフィックス
prefix: "team-shared-"

# list コマンドで表示するデフォルトの属性
list:
  attributes:
    - title
    - owner
    - environment

# 認証情報のメタデータ
credentials:
  - name: "database-password"
    title: "Production Database Password"
    environment: "production"
    owner: "backend-team"
```

### クイックスタート

```bash
# 対話的に設定を生成
gsecutil config init

# またはプロジェクト固有の設定を作成
echo 'project: "my-project-123"' > gsecutil.conf
```

詳細な設定オプションについては、[docs/configuration.md](docs/configuration.md) を参照してください。

## 主な機能

- ✅ **シンプルな CRUD 操作** - シークレットを管理するための直感的なコマンド
- ✅ **クリップボード統合** - シークレットをクリップボードに直接コピー
- ✅ **バージョン管理** - 特定のバージョンへのアクセスとバージョンライフサイクルの管理
- ✅ **設定ファイルサポート** - チームフレンドリーなメタデータと組織化
- ✅ **アクセス管理** - 基本的な IAM ポリシー管理
- ✅ **監査ログ** - 誰がいつシークレットにアクセスしたかを表示
- ✅ **複数の入力方法** - 対話的、インライン、またはファイルベース
- ✅ **クロスプラットフォーム** - Linux、macOS、Windows（amd64/arm64）

## ドキュメント

- **[設定ガイド](docs/configuration.md)** - 詳細な設定オプションと例
- **[コマンドリファレンス](docs/commands.md)** - 完全なコマンドドキュメント
- **[監査ログの設定](docs/audit-logging.md)** - 監査ログの有効化と使用
- **[トラブルシューティングガイド](docs/troubleshooting.md)** - 一般的な問題と解決策
- **[ビルド手順](BUILD.md)** - ソースからビルド
- **[開発ガイド](WARP.md)** - WARP AI での開発

## よく使うコマンド

```bash
# シークレットの詳細を表示
gsecutil describe my-secret

# バージョン履歴を表示
gsecutil describe my-secret --show-versions

# 監査ログを表示
gsecutil auditlog my-secret

# アクセスを管理
gsecutil access list my-secret
gsecutil access grant my-secret --principal user:alice@example.com

# 設定を検証
gsecutil config validate

# 設定を表示
gsecutil config show
```

## ライセンス

このプロジェクトは MIT ライセンスの下でライセンスされています - 詳細は LICENSE ファイルを参照してください。

## 関連

- [Google Cloud SDK](https://cloud.google.com/sdk)
- [Secret Manager ドキュメント](https://cloud.google.com/secret-manager/docs)
