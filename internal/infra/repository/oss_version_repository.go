package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lib/pq"

	"github.com/ramsesyok/mosscat/internal/domain/model"
	domrepo "github.com/ramsesyok/mosscat/internal/domain/repository"
)

// OssVersionRepository は domrepo.OssVersionRepository の実装。
type OssVersionRepository struct {
	DB *sql.DB
}

var _ domrepo.OssVersionRepository = (*OssVersionRepository)(nil)

// Search は指定された OSS コンポーネントのバージョン一覧を返す。
func (r *OssVersionRepository) Search(ctx context.Context, f domrepo.OssVersionFilter) ([]model.OssVersion, int, error) {
	var args []any
	wheres := []string{"oss_id = ?"}
	args = append(args, f.OssID)
	if f.ReviewStatus != "" {
		wheres = append(wheres, "review_status = ?")
		args = append(args, f.ReviewStatus)
	}
	if f.ScopeStatus != "" {
		wheres = append(wheres, "scope_status = ?")
		args = append(args, f.ScopeStatus)
	}
	whereSQL := whereClause(wheres)

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM oss_versions %s", whereSQL)
	row := r.DB.QueryRowContext(ctx, countQuery, args...)
	var total int
	if err := row.Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (f.Page - 1) * f.Size
	listQuery := fmt.Sprintf(`SELECT id, oss_id, version, release_date, license_expression_raw, license_concluded, purl, cpe_list, hash_sha256, review_status, last_reviewed_at, scope_status, created_at, updated_at FROM oss_versions %s ORDER BY created_at DESC LIMIT ? OFFSET ?`, whereSQL)
	argsWithLimit := append(args, f.Size, offset)
	rows, err := r.DB.QueryContext(ctx, listQuery, argsWithLimit...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var versions []model.OssVersion
	for rows.Next() {
		var v model.OssVersion
		var releaseDate sql.NullTime
		var licenseRaw, licenseConc, purl, hash sql.NullString
		var lastReviewed sql.NullTime
		var cpeList pq.StringArray
		if err := rows.Scan(&v.ID, &v.OssID, &v.Version, &releaseDate, &licenseRaw, &licenseConc, &purl, &cpeList, &hash, &v.ReviewStatus, &lastReviewed, &v.ScopeStatus, &v.CreatedAt, &v.UpdatedAt); err != nil {
			return nil, 0, err
		}
		v.ReleaseDate = timePtr(releaseDate)
		v.LicenseExpressionRaw = strPtr(licenseRaw)
		v.LicenseConcluded = strPtr(licenseConc)
		v.Purl = strPtr(purl)
		v.CpeList = []string(cpeList)
		v.HashSha256 = strPtr(hash)
		v.LastReviewedAt = timePtr(lastReviewed)
		versions = append(versions, v)
	}
	return versions, total, rows.Err()
}

// Get は ID でバージョンを取得する。
func (r *OssVersionRepository) Get(ctx context.Context, id string) (*model.OssVersion, error) {
	row := r.DB.QueryRowContext(ctx, `SELECT id, oss_id, version, release_date, license_expression_raw, license_concluded, purl, cpe_list, hash_sha256, review_status, last_reviewed_at, scope_status, created_at, updated_at FROM oss_versions WHERE id = ?`, id)
	var v model.OssVersion
	var releaseDate sql.NullTime
	var licenseRaw, licenseConc, purl, hash sql.NullString
	var lastReviewed sql.NullTime
	var cpeList pq.StringArray
	if err := row.Scan(&v.ID, &v.OssID, &v.Version, &releaseDate, &licenseRaw, &licenseConc, &purl, &cpeList, &hash, &v.ReviewStatus, &lastReviewed, &v.ScopeStatus, &v.CreatedAt, &v.UpdatedAt); err != nil {
		return nil, err
	}
	v.ReleaseDate = timePtr(releaseDate)
	v.LicenseExpressionRaw = strPtr(licenseRaw)
	v.LicenseConcluded = strPtr(licenseConc)
	v.Purl = strPtr(purl)
	v.CpeList = []string(cpeList)
	v.HashSha256 = strPtr(hash)
	v.LastReviewedAt = timePtr(lastReviewed)
	return &v, nil
}

// Create は新しいバージョンを登録する。
func (r *OssVersionRepository) Create(ctx context.Context, v *model.OssVersion) error {
	_, err := r.DB.ExecContext(ctx,
		`INSERT INTO oss_versions (id, oss_id, version, release_date, license_expression_raw, license_concluded, purl, cpe_list, hash_sha256, review_status, last_reviewed_at, scope_status, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		v.ID, v.OssID, v.Version, v.ReleaseDate, v.LicenseExpressionRaw, v.LicenseConcluded, v.Purl, pq.Array(v.CpeList), v.HashSha256, v.ReviewStatus, v.LastReviewedAt, v.ScopeStatus, v.CreatedAt, v.UpdatedAt,
	)
	return err
}

// Update は既存バージョンを更新する。
func (r *OssVersionRepository) Update(ctx context.Context, v *model.OssVersion) error {
	_, err := r.DB.ExecContext(ctx,
		`UPDATE oss_versions SET release_date = ?, license_expression_raw = ?, license_concluded = ?, purl = ?, cpe_list = ?, hash_sha256 = ?, review_status = ?, last_reviewed_at = ?, scope_status = ?, updated_at = ? WHERE id = ?`,
		v.ReleaseDate, v.LicenseExpressionRaw, v.LicenseConcluded, v.Purl, pq.Array(v.CpeList), v.HashSha256, v.ReviewStatus, v.LastReviewedAt, v.ScopeStatus, v.UpdatedAt, v.ID,
	)
	return err
}

// Delete は ID 指定でバージョンを削除する。
func (r *OssVersionRepository) Delete(ctx context.Context, id string) error {
	_, err := r.DB.ExecContext(ctx, `DELETE FROM oss_versions WHERE id = ?`, id)
	return err
}
