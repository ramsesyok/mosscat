/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { ReviewStatus } from './ReviewStatus';
import type { ScopeStatus } from './ScopeStatus';
import type { SupplierType } from './SupplierType';
/**
 * バージョン更新リクエスト（部分）
 */
export type OssVersionUpdateRequest = {
    /**
     * リリース日
     */
    releaseDate?: string | null;
    /**
     * 生ライセンス式
     */
    licenseExpressionRaw?: string | null;
    /**
     * 確定ライセンス式
     */
    licenseConcluded?: string | null;
    /**
     * package-url
     */
    purl?: string | null;
    /**
     * CPE 配列
     */
    cpeList?: Array<string>;
    /**
     * SHA-256 ハッシュ
     */
    hashSha256?: string | null;
    /**
     * 改変有無
     */
    modified?: boolean;
    /**
     * 改変概要
     */
    modificationDescription?: string | null;
    /**
     * レビュー状態
     */
    reviewStatus?: ReviewStatus;
    /**
     * 推奨スコープ判定
     */
    scopeStatus?: ScopeStatus;
    /**
     * 供給形態
     */
    supplierType?: SupplierType | null;
    /**
     * フォーク元 URL
     */
    forkOriginUrl?: string | null;
};

