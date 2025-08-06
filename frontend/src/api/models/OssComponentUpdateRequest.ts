/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { Layer } from './Layer';
import type { UsageRole } from './UsageRole';
/**
 * OSSコンポーネント更新リクエスト（部分）
 */
export type OssComponentUpdateRequest = {
    /**
     * 新しい名称
     */
    name?: string;
    /**
     * 公式サイト URL
     */
    homepageUrl?: string | null;
    /**
     * リポジトリ URL
     */
    repositoryUrl?: string | null;
    /**
     * 説明 / 備考
     */
    description?: string | null;
    /**
     * 主言語
     */
    primaryLanguage?: string | null;
    /**
     * 技術レイヤ配列
     */
    layers?: Array<Layer>;
    /**
     * デフォルト usageRole
     */
    defaultUsageRole?: UsageRole | null;
    /**
     * 非推奨フラグ
     */
    deprecated?: boolean;
    /**
     * 置換後のタグ ID 配列
     */
    tagIds?: Array<string>;
};

