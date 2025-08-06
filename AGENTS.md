# AGENTS.md (AIエージェント運用メタ情報)

**目的:** このファイルはリポジトリ内で *AI コーディングエージェント / LLM* が安全かつ一貫性を保って変更・生成を行うための最小限かつ機械が読みやすいガイドラインを提供する。詳細背景は `go-ai-agent-design.md` を参照し、本ファイルは “行動制約” と “主要コマンド” を列挙する。

---

## 1. プロジェクトメタデータ

| Key             | Value                     |
| --------------- | ------------------------- |
| project         | mosscat-backend       |
| language        | go                        |
| go.min\_version | 1.24.x                    |
| framework.http  | echo v4                   |
| openapi.spec    | internal/api/openapi.yaml |
| codegen.tool    | oapi-codegen v2.x         |
| db.primary      | postgres (16)             |
| db.optional     | sqlite (ローカル/軽量用途)        |
| auth.method     | basic (将来: ldap/jwt)      |
| logging         | slog (構造化)                |
| migration.tool  | golang-migrate            |

---

## 2. ディレクトリ役割

| Path                       | Role               | AI編集可否         |
| -------------------------- | ------------------ | -------------- |
| cmd/api                    | 起動エントリ             | △ (最小変更)       |
| internal/api/openapi.yaml  | OpenAPI仕様          | × (人間レビュー後に更新) |
| internal/api/gen           | oapi-codegen 生成コード | × (自動生成のみ)     |
| internal/api/handler       | Echoハンドラ実装層        | 〇              |
| internal/domain/model      | ドメインエンティティ         | 〇 (互換性注意)      |
| internal/domain/service    | ビジネスロジック           | 〇              |
| internal/domain/repository | リポジトリIF            | 〇 (破壊的変更要注意)   |
| internal/policy            | スコープ判定             | 〇              |
| internal/infra/db          | 接続・Tx管理            | 〇              |
| internal/infra/repository  | SQL実装              | 〇              |
| internal/pkg/errors        | エラー型               | 〇              |
| internal/pkg/response      | Problem変換          | 〇              |
| internal/logger            | ロガー初期化             | △              |
| internal/config            | 設定ロード              | △              |
| scripts                    | 補助スクリプト            | △              |
| migrations                 | SQLマイグレーション        | 〇 (新規追加のみ)     |
| AGENTS.md                  | 本ファイル              | △ (構造維持)       |
| api-design.md              | 詳細設計               | × (人間主導)       |

凡例: 〇=自由, △=要慎重, ×=AIエージェント直接編集禁止。

---

## 3. 生成 (Codegen) ポリシー

**禁止:** `internal/api/gen` 以下を手動編集。

| タスク            | コマンド                                                     | 目的            |
| -------------- | -------------------------------------------------------- | ------------- |
| types+server生成 | `go generate`                                          | OpenAPI仕様変更反映 |
| 差分検査           | `go generate && git diff --exit-code internal/api/gen` | 未生成検出         |

> OpenAPI スキーマの *modelsのみ変更* の場合でも `go generate` を実行し、`api.gen.go` の差分を確認。

### SQLite 併用指針

* DSN が `postgres://` で始まらない場合 SQLite と解釈。DDL 方言差を避けるため ENUM や複雑機能未使用。
* 配列系 (layers, cpe\_list) は PostgreSQL: TEXT\[] / SQLite: JSON 文字列。

---

## 4. 命名 / コーディング規約（抜粋）

| 対象        | 規約                                        |
| --------- | ----------------------------------------- |
| Struct名   | PascalCase (`OssComponent`, `OssVersion`) |
| Interface | `<Name>Repository`, `<Name>Service`       |
| Receiver  | 1文字 (`r`, `s`, `h`)                       |
| エラー       | `errors.Wrap` 等で文脈保持。最終HTTP化は handler 内   |
| Context   | 第一引数に `ctx context.Context`               |
| ロギングキー    | `req_id`, `path`, `status`, `latency_ms`  |

---

## 5. 変更禁止 / 注意事項

| 項目         | ルール                      |
| ---------- | ------------------------ |
| 生成コード      | 手動編集禁止。差分あれば仕様同期をまず実施    |
| OpenAPI仕様  | エンドポイント追加前に人間レビュー必須      |
| ハードコード秘密   | 禁止 (パスワード, APIキー)        |
| `panic`    | ライブラリエラー以外禁止。AppErrorで処理 |
| SQL        | プレースホルダ必須。文字列連結で条件付与禁止   |
| Passwordログ | 出力禁止                     |

---

## 6. サービス層基本パターン

1. 入力DTO→ドメイン変換
2. バリデーション / 正規化 (名称小文字化など)
3. Tx開始 (`TxManager.WithinTx`)
4. Repository 呼び出し
5. ポリシー/判定適用 (スコープ)
6. AuditLog 記録（失敗しても主要処理は成功させる）
7. ドメイン→DTO 変換返却

---

## 7. スコープ初期判定 (簡易)

```
role in [BUILD_ONLY,DEV_ONLY,TEST_ONLY] => OUT_SCOPE
role == SERVER_ENV && !policy.serverEnvIncluded => OUT_SCOPE
role == RUNTIME_REQUIRED && !policy.runtimeRequiredDefaultInScope => REVIEW_NEEDED
その他 => IN_SCOPE
```

---

## 8. 共通エラーコード

| Code       | HTTP | 意味       |
| ---------- | ---- | -------- |
| NOT\_FOUND | 404  | 対象リソース無し |
| CONFLICT   | 409  | 一意制約衝突   |
| VALIDATION | 400  | 入力検証失敗   |
| INTERNAL   | 500  | 予期せぬ内部障害 |

---

## 9. ページング標準

| Query | 説明     |
| ----- | ------ |
| page  | 1 起点   |
| size  | 1..200 |

レスポンス: `{items, page, size, total}`

---

## 11. テスト方針（最小）

| 種別                    | 必須カバレッジ目安      | 主対象            |
| --------------------- | -------------- | -------------- |
| Unit (service/policy) | >=60%          | ビジネス分岐         |
| Repository            | 主要CRUD         | SQL制約 / エラー変換  |
| Handler               | 主要成功 + 404/400 | DTO/Problem 変換 |

失敗時: `t.Helper()` 使用、`require` 系で即時失敗。

和田卓人（twada）さんが行っているテスト駆動開発を基本としてください。
---

## 12. ワークフロー (AI 推奨手順)

1. 仕様差分確認 (openapi.yaml 変更?)
2. `go generate` 実行し生成差分確認
3. 新/変更ハンドラ stub を `handler` に追加
4. Service / Repository 実装・テスト更新
5. `go vet` / `go test` / generate 差分ゼロを確認
6. 変更概要と影響範囲を出力 (人間レビュー用)

---

## 13. 提示すべきレビュー出力フォーマット例 (AI→人間)

```
# Change Summary
- Added: handler/ListOssComponents pagination logic
- Modified: repository/oss_component_repository.go (Search: tag filter)
- DB: No schema change
- OpenAPI: No change
- Tests: Added TestSearchWithTag
```

---

## 14. AI ハンドラ生成テンプレ (短縮)

```
目的: ListOssComponents 実装
前提: services.OssComponent.Search(ctx, filter) -> ([]model.OssComponent,total,error)
入力: name, layers, tag, page, size
出力: 200 JSON (Paged)
失敗: domain error -> Problem
制約: SQL組立は repository 層で行うこと (handler でやらない)
```

---

## 15. DoD (Definition of Done) チェックリスト

* [ ] `go generate` 後に生成差分なし
* [ ] 主要ビジネス分岐のテスト追加 / 既存テスト緑
* [ ] 新規公開エンドポイントはエラーパス(404/400)テスト完備
* [ ] ログに機微情報なし
* [ ] Lint (`go vet`) 問題なし
* [ ] README / AGENTS.md 必要箇所更新済み (必要なら)

---

## 16. NG パターン対策

| NG              | 対応                          |
| --------------- | --------------------------- |
| 生SQL文字列連結       | placeholder 使用 / builder 使用 |
| ハンドラでビジネスロジック   | service 層へ移動                |
| panic 乱用        | error 戻り値化                  |
| 生成コード直接編集       | OpenAPI 更新→再生成              |
| テストで sleep 固定待機 | コンテキスト or 即時検証              |

---

## 17. 変更提案の出し方 (AI)

AI は修正 PR 相当を提案する際、以下を JSON で併記可能:

```json
{
  "summary": "Add tag filter to list components",
  "files": [
    {"path": "internal/domain/repository/oss_component_repository.go", "action": "modify"},
    {"path": "internal/api/handler/oss_component_handler.go", "action": "modify"},
    {"path": "internal/domain/service/oss_component_service_test.go", "action": "add"}
  ]
}
```

---

## 18. 最小メタ更新ルール

| 変更種別     | AGENTS.md 更新要否            |
| -------- | ------------------------- |
| 新エンドポイント | 〇 (ワークフロー/コマンド不変なら概要追記のみ) |
| Enum 値追加 | △ (必要なら簡潔追記)              |
| 内部実装細部   | ×                         |

---

## 19. 安全ガード

* AGENTS.md 自体を大幅書換する提案は要人間承認コメントを付与すること。
* OpenAPI スキーマ未変更で生成差分が出た場合、まずスキーマ/生成設定の不整合を報告する。

---

## 20. 参照

| 種類        | Path                           |
| --------- | ------------------------------ |
| 詳細設計      | go-ai-agent-design.md          |
| OpenAPI   | internal/api/openapi.yaml      |
| マイグレーション  | migrations/                    |
| 主要コード生成設定 | internal/api/oapi-codegen.yaml |

---

**EOF**
