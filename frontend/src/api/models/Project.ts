/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
/**
 * プロジェクト（納品単位）
 */
export type Project = {
    /**
     * プロジェクト ID
     */
    id: string;
    /**
     * 社内識別コード（ユニーク）
     */
    projectCode: string;
    /**
     * プロジェクト名称
     */
    name: string;
    /**
     * 担当部署
     */
    department?: string | null;
    /**
     * 責任者 / 担当者名
     */
    manager?: string | null;
    /**
     * 納品予定日
     */
    deliveryDate?: string | null;
    /**
     * 説明 / 備考
     */
    description?: string | null;
    /**
     * 紐付く OSS 利用件数集計
     */
    ossUsageCount?: number;
    /**
     * 作成日時
     */
    createdAt: string;
    /**
     * 更新日時
     */
    updatedAt: string;
};

