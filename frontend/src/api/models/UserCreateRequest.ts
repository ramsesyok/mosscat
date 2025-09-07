/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { Role } from './Role';
/**
 * ユーザー作成リクエスト
 */
export type UserCreateRequest = {
    /**
     * ログインID
     */
    username: string;
    /**
     * 表示名
     */
    displayName?: string;
    /**
     * メールアドレス
     */
    email?: string | null;
    /**
     * 初期パスワード（JWTでも登録時のみ使用可）
     */
    password?: string;
    /**
     * 付与するロール一覧
     */
    roles: Array<Role>;
    /**
     * 有効化フラグ
     */
    active?: boolean;
};

