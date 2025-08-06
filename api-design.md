## 0. サービス概要（目的サマリ）

**目的:** 社内プロジェクトが *納品対象となる OSS コンポーネント* を統一管理し、バージョン・ライセンス・改変有無・利用形態 (UsageRole)・スコープ判定 (IN/OUT/REVIEW) を登録・検索・エクスポートできる基盤を提供する。

| 観点      | 要約                                                                      |
| ------- | ----------------------------------------------------------------------- |
| 主ユースケース | プロジェクトで利用する OSS を登録し、納品一覧（IN\_SCOPE）だけを抽出 (将来: SPDX / NOTICE)           |
| ドメイン中核  | **OssComponent (親)** と **OssVersion (子)**, **ProjectUsage (結合 + スコープ)** |
| スコープ分類  | UsageRole + ポリシー (ScopePolicy) により自動初期判定 / 人手で上書き                       |
| 除外対象    | 開発専用 (DEV/BUILD/TEST), サーバ環境のみ (SERVER\_ENV), ポリシー設定で制御                 |
| 非機能     | 単一ノード / オフライン / PostgreSQL 本番 + SQLite 簡易運用                             |
| 将来拡張    | 脆弱性(NVD)・SBOM Import・NOTICE 生成・LDAP/JWT 認証                              |

---

## 1. アーキテクチャ要点

| 層                             | 役割                                   | 実装                               |
| ----------------------------- | ------------------------------------ | -------------------------------- |
| API (Echo)                    | OpenAPI に沿った Handler                 | oapi-codegen server stub + 手書き実装 |
| DTO/Gen                       | 自動生成型                                | `internal/api/gen` (編集禁止)        |
| Domain (model/service/policy) | ビジネスロジック / 正規化 / スコープ判定              | 手書き                              |
| Repository                    | DB アクセス (SQL)                        | `infra/repository`               |
| Infra                         | DB接続 / Tx / Migration                | `infra/db`, `migration`          |
| Shared                        | errors / auth / pagination / logging | `pkg/*`                          |

**データフロー:** HTTP → Handler → Service (Tx) → Repository → DB → 変換 → Handler Response

---

## 2. ディレクトリ (抜粋)

```
cmd/api/main.go
internal/api/openapi.yaml
internal/api/gen/        # 生成 (types, server)
internal/api/handler/    # ServerInterface 実装
internal/domain/model/
internal/domain/service/
internal/domain/policy/
internal/domain/repository/ (interface)
internal/infra/repository/ (sql impl)
internal/infra/db/
internal/infra/migration/
pkg/errors/ pkg/auth/ pkg/response/
AGENTS.md
```

---

## 3. コード生成方針 (types/server 分離)

| 生成物    | 目的       | 再生成条件                  | コマンド                   |
| ------ | -------- | ---------------------- | ---------------------- |
| types  | Schema 型 | components/schemas 変更  | `make generate-types`  |
| server | ルート & IF | paths / operationId 変更 | `make generate-server` |

`AGENTS.md` に記載されたコマンドを唯一の真実とし、生成差分はコミット前にゼロ確認。生成ファイル直接編集禁止（CI で検査）。

---

## 4. 主要エンティティ（最小）

| Entity       | 主キー       | 主フィールド                                                                                                                          | 備考                       |
| ------------ | --------- | ------------------------------------------------------------------------------------------------------------------------------- | ------------------------ |
| OssComponent | id (UUID) | name, normalized\_name, layers\[], deprecated                                                                                   | タグは多対多                   |
| OssVersion   | id        | oss\_id(FK), version, license\_expression\_raw, purl, cpe\_list\[], hash\_sha256, review\_status, scope\_status, supplier\_type | (scope\_status は推奨レベル)   |
| Project      | id        | project\_code, name, delivery\_date                                                                                             |                          |
| ProjectUsage | id        | project\_id, oss\_id, oss\_version\_id, usage\_role, scope\_status, inclusion\_note                                             | 最終スコープ判定元                |
| ScopePolicy  | id        | runtime\_required\_default\_in\_scope, server\_env\_included, auto\_mark\_forks\_in\_scope                                      | 1レコード想定                  |
| Tag          | id        | name                                                                                                                            | 中間: oss\_component\_tags |
| AuditLog     | id        | entity\_type, entity\_id, action, user\_name, summary                                                                           | 監査最小                     |

---

## 5. スコープ初期判定ロジック (擬似)

```go
switch usageRole {
case "BUILD_ONLY", "DEV_ONLY", "TEST_ONLY": return OUT_SCOPE
case "SERVER_ENV":  if policy.ServerEnvIncluded { return IN_SCOPE } else { return OUT_SCOPE }
case "RUNTIME_REQUIRED": if policy.RuntimeRequiredDefaultInScope { return IN_SCOPE } else { return REVIEW_NEEDED }
default: return IN_SCOPE // Bundled / Linked 系
}
```

最終的な納品抽出は `ProjectUsage.scope_status == IN_SCOPE` を条件とする。

---

## 6. PostgreSQL / SQLite 併用指針（簡略）

| 項目             | PostgreSQL        | SQLite    | 実装メモ              |
| -------------- | ----------------- | --------- | ----------------- |
| UUID           | アプリ生成 or pgcrypto | アプリ生成     | `uuid.New()` 統一   |
| 配列(layers,cpe) | TEXT\[]           | JSON テキスト | Repository で変換    |
| トランザクション       | 標準                | 標準        | TxManager 共通      |
| Migration      | migrate           | migrate   | 同一 SQL (方言回避)     |
| 限界             | 高並行               | 単ユーザ向け    | SQLite は開発/軽量用途のみ |

---

## 7. Handler 実装基本形

```go
func (h *Handler) ListOssComponents(c echo.Context, p gen.ListOssComponentsParams) error {
 ctx := c.Request().Context()
 res, total, err := h.services.OssComponent.Search(ctx, toFilter(p))
 if err != nil { return h.problem(c, err) }
 return c.JSON(http.StatusOK, toPaged(res, p.Page, p.Size, total))
}
```

---

## 8. エラー標準化

| Code       | 意味      |
| ---------- | ------- |
| NOT\_FOUND | 資源なし    |
| CONFLICT   | 一意制約衝突  |
| VALIDATION | 入力不正    |
| INTERNAL   | 予期しない失敗 |

変換: `errors.AppError` → OpenAPI `Problem`.

---

## 9. テスト最小セット

| 種類        | 内容                            |
| --------- | ----------------------------- |
| Unit      | ScopeDecider, Service (成功/異常) |
| Repo      | CRUD + 一意制約 + ページング           |
| API       | 代表エンドポイント 200/400/404         |
| Migration | Up/Down 成功確認                  |

---

## 10. 開発フロー (要約)

1. OpenAPI 変更 → `make generate-server` (+ schemas 変更なら `generate-types`)
2. ハンドラ stub 差分確認 → 実装/更新
3. サービス & リポジトリ改修 / テスト追加
4. `go vet`, `go test`, `generate-check` パス
5. PR: 生成ファイル差分最小 / 変更理由明記

---

## 11. AI エージェントへの指示テンプレート（短縮）

**Handler 生成:**

```
OpenAPI の ListOssComponents に対する Echo ハンドラを実装。入力 params -> filter 変換し Service 呼出。エラーは h.problem。
```

**Repository 検索:**

```
OssComponent.Search を実装。name 部分一致(normalized_name ILIKE)、layer 配列 OR 条件、tag JOIN は LEFT (タグ未指定時 JOIN 省略)。
```

**Scope 判定テスト:**

```
InitialScope の各 usageRole ケース網羅テスト（テーブルドリブン）。
```

---

## 12. Definition of Done (簡略)

| 項目     | 条件                             |
| ------ | ------------------------------ |
| 生成差分   | `make generate` 後に未コミット変更なし    |
| テスト    | 主要ユニット & API パス (最低 60% カバレッジ) |
| Lint   | `go vet` 問題なし                  |
| エラー整合  | 全エンドポイントで Problem 形式           |
| スコープ判定 | 新 UsageRole 追加時テスト更新           |

---

## 13. 将来拡張（占位のみ）

| 項目          | 予定                              |
| ----------- | ------------------------------- |
| 脆弱性         | NVD JSON Import + Version マッチング |
| SBOM Export | SPDX/CycloneDX 出力               |
| NOTICE      | ライセンステキスト集約生成                   |
| 認証強化        | LDAP / JWT                      |

---

## 14. 参考: AGENTS.md 役割

`AGENTS.md` は **機械可読メタ (コマンド / 生成禁止 / パス)**、本書は **背景と運用意図** を簡潔保持。差異が生じたら `AGENTS.md` を優先し、本書を更新する。

