/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { SupplierType } from './SupplierType';
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
    /**
     * 社内改変有無
     */
    modified?: boolean;
    /**
     * 改変概要
     */
    modificationDescription?: string | null;
    /**
     * 供給形態
     */
    supplierType?: SupplierType | null;
    /**
     * フォーク元 URL
     */
    forkOriginUrl?: string | null;
};

