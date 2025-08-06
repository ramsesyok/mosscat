/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { OssVersion } from '../models/OssVersion';
import type { OssVersionCreateRequest } from '../models/OssVersionCreateRequest';
import type { OssVersionUpdateRequest } from '../models/OssVersionUpdateRequest';
import type { PagedResult_OssVersion } from '../models/PagedResult_OssVersion';
import type { ReviewStatus } from '../models/ReviewStatus';
import type { ScopeStatus } from '../models/ScopeStatus';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class OssVersionsService {
    /**
     * 指定 OSS のバージョン一覧
     * @returns PagedResult_OssVersion OK
     * @throws ApiError
     */
    public static listOssVersions({
        ossId,
        page = 1,
        size = 50,
        reviewStatus,
        scopeStatus,
    }: {
        ossId: string,
        /**
         * 1 始まりのページ番号
         */
        page?: number,
        /**
         * 1ページ件数 (最大 200)
         */
        size?: number,
        reviewStatus?: ReviewStatus,
        scopeStatus?: ScopeStatus,
    }): CancelablePromise<PagedResult_OssVersion> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/oss/{ossId}/versions',
            path: {
                'ossId': ossId,
            },
            query: {
                'page': page,
                'size': size,
                'reviewStatus': reviewStatus,
                'scopeStatus': scopeStatus,
            },
            errors: {
                401: `認証失敗`,
                403: `権限不足`,
            },
        });
    }
    /**
     * バージョン追加
     * @returns OssVersion Created
     * @throws ApiError
     */
    public static createOssVersion({
        ossId,
        requestBody,
    }: {
        ossId: string,
        requestBody: OssVersionCreateRequest,
    }): CancelablePromise<OssVersion> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/oss/{ossId}/versions',
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
     * バージョン詳細
     * @returns OssVersion OK
     * @throws ApiError
     */
    public static getOssVersion({
        ossId,
        versionId,
    }: {
        ossId: string,
        versionId: string,
    }): CancelablePromise<OssVersion> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/oss/{ossId}/versions/{versionId}',
            path: {
                'ossId': ossId,
                'versionId': versionId,
            },
            errors: {
                401: `認証失敗`,
                403: `権限不足`,
                404: `リソースが存在しない`,
            },
        });
    }
    /**
     * バージョン更新
     * @returns OssVersion OK
     * @throws ApiError
     */
    public static updateOssVersion({
        ossId,
        versionId,
        requestBody,
    }: {
        ossId: string,
        versionId: string,
        requestBody: OssVersionUpdateRequest,
    }): CancelablePromise<OssVersion> {
        return __request(OpenAPI, {
            method: 'PATCH',
            url: '/oss/{ossId}/versions/{versionId}',
            path: {
                'ossId': ossId,
                'versionId': versionId,
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
     * バージョン削除 (論理/物理は実装方針による)
     * @returns void
     * @throws ApiError
     */
    public static deleteOssVersion({
        ossId,
        versionId,
    }: {
        ossId: string,
        versionId: string,
    }): CancelablePromise<void> {
        return __request(OpenAPI, {
            method: 'DELETE',
            url: '/oss/{ossId}/versions/{versionId}',
            path: {
                'ossId': ossId,
                'versionId': versionId,
            },
            errors: {
                401: `認証失敗`,
                403: `権限不足`,
                404: `リソースが存在しない`,
            },
        });
    }
}
