/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { ScopeStatus } from './ScopeStatus';
import type { SupplierType } from './SupplierType';
import type { UsageRole } from './UsageRole';
/**
 * プロジェクト利用更新リクエスト
 */
export type ProjectUsageUpdateRequest = {
    /**
     * 新しい OSS バージョン ID
     */
    ossVersionId?: string;
    /**
     * 利用形態
     */
    usageRole?: UsageRole;
    /**
     * 直接依存フラグ
     */
    directDependency?: boolean;
    /**
     * 理由メモ
     */
    inclusionNote?: string | null;
    /**
     * スコープ判定 (手動上書き)
     */
    scopeStatus?: ScopeStatus;
    /**
     * 社内改変有無
     */
    modified?: boolean;
    /**
     * 改変概要
     */
    modificationDescription?: string | null;
    /**
     * 供給形態
     */
    supplierType?: SupplierType | null;
    /**
     * フォーク元 URL
     */
    forkOriginUrl?: string | null;
};

