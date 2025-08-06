/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
/**
 * スコープポリシー更新リクエスト
 */
export type ScopePolicyUpdateRequest = {
    /**
     * RUNTIME_REQUIRED 自動 IN_SCOPE フラグ
     */
    runtimeRequiredDefaultInScope?: boolean;
    /**
     * SERVER_ENV を IN_SCOPE とするか
     */
    serverEnvIncluded?: boolean;
    /**
     * INTERNAL_FORK 自動 IN_SCOPE フラグ
     */
    autoMarkForksInScope?: boolean;
};

