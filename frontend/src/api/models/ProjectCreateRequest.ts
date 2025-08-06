/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
/**
 * プロジェクト作成リクエスト
 */
export type ProjectCreateRequest = {
    /**
     * ユニークなプロジェクトコード
     */
    projectCode: string;
    /**
     * 名称
     */
    name: string;
    /**
     * 部署
     */
    department?: string | null;
    /**
     * 責任者
     */
    manager?: string | null;
    /**
     * 納品日
     */
    deliveryDate?: string | null;
    /**
     * 説明 / 備考
     */
    description?: string | null;
};

