/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { ScopePolicy } from '../models/ScopePolicy';
import type { ScopePolicyUpdateRequest } from '../models/ScopePolicyUpdateRequest';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class ScopePolicyService {
    /**
     * 現行スコープポリシー取得
     * @returns ScopePolicy OK
     * @throws ApiError
     */
    public static getScopePolicy(): CancelablePromise<ScopePolicy> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/scope/policy',
            errors: {
                401: `認証失敗`,
                403: `権限不足`,
            },
        });
    }
    /**
     * スコープポリシー更新 (管理者)
     * @returns ScopePolicy OK
     * @throws ApiError
     */
    public static updateScopePolicy({
        requestBody,
    }: {
        requestBody: ScopePolicyUpdateRequest,
    }): CancelablePromise<ScopePolicy> {
        return __request(OpenAPI, {
            method: 'PATCH',
            url: '/scope/policy',
            body: requestBody,
            mediaType: 'application/json',
            errors: {
                401: `認証失敗`,
                403: `権限不足`,
            },
        });
    }
}
