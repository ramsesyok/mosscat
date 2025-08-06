/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { OssComponent } from './OssComponent';
/**
 * OSSコンポーネントページング結果
 */
export type PagedResult_OssComponent = {
    /**
     * 結果アイテム配列
     */
    items?: Array<OssComponent>;
    /**
     * 現在ページ (1 始まり)
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

