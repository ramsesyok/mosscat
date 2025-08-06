/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { PagedResult_User } from '../models/PagedResult_User';
import type { Role } from '../models/Role';
import type { User } from '../models/User';
import type { UserCreateRequest } from '../models/UserCreateRequest';
import type { UserUpdateRequest } from '../models/UserUpdateRequest';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class UsersService {
    /**
     * 現在ログイン中ユーザー情報取得
     * @returns User OK
     * @throws ApiError
     */
    public static getCurrentUser(): CancelablePromise<User> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/me',
            errors: {
                401: `認証失敗`,
            },
        });
    }
    /**
     * ユーザー一覧
     * @returns PagedResult_User OK
     * @throws ApiError
     */
    public static listUsers({
        page = 1,
        size = 50,
        username,
        role,
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
         * 部分一致検索
         */
        username?: string,
        /**
         * ロールで絞り込み
         */
        role?: Role,
    }): CancelablePromise<PagedResult_User> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/users',
            query: {
                'page': page,
                'size': size,
                'username': username,
                'role': role,
            },
            errors: {
                401: `認証失敗`,
                403: `権限不足`,
            },
        });
    }
    /**
     * ユーザー作成
     * @returns User Created
     * @throws ApiError
     */
    public static createUser({
        requestBody,
    }: {
        requestBody: UserCreateRequest,
    }): CancelablePromise<User> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/users',
            body: requestBody,
            mediaType: 'application/json',
            errors: {
                401: `認証失敗`,
                403: `権限不足`,
            },
        });
    }
    /**
     * ユーザー詳細
     * @returns User OK
     * @throws ApiError
     */
    public static getUser({
        userId,
    }: {
        userId: string,
    }): CancelablePromise<User> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/users/{userId}',
            path: {
                'userId': userId,
            },
            errors: {
                401: `認証失敗`,
                403: `権限不足`,
                404: `リソースが存在しない`,
            },
        });
    }
    /**
     * ユーザー更新
     * @returns User OK
     * @throws ApiError
     */
    public static updateUser({
        userId,
        requestBody,
    }: {
        userId: string,
        requestBody: UserUpdateRequest,
    }): CancelablePromise<User> {
        return __request(OpenAPI, {
            method: 'PATCH',
            url: '/users/{userId}',
            path: {
                'userId': userId,
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
     * ユーザー削除 (論理/物理は実装次第)
     * @returns void
     * @throws ApiError
     */
    public static deleteUser({
        userId,
    }: {
        userId: string,
    }): CancelablePromise<void> {
        return __request(OpenAPI, {
            method: 'DELETE',
            url: '/users/{userId}',
            path: {
                'userId': userId,
            },
            errors: {
                401: `認証失敗`,
                403: `権限不足`,
                404: `リソースが存在しない`,
            },
        });
    }
}
