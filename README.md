<img src="./assets/logo_title.png" width="248px" height="248px"></image>

mosscat は社内向けの OSS 利用状況を一元管理するためのバックエンドサービスです。Go 言語と [Echo](https://echo.labstack.com/) を用いた REST API として実装されており、PostgreSQL もしくは SQLite をデータストアとして利用できます。

## 機能概要

- OSS コンポーネントおよびバージョンの CRUD
- プロジェクトと OSS 利用状況 (Usage) の管理
- タグ付け、スコープポリシー判定、監査ログ取得
- OpenAPI (\`internal/api/openapi.yaml\`) に基づくサーバ実装

## 実行方法

1. Go 1.24 以降がインストールされた環境でリポジトリを取得します。
2. 必要に応じて `config.yaml` の `db.dsn` にデータベース DSN を指定してください (未指定の場合は SQLite のメモリ DB を利用します)。
3. 以下のコマンドでサーバを起動します。

```bash
$ go run .
```

`config.yaml` を用意することで待ち受けホスト・ポートを変更できます。
デフォルトでは `0.0.0.0:8080` で起動します。生成された API ハンドラは Echo のミドルウェアによりリクエスト検証が行われます。
`server.allowed_origins` を設定することで CORS 許可オリジンを指定できます (省略時は `*`)。

## Windows サービスとしての登録と実行

Windows 環境ではビルドしたバイナリをサービスとして登録できます。以下は 64bit Windows 用バイナリを例とした手順です。

```bash
# バイナリをビルド
$ GOOS=windows GOARCH=amd64 go build -o mosscat.exe

# サービス登録
$ mosscat.exe -service install

# サービス開始 (PowerShell または管理ツールから実行)
> Start-Service mosscat

# サービス削除
$ mosscat.exe -service uninstall
```

サービスとして実行された場合も `config.yaml` の内容が利用されます。未設定時は `0.0.0.0:8080` で待ち受けます。
