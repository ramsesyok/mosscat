/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { Role } from './Role';
/**
 * JWT に含める代表的クレーム例（参考情報）
 */
export type JwtClaims = {
    /**
     * ユーザーID
     */
    sub?: string;
    /**
     * ログインID
     */
    username?: string;
    /**
     * 付与ロール一覧
     */
    roles?: Array<Role>;
    /**
     * 有効期限 (UNIX epoch seconds)
     */
    exp?: number;
    /**
     * 発行時刻 (UNIX epoch seconds)
     */
    iat?: number;
};

