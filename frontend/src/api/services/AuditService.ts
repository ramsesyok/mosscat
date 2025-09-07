/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class AuditService {
    /**
     * 監査ログ簡易検索 (Phase1簡易)
     * @returns any OK
     * @throws ApiError
     */
    public static searchAuditLogs({
        entityType,
        entityId,
        from,
        to,
    }: {
        entityType?: string,
        entityId?: string,
        from?: string,
        to?: string,
    }): CancelablePromise<{
        items?: Array<{
            id?: string;
            entityType?: string;
            entityId?: string;
            action?: string;
            at?: string;
            user?: string;
            summary?: string;
        }>;
    }> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/audit',
            query: {
                'entityType': entityType,
                'entityId': entityId,
                'from': from,
                'to': to,
            },
            errors: {
                401: `認証失敗`,
                403: `権限不足`,
            },
        });
    }
}
