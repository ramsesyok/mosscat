/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { OssVersion } from './OssVersion';
/**
 * OSS バージョンページング結果
 */
export type PagedResult_OssVersion = {
    /**
     * 結果アイテム配列
     */
    items?: Array<OssVersion>;
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

