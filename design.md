# OSSカタログ Webアプリ設計書

## 0. サービス概要（目的サマリ）

**目的:** 社内プロジェクトが *納品対象となる OSS コンポーネント* を統一管理し、バージョン・ライセンス・改変有無・利用形態 (UsageRole)・スコープ判定 (IN/OUT/REVIEW) を登録・検索・エクスポートできる基盤を提供する。

| 観点        | 要約                                                                      |
| --------- | ----------------------------------------------------------------------- |
| 主ユースケース   | プロジェクトで利用する OSS を登録し、納品一覧（IN\_SCOPE）だけを抽出 (将来: SPDX / NOTICE)           |
| ドメイン中核    | **OssComponent (親)** と **OssVersion (子)**, **ProjectUsage (結合 + スコープ)** |
| スコープ分類    | UsageRole + ポリシー (ScopePolicy) により自動初期判定 / 人手で上書き                       |
| 除外対象      | 開発専用 (DEV/BUILD/TEST), サーバ環境のみ (SERVER\_ENV), ポリシー設定で制御                 |
| 非機能       | 単一ノード / オフライン / PostgreSQL 本番 + SQLite 簡易運用                             |
| フロントエンド概要 | Vue 3 SPA（Vuetify UI + Pinia 状態管理）による管理画面構築                             |
| 将来拡張      | 脆弱性(NVD)・SBOM Import・NOTICE 生成・LDAP/JWT 認証                              |

---

## 1. アーキテクチャ構成（バックエンド + フロントエンド）

### 🔧 バックエンド構成

| 層                             | 役割                                   | 実装                               |
| ----------------------------- | ------------------------------------ | -------------------------------- |
| API (Echo)                    | OpenAPI に沿った Handler                 | oapi-codegen server stub + 手書き実装 |
| DTO/Gen                       | 自動生成型                                | `internal/api/gen` (編集禁止)        |
| Domain (model/service/policy) | ビジネスロジック / 正規化 / スコープ判定              | 手書き                              |
| Repository                    | DB アクセス (SQL)                        | `infra/repository`               |
| Infra                         | DB接続 / Tx / Migration                | `infra/db`, `migration`          |
| Shared                        | errors / auth / pagination / logging | `pkg/*`                          |

### 🌐 フロントエンド構成

| 分類         | 内容                                |
| ---------- | --------------------------------- |
| フレームワーク    | Vue 3 (script setup) + TypeScript |
| UI         | Vuetify 3                         |
| 状態管理       | Pinia                             |
| 型生成        | openapi-typescript-codegen        |
| 認証処理       | Basic / JWT 自動切替                  |
| 多言語対応      | vue-i18n                          |
| エクスポート     | CSV / SPDX-JSON（プロジェクト単位）         |
| レイアウト      | ナビドロワー + AppBar + RouterView      |
| 生成禁止ディレクトリ | `src/api`（型定義/クライアント生成）           |

---

## 2. ディレクトリ構成

```
mosscat/
├─ frontend/
│  ├─ src/
│  │  ├─ api/                # openapi 型安全クライアント (generated)
│  │  ├─ components/
│  │  ├─ pages/
│  │  ├─ stores/
│  │  ├─ router/
│  │  ├─ plugins/
│  │  ├─ locales/
│  │  └─ App.vue
│  ├─ public/
│  ├─ vite.config.ts
│  └─ scripts/generate-client.mjs
├─ cmd/api/main.go
├─ internal/api/openapi.yaml
├─ internal/api/gen/
├─ internal/api/handler/
├─ internal/domain/
├─ internal/infra/
├─ pkg/
├─ AGENTS.md
└─ design.md  ← ★このファイル
```

---

## 3. 主要機能とコンポーネント

### 🧠 スコープ初期判定ロジック（Go）

```go
switch usageRole {
case "BUILD_ONLY", "DEV_ONLY", "TEST_ONLY": return OUT_SCOPE
case "SERVER_ENV":  if policy.ServerEnvIncluded { return IN_SCOPE } else { return OUT_SCOPE }
case "RUNTIME_REQUIRED": if policy.RuntimeRequiredDefaultInScope { return IN_SCOPE } else { return REVIEW_NEEDED }
default: return IN_SCOPE // Bundled / Linked 系
}
```

### 🖥️ フロントエンド UI（Vue）

* `OssListPage.vue`: OSS 一覧、検索バー、詳細ダイアログ連携
* `OssDetailDialog.vue`: OSSのフォーム編集 + バージョン一覧タブ
* `ProjectListPage.vue`: プロジェクト一覧 + エクスポートメニュー
* `ScopeStatusChip.vue`: IN/OUT/REVIEW に応じた色 + Tooltip
* `TagAutocomplete.vue`: APIベースタグ候補 + 新規追加

---

## 4. エクスポートフロー（SPA）

1. ユーザが「Export > CSV」を選択
2. `ProjectApi.exportProjectArtifacts()` を呼び出し
3. Blob をダウンロード
4. 完了後スナックバー表示、長時間はプログレスバー表示

---

## 5. テストポリシー

| 種別           | 内容                                      |
| ------------ | --------------------------------------- |
| Go Unit      | ScopeDecider / Service の単体テスト           |
| Go API       | 200/400/404 のエンドポイントテスト                 |
| Go Repo      | SQL 実装の CRUD / ページング                    |
| Vue Frontend | vitest + vue-test-utils によるモック/スナップショット |
| Migration    | up/down 成功確認                            |

---

## 6. コマンド・自動生成方針

### Goサイド

| 対象     | 条件                     | コマンド                   |
| ------ | ---------------------- | ---------------------- |
| types  | components/schemas 変更時 | `make generate-types`  |
| server | paths/operationId 変更時  | `make generate-server` |

### Vueサイド

```bash
npm run gen:client
```

```js
// scripts/generate-client.mjs
await generate({
  input: '../backend/internal/api/openapi.yaml',
  output: 'src/api',
  httpClient: 'axios',
  useOptions: true,
  useUnionTypes: true,
});
```

---

## 7. AGENTS.md 方針との関係

| 対象          | 禁止フォルダ              | 注意点                     |
| ----------- | ------------------- | ----------------------- |
| Go生成ファイル    | `internal/api/gen/` | 編集禁止、CIで差分チェック          |
| Vue生成クライアント | `frontend/src/api/` | 自動生成対象。コード補完対象外、AI 編集禁止 |
| カバレッジ       | `go test`/`vitest`  | 最低60%目標（将来）             |

---

## 8. 今後の拡張項目（全体）

| 項目       | 説明                                 |
| -------- | ---------------------------------- |
| SBOM     | SPDX / CycloneDX エクスポート機能          |
| NOTICE生成 | ライセンス文章集約機能                        |
| 脆弱性      | NVD JSON インポート → バージョン一致判定         |
| 認証強化     | Basic → JWT / LDAP 対応              |
| パフォーマンス  | OSS一覧・プロジェクト一覧の仮想スクロール（5,000行以上対応） |
| アクセシビリティ | Vuetify準拠 + キーボード対応必須              |
| PWA      | オフライン動作検討（ドキュメント用途）                |

---

## 9. Definition of Done（共通）

| 項目        | 条件                                         |
| --------- | ------------------------------------------ |
| コード整合     | 自動生成差分がない (`make generate` / `gen:client`) |
| テスト通過     | Go + Vue の主要処理の単体テストを通過                    |
| Lint適合    | `go vet`, ESLint による静的解析 OK                |
| エラーハンドリング | OpenAPI/Problem形式 (Go) + Snackbar表現 (Vue)  |
| UI反映      | スコープ色、タグ作成、エクスポート、認証切替などが動作確認済             |

