# gsecutil - Google Secret Manager ユーティリティ

Google Secret Manager の簡略化されたコマンドラインラッパーです。プロジェクト単位のパスワードマネージャーのように動作し、直感的なコマンド、クリップボード統合、バージョン管理、チームフレンドリーな設定ファイル、監査ログでシークレットを保存、取得、管理できます。

## 🌍 言語バージョン

- **English** - [README.md](README.md)
- **日本語** - [README.ja.md](README.ja.md)（現在）
- **中文** - [README.zh.md](README.zh.md)
- **Español** - [README.es.md](README.es.md)
- **हिंदी** - [README.hi.md](README.hi.md)
- **Português** - [README.pt.md](README.pt.md)

> **注意**: 英語以外のすべてのバージョンは機械翻訳されています。最も正確な情報については、英語版を参照してください。

## クイックスタート

### インストール

お使いのプラットフォーム用の最新バイナリを[リリースページ](https://github.com/superdaigo/gsecutil/releases)からダウンロードするか、Go でインストールしてください：

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

各プロジェクトには通常、プロジェクト ID、シークレット命名規則、メタデータ属性を保存する独自の設定ファイルがあります。

### 1. 設定ファイルを作成する

対話型セットアップを実行して設定ファイルを生成します。Google Cloud プロジェクト ID、シークレット名のプレフィックス、デフォルトの一覧表示属性、オプションのサンプル認証情報を入力するよう求められます。生成されたファイルはデフォルトでカレントディレクトリに `gsecutil.conf` として保存されます（`--home` を使用すると `~/.config/gsecutil/gsecutil.conf` に保存されます）。

```bash
gsecutil config init
```

設定ファイルは以下の順序で検索されます：
1. `--config` フラグ（指定されている場合）
2. カレントディレクトリ：`gsecutil.conf`
3. ホームディレクトリ：`~/.config/gsecutil/gsecutil.conf`

### 2. シークレットを管理する

```bash
# シークレットを作成
gsecutil create database-password

# 最新バージョンを取得
gsecutil get database-password

# クリップボードにコピー
gsecutil get database-password --clipboard

# すべてのシークレットを一覧表示
gsecutil list

# シークレットを更新
gsecutil update database-password

# シークレットを削除
gsecutil delete database-password
```

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

# 認証情報メタデータ（名前はベア名 — プレフィックスは自動的に付加）
credentials:
  - name: "database-password"    # "team-shared-database-password" にアクセス
    title: "Production Database Password"
    environment: "production"
    owner: "backend-team"
```

> **プレフィックスは透過的:** プレフィックスが設定されている場合、コマンド、設定、CSV ファイルでは常にベア名を使用します。プレフィックスは自動的に付加・削除されます。

詳細な設定オプションについては、[docs/configuration.md](docs/configuration.md) を参照してください。

## ドキュメント

- **[設定ガイド](docs/configuration.md)** - 詳細な設定オプションと例
- **[コマンドリファレンス](docs/commands.md)** - 完全なコマンドドキュメント
- **[監査ログの設定](docs/audit-logging.md)** - 監査ログの有効化と使用
- **[トラブルシューティングガイド](docs/troubleshooting.md)** - 一般的な問題と解決策
- **[ビルド手順](BUILD.md)** - ソースからビルド
- **[開発ガイド](WARP.md)** - WARP AI での開発

## ライセンス

このプロジェクトは MIT ライセンスの下でライセンスされています - 詳細は LICENSE ファイルを参照してください。

## 関連

- [Google Cloud SDK](https://cloud.google.com/sdk)
- [Secret Manager ドキュメント](https://cloud.google.com/secret-manager/docs)
