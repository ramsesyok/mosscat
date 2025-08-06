/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { OssComponent } from '../models/OssComponent';
import type { OssComponentCreateRequest } from '../models/OssComponentCreateRequest';
import type { OssComponentUpdateRequest } from '../models/OssComponentUpdateRequest';
import type { PagedResult_OssComponent } from '../models/PagedResult_OssComponent';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class OssService {
    /**
     * OSSコンポーネント一覧取得
     * @returns PagedResult_OssComponent OK
     * @throws ApiError
     */
    public static listOssComponents({
        page = 1,
        size = 50,
        sort,
        name,
        layers,
        tag,
        inScopeOnly,
    }: {
        /**
         * 1 始まりのページ番号
         */
        page?: number,
        /**
         * 1ページ件数 (最大 200)
         */
        size?: number,
        /**
         * ソート指定 (例: name,asc / createdAt,desc)
         */
        sort?: string,
        /**
         * 部分一致
         */
        name?: string,
        /**
         * カンマ区切り Layer フィルタ (例 LIB,DB)
         */
        layers?: string,
        /**
         * タグ名 (正確一致)
         */
        tag?: string,
        /**
         * true の場合 IN_SCOPE のバージョンを一つ以上持つもののみ
         */
        inScopeOnly?: boolean,
    }): CancelablePromise<PagedResult_OssComponent> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/oss',
            query: {
                'page': page,
                'size': size,
                'sort': sort,
                'name': name,
                'layers': layers,
                'tag': tag,
                'inScopeOnly': inScopeOnly,
            },
            errors: {
                401: `認証失敗`,
                403: `権限不足`,
            },
        });
    }
    /**
     * OSSコンポーネント作成
     * @returns OssComponent Created
     * @throws ApiError
     */
    public static createOssComponent({
        requestBody,
    }: {
        requestBody: OssComponentCreateRequest,
    }): CancelablePromise<OssComponent> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/oss',
            body: requestBody,
            mediaType: 'application/json',
            errors: {
                400: `不正なリクエスト`,
                401: `認証失敗`,
                403: `権限不足`,
            },
        });
    }
    /**
     * OSSコンポーネント詳細
     * @returns OssComponent OK
     * @throws ApiError
     */
    public static getOssComponent({
        ossId,
    }: {
        ossId: string,
    }): CancelablePromise<OssComponent> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/oss/{ossId}',
            path: {
                'ossId': ossId,
            },
            errors: {
                401: `認証失敗`,
                403: `権限不足`,
                404: `リソースが存在しない`,
            },
        });
    }
    /**
     * OSSコンポーネント更新 (部分)
     * @returns OssComponent OK
     * @throws ApiError
     */
    public static updateOssComponent({
        ossId,
        requestBody,
    }: {
        ossId: string,
        requestBody: OssComponentUpdateRequest,
    }): CancelablePromise<OssComponent> {
        return __request(OpenAPI, {
            method: 'PATCH',
            url: '/oss/{ossId}',
            path: {
                'ossId': ossId,
            },
            body: requestBody,
            mediaType: 'application/json',
            errors: {
                401: `認証失敗`,
                403: `権限不足`,
                404: `リソースが存在しない`,
            },
        });
    }
    /**
     * OSSコンポーネントを非推奨 (deprecated=true) に設定
     * @returns void
     * @throws ApiError
     */
    public static deprecateOssComponent({
        ossId,
    }: {
        ossId: string,
    }): CancelablePromise<void> {
        return __request(OpenAPI, {
            method: 'DELETE',
            url: '/oss/{ossId}',
            path: {
                'ossId': ossId,
            },
            errors: {
                401: `認証失敗`,
                403: `権限不足`,
                404: `リソースが存在しない`,
            },
        });
    }
}
