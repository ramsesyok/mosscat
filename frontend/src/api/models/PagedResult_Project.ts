/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { Project } from './Project';
/**
 * プロジェクト一覧ページング結果
 */
export type PagedResult_Project = {
    /**
     * 結果アイテム配列
     */
    items?: Array<Project>;
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

