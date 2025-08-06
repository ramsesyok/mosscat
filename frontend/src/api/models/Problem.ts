/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
/**
 * RFC 9457 / RFC 7807 型エラー応答ボディ
 */
export type Problem = {
    /**
     * 問題タイプ識別 URI（拡張分類用）
     */
    type?: string | null;
    /**
     * エラーの概要メッセージ（人間可読）
     */
    title: string;
    /**
     * HTTP ステータスコード
     */
    status: number;
    /**
     * 追加詳細メッセージ
     */
    detail?: string | null;
    /**
     * エラーが発生した具体的インスタンス URI（トレースID等）
     */
    instance?: string | null;
    /**
     * アプリケーション独自エラーコード
     */
    code?: string | null;
    /**
     * フィールド単位バリデーションエラー配列
     */
    errors?: Array<{
        /**
         * エラーが発生したフィールド名（JSON Pointer など）
         */
        field?: string;
        /**
         * フィールドに対するエラーメッセージ
         */
        message?: string;
    }>;
};

