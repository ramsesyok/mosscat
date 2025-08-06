/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
/**
 * スコープ自動判定ポリシー設定
 */
export type ScopePolicy = {
    /**
     * ポリシー ID
     */
    id?: string;
    /**
     * RUNTIME_REQUIRED を自動 IN_SCOPE にするフラグ
     */
    runtimeRequiredDefaultInScope?: boolean;
    /**
     * SERVER_ENV を IN_SCOPE とみなすか
     */
    serverEnvIncluded?: boolean;
    /**
     * INTERNAL_FORK を自動 IN_SCOPE 化するか
     */
    autoMarkForksInScope?: boolean;
    /**
     * 最終更新日時
     */
    updatedAt?: string;
    /**
     * 更新ユーザ
     */
    updatedBy?: string;
};

