/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { Role } from './Role';
/**
 * ユーザー更新リクエスト（部分）
 */
export type UserUpdateRequest = {
    /**
     * 表示名
     */
    displayName?: string;
    /**
     * メールアドレス
     */
    email?: string | null;
    /**
     * パスワード再設定用（任意）
     */
    password?: string;
    /**
     * 付与ロール一覧（管理者のみ変更可）
     */
    roles?: Array<Role>;
    /**
     * 有効/無効切替
     */
    active?: boolean;
};

