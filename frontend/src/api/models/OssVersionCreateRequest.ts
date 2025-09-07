/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
/**
 * バージョン作成リクエスト
 */
export type OssVersionCreateRequest = {
    /**
     * バージョン文字列
     */
    version: string;
    /**
     * リリース日
     */
    releaseDate?: string | null;
    /**
     * 生ライセンス式
     */
    licenseExpressionRaw?: string | null;
    /**
     * package-url
     */
    purl?: string | null;
    /**
     * CPE 配列
     */
    cpeList?: Array<string>;
    /**
     * アーカイブ SHA-256
     */
    hashSha256?: string | null;
};

