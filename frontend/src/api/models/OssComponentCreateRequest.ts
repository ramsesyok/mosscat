/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { Layer } from './Layer';
import type { UsageRole } from './UsageRole';
/**
 * OSSコンポーネント作成リクエスト
 */
export type OssComponentCreateRequest = {
    /**
     * OSS 表示名
     */
    name: string;
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
     * 初期付与するタグ ID 群
     */
    tagIds?: Array<string>;
};

