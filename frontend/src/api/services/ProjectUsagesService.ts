/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { PagedResult_ProjectUsage } from '../models/PagedResult_ProjectUsage';
import type { ProjectUsage } from '../models/ProjectUsage';
import type { ProjectUsageCreateRequest } from '../models/ProjectUsageCreateRequest';
import type { ProjectUsageUpdateRequest } from '../models/ProjectUsageUpdateRequest';
import type { ScopeStatus } from '../models/ScopeStatus';
import type { ScopeStatusUpdateRequest } from '../models/ScopeStatusUpdateRequest';
import type { UsageRole } from '../models/UsageRole';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class ProjectUsagesService {
    /**
     * プロジェクト中利用 OSS 一覧
     * @returns PagedResult_ProjectUsage OK
     * @throws ApiError
     */
    public static listProjectUsages({
        projectId,
        page = 1,
        size = 50,
        scopeStatus,
        usageRole,
        direct,
    }: {
        projectId: string,
        /**
         * 1 始まりのページ番号
         */
        page?: number,
        /**
         * 1ページ件数 (最大 200)
         */
        size?: number,
        scopeStatus?: ScopeStatus,
        usageRole?: UsageRole,
        /**
         * 直接依存のみ true
         */
        direct?: boolean,
    }): CancelablePromise<PagedResult_ProjectUsage> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/projects/{projectId}/usages',
            path: {
                'projectId': projectId,
            },
            query: {
                'page': page,
                'size': size,
                'scopeStatus': scopeStatus,
                'usageRole': usageRole,
                'direct': direct,
            },
            errors: {
                401: `認証失敗`,
                403: `権限不足`,
            },
        });
    }
    /**
     * プロジェクト利用追加
     * @returns ProjectUsage Created
     * @throws ApiError
     */
    public static createProjectUsage({
        projectId,
        requestBody,
    }: {
        projectId: string,
        requestBody: ProjectUsageCreateRequest,
    }): CancelablePromise<ProjectUsage> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/projects/{projectId}/usages',
            path: {
                'projectId': projectId,
            },
            body: requestBody,
            mediaType: 'application/json',
            errors: {
                401: `認証失敗`,
                403: `権限不足`,
            },
        });
    }
    /**
     * 利用情報更新
     * @returns ProjectUsage OK
     * @throws ApiError
     */
    public static updateProjectUsage({
        projectId,
        usageId,
        requestBody,
    }: {
        projectId: string,
        usageId: string,
        requestBody: ProjectUsageUpdateRequest,
    }): CancelablePromise<ProjectUsage> {
        return __request(OpenAPI, {
            method: 'PATCH',
            url: '/projects/{projectId}/usages/{usageId}',
            path: {
                'projectId': projectId,
                'usageId': usageId,
            },
            body: requestBody,
            mediaType: 'application/json',
            errors: {
                401: `認証失敗`,
                403: `権限不足`,
            },
        });
    }
    /**
     * 利用削除
     * @returns void
     * @throws ApiError
     */
    public static deleteProjectUsage({
        projectId,
        usageId,
    }: {
        projectId: string,
        usageId: string,
    }): CancelablePromise<void> {
        return __request(OpenAPI, {
            method: 'DELETE',
            url: '/projects/{projectId}/usages/{usageId}',
            path: {
                'projectId': projectId,
                'usageId': usageId,
            },
            errors: {
                401: `認証失敗`,
                403: `権限不足`,
            },
        });
    }
    /**
     * スコープ判定更新
     * @returns ProjectUsage OK
     * @throws ApiError
     */
    public static updateProjectUsageScope({
        projectId,
        usageId,
        requestBody,
    }: {
        projectId: string,
        usageId: string,
        requestBody: ScopeStatusUpdateRequest,
    }): CancelablePromise<ProjectUsage> {
        return __request(OpenAPI, {
            method: 'PATCH',
            url: '/projects/{projectId}/usages/{usageId}/scope',
            path: {
                'projectId': projectId,
                'usageId': usageId,
            },
            body: requestBody,
            mediaType: 'application/json',
            errors: {
                401: `認証失敗`,
                403: `権限不足`,
            },
        });
    }
}
