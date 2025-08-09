/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { ScopeStatus } from './ScopeStatus';
import type { SupplierType } from './SupplierType';
import type { UsageRole } from './UsageRole';
/**
 * プロジェクトにおける特定 OSS バージョンの利用レコード
 */
export type ProjectUsage = {
    /**
     * 利用 ID
     */
    id: string;
    /**
     * プロジェクト ID
     */
    projectId: string;
    /**
     * OSSコンポーネント ID (利便性重複)
     */
    ossId: string;
    /**
     * OSS バージョン ID
     */
    ossVersionId: string;
    /**
     * 利用形態
     */
    usageRole: UsageRole;
    /**
     * スコープ判定
     */
    scopeStatus: ScopeStatus;
    /**
     * IN/OUT 判断理由メモ
     */
    inclusionNote?: string | null;
    /**
     * 直接依存なら true (間接は false)
     */
    directDependency: boolean;
    /**
     * 登録日時
     */
    addedAt: string;
    /**
     * スコープ判定更新日時
     */
    evaluatedAt?: string | null;
    /**
     * 判定実施ユーザ
     */
    evaluatedBy?: string | null;
    /**
     * 社内改変有無
     */
    modified: boolean;
    /**
     * 改変内容概要
     */
    modificationDescription?: string | null;
    /**
     * 供給形態
     */
    supplierType?: SupplierType | null;
    /**
     * フォーク元 URL (INTERNAL_FORK の場合)
     */
    forkOriginUrl?: string | null;
};

