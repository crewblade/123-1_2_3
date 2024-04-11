package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/crewblade/banner-management-service/internal/domain/models"
	"github.com/crewblade/banner-management-service/internal/lib/api/errs"
	"github.com/crewblade/banner-management-service/internal/lib/utils"
	"github.com/lib/pq"
)

func (s *Storage) GetBanners(ctx context.Context, featureID, tagID, limit, offset *int) ([]models.Banner, error) {
	const op = "repo.postgres.GetBanners"

	stmt, err := s.db.PrepareContext(ctx, `
    SELECT id, content, feature_id, tag_ids, is_active, created_at, updated_at
	FROM banners 
	WHERE (feature_id = $1 OR $1 IS NULL) 
		AND ($2 = ANY(tag_ids) OR $2 IS NULL) 
	LIMIT $3 OFFSET $4;
    `)
	if err != nil {
		return nil, fmt.Errorf("%s: prepare context %w", op, err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, featureID, tagID, *limit, *offset)
	if err != nil {
		return nil, fmt.Errorf("%s: query context %w", op, err)
	}
	defer rows.Close()

	var banners []models.Banner
	for rows.Next() {
		var banner models.Banner
		var tagIDsString string
		err := rows.Scan(
			&banner.BannerID,
			&banner.Content,
			&banner.FeatureID,
			&tagIDsString,
			&banner.IsActive,
			&banner.CreatedAt,
			&banner.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: scan row %w", op, err)
		}
		tagIDs, err := utils.StringToIntArray(tagIDsString)
		if err != nil {
			return nil, fmt.Errorf("%s: error parsing tag IDs: %w", op, err)
		}
		banner.TagIDs = tagIDs
		banners = append(banners, banner)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows error %w", op, err)
	}

	return banners, nil

}

func (s *Storage) SaveBanner(
	ctx context.Context,
	tagIDs []int,
	featureID int,
	content json.RawMessage,
	isActive bool,
) (int, error) {
	const op = "repo.postgres.SaveBanner"
	const errValue = -1
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return errValue, fmt.Errorf("%s: begin transaction: %w", op, err)
	}
	defer tx.Rollback()

	exists, err := isBannerExistsInTx(ctx, tx, featureID, tagIDs, -1)
	if err != nil {
		return errValue, fmt.Errorf("failed to check banner existence: %w", err)
	}
	if exists {
		return errValue, fmt.Errorf("banner with the same feature_id and tag_id already exists")
	}

	stmt, err := s.db.PrepareContext(ctx, "INSERT INTO banners(tag_ids, feature_id, content, is_active) VALUES ($1, $2, $3, $4) RETURNING id")
	if err != nil {
		return errValue, fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	defer stmt.Close()

	var bannerID int
	err = stmt.QueryRowContext(ctx, pq.Array(tagIDs), featureID, content, isActive).Scan(&bannerID)
	if err != nil {
		return errValue, fmt.Errorf("%s: execute statement: %w", op, err)
	}

	if err := tx.Commit(); err != nil {
		return errValue, fmt.Errorf("%s: commit transaction: %w", op, err)
	}

	return bannerID, nil
}

func (s *Storage) DeleteBanner(ctx context.Context, bannerID int) error {
	const op = "repo.postgres.DeleteBanner"

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%s: begin transaction: %w", op, err)
	}
	defer tx.Rollback()

	existsID, err := isBannerIDExistsInTx(ctx, tx, bannerID)
	if err != nil {
		return fmt.Errorf("failed to check bannerID existence: %w", err)
	}
	if !existsID {
		return errs.ErrBannerNotFound
	}

	stmt, err := tx.PrepareContext(ctx, "DELETE FROM banners WHERE id = $1")
	if err != nil {
		return fmt.Errorf("%s: prepare context: %w", op, err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, bannerID)
	if err != nil {
		return fmt.Errorf("%s: execute statement: %w", op, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%s: commit transaction: %w", op, err)
	}

	return nil
}

func (s *Storage) UpdateBanner(
	ctx context.Context,
	bannerID int,
	tagIDs []int,
	featureID int,
	content json.RawMessage,
	isActive bool,
) error {
	const op = "repo.postgres.UpdateBanner"

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%s: start transaction: %w", op, err)
	}
	defer tx.Rollback()

	existsID, err := isBannerIDExistsInTx(ctx, tx, bannerID)
	if err != nil {
		return fmt.Errorf("failed to check bannerID existence: %w", err)
	}
	if !existsID {
		return errs.ErrBannerNotFound
	}

	exists, err := isBannerExistsInTx(ctx, tx, featureID, tagIDs, bannerID)
	if err != nil {
		return fmt.Errorf("failed to check banner existence: %w", err)
	}
	if exists {
		return fmt.Errorf("banner with the same feature_id and tag_id already exists")
	}

	stmt, err := tx.PrepareContext(ctx, "UPDATE banners SET tag_ids = $1, feature_id = $2, content = $3, is_active = $4, updated_at = NOW() WHERE id = $5")
	if err != nil {
		return fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, pq.Array(tagIDs), featureID, content, isActive, bannerID)
	if err != nil {
		return fmt.Errorf("%s: execute statement: %w", op, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%s: commit transaction: %w", op, err)
	}

	return nil
}

func (s *Storage) GetUserBanner(
	ctx context.Context,
	tagID int,
	featureID int,
) (json.RawMessage, bool, error) {

	const op = "repo.postgres.GetUserBanner"

	stmt, err := s.db.PrepareContext(ctx, "SELECT content, is_active FROM banners WHERE feature_id = $1 AND $2 = ANY(tag_ids)")
	if err != nil {
		return nil, false, fmt.Errorf("%s: prepare context %w", op, err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, featureID, tagID)
	var content json.RawMessage
	var isActive bool
	err = row.Scan(&content, &isActive)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, false, errs.ErrBannerNotFound
		}

		return nil, false, fmt.Errorf("%s: row scan %w", op, err)
	}

	return content, isActive, nil

}

func isBannerIDExistsInTx(ctx context.Context, tx *sql.Tx, bannerID int) (bool, error) {
	const op = "repo.postgres.isBannerIDExists"

	var exists bool
	err := tx.QueryRowContext(ctx, "SELECT EXISTS (SELECT 1 FROM banners WHERE id = $1)", bannerID).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("%s: query row scan: %w", op, err)
	}

	return exists, nil
}

func isBannerExistsInTx(ctx context.Context, tx *sql.Tx, featureID int, tagIDs []int, bannerID int) (bool, error) {
	const op = "repo.postgres.isBannerExistsInTx"

	query := `SELECT COUNT(*)
FROM banners
WHERE feature_id = $1 AND EXISTS (
    SELECT 1
    FROM UNNEST($2::INT[]) AS tag
    WHERE tag = ANY(tag_ids)
	)
	AND id != $3;
;
`
	var count int
	err := tx.QueryRowContext(ctx, query, featureID, pq.Array(tagIDs), bannerID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("%s: failed to execute query: %w", op, err)
	}

	return count > 0, nil
}
