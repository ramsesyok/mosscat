/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { UsageRole } from './UsageRole';
/**
 * プロジェクト利用作成リクエスト
 */
export type ProjectUsageCreateRequest = {
    /**
     * OSSコンポーネント ID
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
     * 直接依存なら true
     */
    directDependency?: boolean;
    /**
     * 初期理由メモ
     */
    inclusionNote?: string | null;
};

