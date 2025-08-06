/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { ScopeStatus } from './ScopeStatus';
/**
 * スコープ判定更新リクエスト
 */
export type ScopeStatusUpdateRequest = {
    /**
     * 更新後ステータス
     */
    scopeStatus: ScopeStatus;
    /**
     * 理由コード (マスタ化想定)
     */
    reasonCode?: string | null;
    /**
     * 自由記述理由
     */
    reasonNote?: string | null;
};

