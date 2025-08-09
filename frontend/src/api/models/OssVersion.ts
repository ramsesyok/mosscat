/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { ReviewStatus } from './ReviewStatus';
import type { ScopeStatus } from './ScopeStatus';
/**
 * 個別バージョン情報
 */
export type OssVersion = {
    /**
     * バージョン ID
     */
    id: string;
    /**
     * 親 OSSコンポーネント ID
     */
    ossId: string;
    /**
     * バージョン文字列 (SemVer 互換推奨)
     */
    version: string;
    /**
     * 上流リリース日
     */
    releaseDate?: string | null;
    /**
     * SPDX 互換ライセンス式（上流取得時点）
     */
    licenseExpressionRaw?: string | null;
    /**
     * 社内審査確定ライセンス式
     */
    licenseConcluded?: string | null;
    /**
     * package-url (例: pkg:maven/...)
     */
    purl?: string | null;
    /**
     * CPE 文字列配列（脆弱性紐付け用）
     */
    cpeList?: Array<string>;
    /**
     * 配布アーカイブ等の SHA-256 ハッシュ
     */
    hashSha256?: string | null;
    /**
     * レビュー状態
     */
    reviewStatus: ReviewStatus;
    /**
     * 最終レビュー日時
     */
    lastReviewedAt?: string | null;
    /**
     * バージョン単位の推奨スコープ判定（参考）
     */
    scopeStatus: ScopeStatus;
    /**
     * 作成日時
     */
    createdAt: string;
    /**
     * 更新日時
     */
    updatedAt: string;
};

