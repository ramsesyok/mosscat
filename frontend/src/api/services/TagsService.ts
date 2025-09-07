/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { Tag } from '../models/Tag';
import type { TagCreateRequest } from '../models/TagCreateRequest';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class TagsService {
    /**
     * タグ一覧
     * @returns Tag OK
     * @throws ApiError
     */
    public static listTags(): CancelablePromise<Array<Tag>> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/tags',
            errors: {
                401: `認証失敗`,
                403: `権限不足`,
            },
        });
    }
    /**
     * タグ作成
     * @returns Tag Created
     * @throws ApiError
     */
    public static createTag({
        requestBody,
    }: {
        requestBody: TagCreateRequest,
    }): CancelablePromise<Tag> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/tags',
            body: requestBody,
            mediaType: 'application/json',
            errors: {
                401: `認証失敗`,
                403: `権限不足`,
            },
        });
    }
    /**
     * タグ削除
     * @returns void
     * @throws ApiError
     */
    public static deleteTag({
        tagId,
    }: {
        tagId: string,
    }): CancelablePromise<void> {
        return __request(OpenAPI, {
            method: 'DELETE',
            url: '/tags/{tagId}',
            path: {
                'tagId': tagId,
            },
            errors: {
                401: `認証失敗`,
                403: `権限不足`,
            },
        });
    }
}
