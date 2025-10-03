# gsecutil - Google Secret Manager ユーティリティ

> **翻訳について**: このREADMEファイルは機械翻訳されています。最新かつ正確な情報については、英語版の[README.md](README.md)をご参照ください。

🚀 **v1.1.0** - 設定ファイルサポート付きGoogle Secret Managerのシンプルなコマンドラインラッパー。`gsecutil`は一般的な秘密操作のための便利なコマンドを提供し、小規模なチームが専用のパスワード管理ツールを必要とせずにGoogle CloudのSecret Managerを使用してパスワードや認証情報を管理しやすくします。

**v1.1.0の新機能**: YAML設定ファイルサポート、プレフィックス機能、チームカスタムメタデータ付きの拡張リストおよび説明コマンド。

## ✨ 機能

### 🔐 **完全な秘密管理**
- **CRUD操作**: 簡略化されたコマンドで秘密の作成、読み取り、更新、削除
- **バージョン管理**: 任意のバージョンへのアクセス、バージョン履歴とメタデータの表示
- **クロスプラットフォーム**サポート（Linux、macOS、Windows、ARM64サポート）
- **クリップボード統合** - 秘密の値を直接クリップボードにコピー
- **インタラクティブ＆ファイル入力** - セキュアプロンプトまたはファイルベースの秘密読み込み

### 🛡️ **高度なアクセス管理**
*(v1.0.0で導入)*
- **完全なIAMポリシー分析** - 任意のレベルで秘密にアクセスできるユーザーを表示
- **マルチレベル権限チェック** - 秘密レベルとプロジェクトレベルのアクセス分析
- **IAM条件認識** - CEL式を使用した条件付きアクセスポリシーの完全サポート
- **プリンシパル管理** - ユーザー、グループ、サービスアカウントのアクセス権限の付与/取り消し
- **プロジェクト全体分析** - Secret Managerアクセスを提供するEditor/Ownerロールの特定

### 📊 **監査とコンプライアンス**
- **包括的な監査ログ** - 誰が秘密にアクセスし、いつ、どの操作を実行したかを追跡
- **プリンシパルベースのフィルタリング** - 特定のユーザー/グループがアクセス可能な全ての秘密を表示
- **柔軟なフィルタリング** - 秘密、プリンシパル、操作タイプ、時間範囲による絞り込み
- **条件評価** - 条件付きアクセスがいつ適用されるかを理解

### 🎯 **本番環境対応**
- **一貫したAPI** - 全てのコマンドで統一されたパラメータ命名
- **エンタープライズ機能** - IAM条件、プロジェクトレベル分析、コンプライアンス監査
- **堅牢なエラーハンドリング** - 権限不足やネットワーク問題の優雅な処理
- **柔軟な出力** - JSON、YAML、テーブル形式とリッチフォーマット

## 前提条件

- [Google Cloud SDK (gcloud)](https://cloud.google.com/sdk/docs/install) がインストールされ、認証されている
- Secret Manager APIが有効化されたGoogle Cloudプロジェクト
- Secret Manager操作のための適切なIAM権限

## インストール

### プリビルトバイナリ

[リリースページ](https://github.com/superdaigo/gsecutil/releases)から、お使いのプラットフォーム向けの最新リリースをダウンロードしてください：

| プラットフォーム | アーキテクチャ | ダウンロード |
|----------|--------------|----------|
| Linux | x64 | `gsecutil-linux-amd64-v{version}` |
| Linux | ARM64 | `gsecutil-linux-arm64-v{version}` |
| macOS | Intel | `gsecutil-darwin-amd64-v{version}` |
| macOS | Apple Silicon | `gsecutil-darwin-arm64-v{version}` |
| Windows | x64 | `gsecutil-windows-amd64-v{version}.exe` |

**ダウンロード後:** 一貫した使用のためにバイナリをリネーム：

```bash
# Linux/macOS 例:
mv gsecutil-linux-amd64-v1.1.0 gsecutil
chmod +x gsecutil

# Windows 例 (PowerShell/Command Prompt):
ren gsecutil-windows-amd64-v1.1.0.exe gsecutil.exe
```

これにより、バージョンに関係なく`gsecutil`を一貫して使用できます。

### Goでインストール

```bash
go install github.com/superdaigo/gsecutil@latest
```

### ソースからビルド

包括的なビルド手順については、[BUILD.md](BUILD.md)をご覧ください。

**クイックビルド:**
```bash
git clone https://github.com/superdaigo/gsecutil.git
cd gsecutil

# 現在のプラットフォーム用にビルド
make build
# または
./build.sh          # Linux/macOS
.\\build.ps1         # Windows

# 全プラットフォーム用にビルド
make build-all
# または
./build.sh all      # Linux/macOS
.\\build.ps1 all     # Windows
```

## 使用方法

### グローバルオプション

- `-p, --project`: Google Cloudプロジェクト ID（`GOOGLE_CLOUD_PROJECT`環境変数でも設定可能）

### コマンド

#### Get Secret（秘密の取得）

Google Secret Managerから秘密の値を取得します。デフォルトでは最新バージョンを取得しますが、任意のバージョンを指定できます：

```bash
# 秘密の最新バージョンを取得
gsecutil get my-secret

# 特定のバージョンを取得（ロールバック、デバッグ、履歴値へのアクセスに便利）
gsecutil get my-secret --version 1
gsecutil get my-secret -v 3

# 秘密を取得してクリップボードにコピー
gsecutil get my-secret --clipboard

# 特定のバージョンをクリップボードで取得
gsecutil get my-secret --version 2 --clipboard

# バージョンメタデータ付きで秘密を取得（バージョン、作成時刻、状態）
gsecutil get my-secret --show-metadata
gsecutil get my-secret -v 1 --show-metadata    # 古いバージョンをメタデータ付きで

# 特定のプロジェクトから秘密を取得
gsecutil get my-secret --project my-gcp-project
```

**バージョンサポート:**
- 🔄 **最新バージョン**: `--version`が指定されていない場合のデフォルト動作
- 📅 **履歴バージョン**: 番号による任意の以前のバージョンへのアクセス（例：`--version 1`、`--version 2`）
- 🔍 **バージョンメタデータ**: `--show-metadata`を使用してバージョンの詳細を表示（作成時刻、状態、ETag）
- 📋 **クリップボードサポート**: `--clipboard`を使用して任意のバージョンで動作

## 設定

### 環境変数

- `GOOGLE_CLOUD_PROJECT`: デフォルトプロジェクトID（`--project`フラグで上書き）

### 認証

`gsecutil`は`gcloud`と同じ認証を使用します。認証されていることを確認してください：

```bash
# gcloudで認証
gcloud auth login

# デフォルトプロジェクトを設定
gcloud config set project YOUR_PROJECT_ID

# サービスアカウント用（CI/CDで）
gcloud auth activate-service-account --key-file=service-account.json
```

### シェル補完

`gsecutil`はbash、zsh、fish、PowerShellのシェル自動補完をサポートしています。これにより、コマンド、フラグ、オプションのタブ補完が可能になり、CLIがより使いやすくなります。

#### セットアップ手順

**Bash:**
```bash
# 一時的（現在のセッションのみ）
source <(gsecutil completion bash)

# 永続インストール（bash-completionパッケージが必要）
# システム全体（sudoが必要）
sudo gsecutil completion bash > /etc/bash_completion.d/gsecutil

# ユーザーローカルインストール
gsecutil completion bash > ~/.local/share/bash-completion/completions/gsecutil

# または~/.bashrcに追加して自動読み込み
echo 'source <(gsecutil completion bash)' >> ~/.bashrc
```

**Zsh:**
```bash
# 一時的（現在のセッションのみ）
source <(gsecutil completion zsh)

# 永続インストール
gsecutil completion zsh > "${fpath[1]}/_gsecutil"

# または~/.zshrcに追加して自動読み込み
echo 'source <(gsecutil completion zsh)' >> ~/.zshrc
```

**Fish:**
```bash
# 一時的（現在のセッションのみ）
gsecutil completion fish | source

# 永続インストール
gsecutil completion fish > ~/.config/fish/completions/gsecutil.fish
```

**PowerShell:**
```powershell
# PowerShellプロファイルに追加
gsecutil completion powershell | Out-String | Invoke-Expression

# またはプロファイルに保存して自動読み込み
gsecutil completion powershell >> $PROFILE
```

#### 機能

インストール後、シェル補完は以下を提供します：
- **コマンド補完**: `gsecutil`サブコマンドのタブ補完（`get`、`create`、`list`など）
- **フラグ補完**: `--project`、`--version`、`--clipboard`などのフラグのタブ補完
- **スマート提案**: 現在のコマンドに基づいたコンテキストを考慮した補完
- **ヘルプテキスト**: コマンドやフラグの簡潔な説明（サポートされている場合）

#### 使用例

```bash
# タイプしてTabを押し、使用可能なコマンドを表示
gsecutil <Tab>
# 表示: access, auditlog, completion, create, delete, describe, get, help, list, update

# 部分コマンドをタイプしてTabで補完
gsecutil des<Tab>
# 補完結果: gsecutil describe

# フラグのタブ補完も動作
gsecutil get my-secret --<Tab>
# 表示: --clipboard, --project, --show-metadata, --version
```

**注意**: 補完が有効になるために、シェルを再起動するか、シェル設定ファイルをソースする必要がある場合があります。

## セキュリティとベストプラクティス

### セキュリティ機能

- **永続ストレージなし**: 秘密の値は`gsecutil`によってログや永続的な保存は行われません
- **セキュアな入力**: インタラクティブプロンプトは非表示パスワード入力を使用
- **OS ネイティブクリップボード**: クリップボード操作はセキュアなOSネイティブAPIを使用
- **gcloud委譲**: すべての操作は認証済みの`gcloud` CLIに委譲

### ベストプラクティス

- **`--force`を慎重に使用**: 自動化環境で`--force`を使用する前は常にレビュー
- **環境変数**: 繰り返しの`--project`フラグを避けるために`GOOGLE_CLOUD_PROJECT`を設定
- **バージョン制御**: 本番環境では特定の秘密バージョンを使用（`--version N`）
- **定期的な監査**: `gsecutil auditlog secret-name`で秘密のアクセスを監視
- **秘密のローテーション**: `gsecutil update`を使用した定期的な秘密のローテーション

## トラブルシューティング

### よくある問題

1. **"gcloud command not found"**
   - Google Cloud SDKがインストールされ、`gcloud`がPATHに含まれていることを確認

2. **認証エラー**
   - `gcloud auth login`を実行して認証
   - プロジェクトアクセスを検証: `gcloud config get-value project`

3. **権限拒否エラー**
   - アカウントが必要なIAMロールを持っていることを確認:
     - `roles/secretmanager.admin`（すべての操作用）
     - `roles/secretmanager.secretAccessor`（読み取り操作用）
     - `roles/secretmanager.secretVersionManager`（作成/更新操作用）

4. **クリップボードが動作しない**
   - グラフィカル環境があることを確認（Linux用）
   - ヘッドレスサーバーでは、クリップボード操作は優雅に失敗する場合があります

### デバッグモード

以下を設定してgcloudコマンドに詳細出力を追加:

```bash
export CLOUDSDK_CORE_VERBOSITY=debug
```

## ドキュメント

- **[BUILD.md](BUILD.md)** - すべてのプラットフォームの包括的なビルド手順
- **[CONTRIBUTING.md](CONTRIBUTING.md)** - 貢献ガイドラインと開発ワークフロー
- **[WARP.md](WARP.md)** - WARP AIターミナル統合のための開発ガイダンス
- **README.md** - このファイル、使用方法と概要

## 貢献

貢献を歓迎します！このプロジェクトへの貢献方法、開発環境とpre-commitフックのセットアップ手順を含む詳細なガイドラインについては、[CONTRIBUTING.md](CONTRIBUTING.md)をご覧ください。

## ライセンス

このプロジェクトはMITライセンスの下でライセンスされています - 詳細についてはLICENSEファイルをご覧ください。

## 関連プロジェクト

- [Google Cloud SDK](https://cloud.google.com/sdk)
- [Secret Manager Documentation](https://cloud.google.com/secret-manager/docs)
- [Cobra CLI Framework](https://github.com/spf13/cobra)
