/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { ProjectUsage } from './ProjectUsage';
/**
 * プロジェクト利用レコードのページング結果
 */
export type PagedResult_ProjectUsage = {
    /**
     * 結果アイテム配列
     */
    items?: Array<ProjectUsage>;
    /**
     * 現在ページ
     */
    page?: number;
    /**
     * ページサイズ
     */
    size?: number;
    /**
     * 総件数
     */
    total?: number;
};

