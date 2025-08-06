/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class ExportService {
    /**
     * プロジェクト納品用エクスポート (プレースホルダ)
     * @returns any Export ファイル (現状 JSON/CSV テキスト)
     * @throws ApiError
     */
    public static exportProjectArtifacts({
        projectId,
        format,
        scopes,
    }: {
        projectId: string,
        format: 'csv' | 'spdx-json',
        scopes?: string,
    }): CancelablePromise<{
        placeholder?: string;
    }> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/projects/{projectId}/export',
            path: {
                'projectId': projectId,
            },
            query: {
                'format': format,
                'scopes': scopes,
            },
            errors: {
                401: `認証失敗`,
                403: `権限不足`,
            },
        });
    }
}
