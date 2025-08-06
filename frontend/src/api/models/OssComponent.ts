/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { Layer } from './Layer';
import type { Tag } from './Tag';
import type { UsageRole } from './UsageRole';
/**
 * OSS の論理的名称（バージョン共通情報）
 */
export type OssComponent = {
    /**
     * OSSコンポーネント ID
     */
    id: string;
    /**
     * 表示名 / ユニーク（大文字小文字区別ポリシーは実装で決定）
     */
    name: string;
    /**
     * 検索用正規化名称（小文字化・記号除去）
     */
    normalizedName?: string;
    /**
     * 公式サイト URL
     */
    homepageUrl?: string | null;
    /**
     * ソースリポジトリ（Git 等）URL
     */
    repositoryUrl?: string | null;
    /**
     * 簡易説明 / メモ
     */
    description?: string | null;
    /**
     * 主言語（例: C, C++, Go, Java）
     */
    primaryLanguage?: string | null;
    /**
     * 技術レイヤ分類（複数指定可）
     */
    layers?: Array<Layer>;
    /**
     * 新規 ProjectUsage 作成時の初期 usageRole 推奨値
     */
    defaultUsageRole?: UsageRole | null;
    /**
     * 非推奨フラグ（新規利用抑止）
     */
    deprecated: boolean;
    /**
     * 作成日時
     */
    createdAt: string;
    /**
     * 更新日時
     */
    updatedAt: string;
    /**
     * 付与タグ一覧
     */
    tags?: Array<Tag>;
};

