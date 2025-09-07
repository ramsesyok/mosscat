/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { Role } from './Role';
/**
 * システム利用ユーザー情報
 */
export type User = {
    /**
     * ユーザーID
     */
    id: string;
    /**
     * ログインID（ユニーク）
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
     * 付与ロール一覧
     */
    roles: Array<Role>;
    /**
     * 有効/無効フラグ
     */
    active: boolean;
    /**
     * 作成日時
     */
    createdAt: string;
    /**
     * 更新日時
     */
    updatedAt: string;
};

