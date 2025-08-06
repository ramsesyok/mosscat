/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
/**
 * ログイン成功時のレスポンス
 */
export type LoginResponse = {
    /**
     * JWT アクセストークン。Authorization: Bearer にセットする
     */
    accessToken: string;
    /**
     * トークンの有効期限（秒）
     */
    expiresIn: number;
};

