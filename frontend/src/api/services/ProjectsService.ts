/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { PagedResult_Project } from '../models/PagedResult_Project';
import type { Project } from '../models/Project';
import type { ProjectCreateRequest } from '../models/ProjectCreateRequest';
import type { ProjectUpdateRequest } from '../models/ProjectUpdateRequest';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class ProjectsService {
    /**
     * プロジェクト一覧
     * @returns PagedResult_Project OK
     * @throws ApiError
     */
    public static listProjects({
        page = 1,
        size = 50,
        code,
        name,
    }: {
        /**
         * 1 始まりのページ番号
         */
        page?: number,
        /**
         * 1ページ件数 (最大 200)
         */
        size?: number,
        code?: string,
        name?: string,
    }): CancelablePromise<PagedResult_Project> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/projects',
            query: {
                'page': page,
                'size': size,
                'code': code,
                'name': name,
            },
            errors: {
                401: `認証失敗`,
                403: `権限不足`,
            },
        });
    }
    /**
     * プロジェクト作成
     * @returns Project Created
     * @throws ApiError
     */
    public static createProject({
        requestBody,
    }: {
        requestBody: ProjectCreateRequest,
    }): CancelablePromise<Project> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/projects',
            body: requestBody,
            mediaType: 'application/json',
            errors: {
                401: `認証失敗`,
                403: `権限不足`,
            },
        });
    }
    /**
     * プロジェクト詳細
     * @returns Project OK
     * @throws ApiError
     */
    public static getProject({
        projectId,
    }: {
        projectId: string,
    }): CancelablePromise<Project> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/projects/{projectId}',
            path: {
                'projectId': projectId,
            },
            errors: {
                401: `認証失敗`,
                403: `権限不足`,
                404: `リソースが存在しない`,
            },
        });
    }
    /**
     * プロジェクト更新
     * @returns Project OK
     * @throws ApiError
     */
    public static updateProject({
        projectId,
        requestBody,
    }: {
        projectId: string,
        requestBody: ProjectUpdateRequest,
    }): CancelablePromise<Project> {
        return __request(OpenAPI, {
            method: 'PATCH',
            url: '/projects/{projectId}',
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
     * プロジェクト削除 (論理予定)
     * @returns void
     * @throws ApiError
     */
    public static deleteProject({
        projectId,
    }: {
        projectId: string,
    }): CancelablePromise<void> {
        return __request(OpenAPI, {
            method: 'DELETE',
            url: '/projects/{projectId}',
            path: {
                'projectId': projectId,
            },
            errors: {
                401: `認証失敗`,
                403: `権限不足`,
                404: `リソースが存在しない`,
            },
        });
    }
}
