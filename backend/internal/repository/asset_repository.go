package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"it-asset-management/backend/internal/models"
)

var (
	ErrNotFound       = errors.New("asset not found")
	ErrInvalidStatus  = errors.New("invalid status transition")
	ErrDuplicateCode  = errors.New("asset_code or serial_number already exists")
)

type AssetRepository interface {
	GetAssets(ctx context.Context, search, status, assetType string, page, limit int) ([]models.Asset, int, error)
	GetStats(ctx context.Context) (*models.DashboardStats, error)
	GetByID(ctx context.Context, id int) (*models.Asset, error)
	Create(ctx context.Context, req *models.CreateAssetRequest) (*models.Asset, error)
	Update(ctx context.Context, id int, req *models.UpdateAssetRequest) (*models.Asset, error)
	Delete(ctx context.Context, id int) error
	Assign(ctx context.Context, id int, req *models.AssignAssetRequest) (*models.Asset, error)
	ReturnAsset(ctx context.Context, id int, req *models.ReturnAssetRequest) (*models.Asset, error)
	Repair(ctx context.Context, id int, req *models.StatusChangeRequest) (*models.Asset, error)
	Retire(ctx context.Context, id int, req *models.StatusChangeRequest) (*models.Asset, error)
	GetHistory(ctx context.Context, assetID int) ([]models.AssetHistory, error)
}

type PostgresAssetRepository struct {
	db *sql.DB
}

func NewPostgresAssetRepository(db *sql.DB) *PostgresAssetRepository {
	return &PostgresAssetRepository{db: db}
}

func (r *PostgresAssetRepository) GetAssets(ctx context.Context, search, status, assetType string, page, limit int) ([]models.Asset, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit

	var conditions []string
	var args []interface{}
	argIdx := 1

	if strings.TrimSpace(search) != "" {
		searchTerm := "%" + strings.TrimSpace(search) + "%"
		conditions = append(conditions, fmt.Sprintf("(asset_code ILIKE $%d OR asset_name ILIKE $%d OR serial_number ILIKE $%d)", argIdx, argIdx, argIdx))
		args = append(args, searchTerm)
		argIdx++
	}

	if strings.TrimSpace(status) != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIdx))
		args = append(args, strings.TrimSpace(status))
		argIdx++
	}

	if strings.TrimSpace(assetType) != "" {
		conditions = append(conditions, fmt.Sprintf("asset_type = $%d", argIdx))
		args = append(args, strings.TrimSpace(assetType))
		argIdx++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM assets %s", whereClause)
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count assets: %w", err)
	}

	selectQuery := fmt.Sprintf(`
		SELECT id, asset_code, asset_name, asset_type, serial_number, 
		       to_char(purchase_date, 'YYYY-MM-DD') AS purchase_date, 
		       status, assigned_to, location, created_at, updated_at
		FROM assets
		%s
		ORDER BY id DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIdx, argIdx+1)

	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, selectQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch assets: %w", err)
	}
	defer rows.Close()

	var assets []models.Asset
	for rows.Next() {
		var a models.Asset
		err := rows.Scan(
			&a.ID, &a.AssetCode, &a.AssetName, &a.AssetType, &a.SerialNumber,
			&a.PurchaseDate, &a.Status, &a.AssignedTo, &a.Location,
			&a.CreatedAt, &a.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan asset row: %w", err)
		}
		assets = append(assets, a)
	}

	if assets == nil {
		assets = []models.Asset{}
	}

	return assets, total, nil
}

func (r *PostgresAssetRepository) GetStats(ctx context.Context) (*models.DashboardStats, error) {
	query := `
		SELECT 
			COUNT(*) AS total,
			COUNT(*) FILTER (WHERE status = 'AVAILABLE') AS available,
			COUNT(*) FILTER (WHERE status = 'ASSIGNED') AS assigned,
			COUNT(*) FILTER (WHERE status = 'UNDER_REPAIR') AS under_repair,
			COUNT(*) FILTER (WHERE status = 'RETIRED') AS retired
		FROM assets
	`
	var stats models.DashboardStats
	err := r.db.QueryRowContext(ctx, query).Scan(
		&stats.Total,
		&stats.Available,
		&stats.Assigned,
		&stats.UnderRepair,
		&stats.Retired,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to compute dashboard stats: %w", err)
	}
	return &stats, nil
}

func (r *PostgresAssetRepository) GetByID(ctx context.Context, id int) (*models.Asset, error) {
	query := `
		SELECT id, asset_code, asset_name, asset_type, serial_number,
		       to_char(purchase_date, 'YYYY-MM-DD') AS purchase_date,
		       status, assigned_to, location, created_at, updated_at
		FROM assets
		WHERE id = $1
	`
	var a models.Asset
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&a.ID, &a.AssetCode, &a.AssetName, &a.AssetType, &a.SerialNumber,
		&a.PurchaseDate, &a.Status, &a.AssignedTo, &a.Location,
		&a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get asset by ID: %w", err)
	}
	return &a, nil
}

func (r *PostgresAssetRepository) Create(ctx context.Context, req *models.CreateAssetRequest) (*models.Asset, error) {
	status := req.Status
	if status == "" {
		status = models.StatusAvailable
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO assets (asset_code, asset_name, asset_type, serial_number, purchase_date, status, assigned_to, location)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`
	var asset models.Asset
	asset.AssetCode = req.AssetCode
	asset.AssetName = req.AssetName
	asset.AssetType = req.AssetType
	asset.SerialNumber = req.SerialNumber
	asset.PurchaseDate = req.PurchaseDate
	asset.Status = status
	asset.AssignedTo = req.AssignedTo
	asset.Location = req.Location

	err = tx.QueryRowContext(
		ctx, query,
		asset.AssetCode, asset.AssetName, asset.AssetType, asset.SerialNumber,
		asset.PurchaseDate, asset.Status, asset.AssignedTo, asset.Location,
	).Scan(&asset.ID, &asset.CreatedAt, &asset.UpdatedAt)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") {
			return nil, ErrDuplicateCode
		}
		return nil, fmt.Errorf("failed to insert asset: %w", err)
	}

	historyQuery := `
		INSERT INTO asset_history (asset_id, action, performed_by, details)
		VALUES ($1, 'CREATED', 'Admin', $2)
	`
	initialNote := fmt.Sprintf("Asset registered with status %s at %s", asset.Status, asset.Location)
	if asset.AssignedTo != nil && *asset.AssignedTo != "" {
		initialNote = fmt.Sprintf("Asset registered and assigned to %s at %s", *asset.AssignedTo, asset.Location)
	}
	_, err = tx.ExecContext(ctx, historyQuery, asset.ID, initialNote)
	if err != nil {
		return nil, fmt.Errorf("failed to log asset creation history: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &asset, nil
}

func (r *PostgresAssetRepository) Update(ctx context.Context, id int, req *models.UpdateAssetRequest) (*models.Asset, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		UPDATE assets
		SET asset_name = $1, asset_type = $2, purchase_date = $3, location = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $5
		RETURNING id, asset_code, asset_name, asset_type, serial_number, 
		          to_char(purchase_date, 'YYYY-MM-DD') AS purchase_date, 
		          status, assigned_to, location, created_at, updated_at
	`
	var a models.Asset
	err = tx.QueryRowContext(
		ctx, query,
		req.AssetName, req.AssetType, req.PurchaseDate, req.Location, id,
	).Scan(
		&a.ID, &a.AssetCode, &a.AssetName, &a.AssetType, &a.SerialNumber,
		&a.PurchaseDate, &a.Status, &a.AssignedTo, &a.Location,
		&a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to update asset: %w", err)
	}

	historyQuery := `
		INSERT INTO asset_history (asset_id, action, performed_by, details)
		VALUES ($1, 'UPDATED', 'Admin', $2)
	`
	details := fmt.Sprintf("Updated asset details (Name: %s, Type: %s, Location: %s)", a.AssetName, a.AssetType, a.Location)
	_, err = tx.ExecContext(ctx, historyQuery, a.ID, details)
	if err != nil {
		return nil, fmt.Errorf("failed to log update history: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &a, nil
}

func (r *PostgresAssetRepository) Delete(ctx context.Context, id int) error {
	query := "DELETE FROM assets WHERE id = $1"
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete asset: %w", err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *PostgresAssetRepository) Assign(ctx context.Context, id int, req *models.AssignAssetRequest) (*models.Asset, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var current models.Asset
	checkQuery := `
		SELECT id, asset_code, status, assigned_to, location 
		FROM assets 
		WHERE id = $1 
		FOR UPDATE
	`
	err = tx.QueryRowContext(ctx, checkQuery, id).Scan(
		&current.ID, &current.AssetCode, &current.Status, &current.AssignedTo, &current.Location,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	// Business Rule: An asset CANNOT be assigned if it is already ASSIGNED, UNDER_REPAIR, or RETIRED
	if current.Status != models.StatusAvailable {
		return nil, fmt.Errorf("%w: asset '%s' cannot be assigned because its current status is %s", ErrInvalidStatus, current.AssetCode, current.Status)
	}

	newLocation := current.Location
	if strings.TrimSpace(req.Location) != "" {
		newLocation = strings.TrimSpace(req.Location)
	}

	updateQuery := `
		UPDATE assets
		SET status = 'ASSIGNED', assigned_to = $1, location = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
		RETURNING id, asset_code, asset_name, asset_type, serial_number,
		          to_char(purchase_date, 'YYYY-MM-DD') AS purchase_date,
		          status, assigned_to, location, created_at, updated_at
	`
	var updated models.Asset
	err = tx.QueryRowContext(ctx, updateQuery, req.AssignedTo, newLocation, id).Scan(
		&updated.ID, &updated.AssetCode, &updated.AssetName, &updated.AssetType, &updated.SerialNumber,
		&updated.PurchaseDate, &updated.Status, &updated.AssignedTo, &updated.Location,
		&updated.CreatedAt, &updated.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update asset assignment: %w", err)
	}

	performedBy := req.PerformedBy
	if performedBy == "" {
		performedBy = "Admin"
	}
	historyDetail := fmt.Sprintf("Assigned to %s (Location: %s)", req.AssignedTo, newLocation)
	if req.Notes != "" {
		historyDetail += fmt.Sprintf(" - Notes: %s", req.Notes)
	}

	historyQuery := `
		INSERT INTO asset_history (asset_id, action, performed_by, details)
		VALUES ($1, 'ASSIGNED', $2, $3)
	`
	_, err = tx.ExecContext(ctx, historyQuery, id, performedBy, historyDetail)
	if err != nil {
		return nil, fmt.Errorf("failed to log assignment history: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit assignment transaction: %w", err)
	}

	return &updated, nil
}

func (r *PostgresAssetRepository) ReturnAsset(ctx context.Context, id int, req *models.ReturnAssetRequest) (*models.Asset, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var current models.Asset
	checkQuery := `
		SELECT id, asset_code, status, assigned_to, location 
		FROM assets 
		WHERE id = $1 
		FOR UPDATE
	`
	err = tx.QueryRowContext(ctx, checkQuery, id).Scan(
		&current.ID, &current.AssetCode, &current.Status, &current.AssignedTo, &current.Location,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if current.Status != models.StatusAssigned {
		return nil, fmt.Errorf("%w: cannot return asset '%s' because its status is %s (must be ASSIGNED)", ErrInvalidStatus, current.AssetCode, current.Status)
	}

	previousAssignee := "Unknown"
	if current.AssignedTo != nil && *current.AssignedTo != "" {
		previousAssignee = *current.AssignedTo
	}

	newLocation := current.Location
	if strings.TrimSpace(req.Location) != "" {
		newLocation = strings.TrimSpace(req.Location)
	}

	updateQuery := `
		UPDATE assets
		SET status = 'AVAILABLE', assigned_to = NULL, location = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
		RETURNING id, asset_code, asset_name, asset_type, serial_number,
		          to_char(purchase_date, 'YYYY-MM-DD') AS purchase_date,
		          status, assigned_to, location, created_at, updated_at
	`
	var updated models.Asset
	err = tx.QueryRowContext(ctx, updateQuery, newLocation, id).Scan(
		&updated.ID, &updated.AssetCode, &updated.AssetName, &updated.AssetType, &updated.SerialNumber,
		&updated.PurchaseDate, &updated.Status, &updated.AssignedTo, &updated.Location,
		&updated.CreatedAt, &updated.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to return asset: %w", err)
	}

	performedBy := req.PerformedBy
	if performedBy == "" {
		performedBy = "Admin"
	}
	historyDetail := fmt.Sprintf("Returned by %s into storage/available (Location: %s)", previousAssignee, newLocation)
	if req.Notes != "" {
		historyDetail += fmt.Sprintf(" - Notes: %s", req.Notes)
	}

	historyQuery := `
		INSERT INTO asset_history (asset_id, action, performed_by, details)
		VALUES ($1, 'RETURNED', $2, $3)
	`
	_, err = tx.ExecContext(ctx, historyQuery, id, performedBy, historyDetail)
	if err != nil {
		return nil, fmt.Errorf("failed to log return history: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit return transaction: %w", err)
	}

	return &updated, nil
}

func (r *PostgresAssetRepository) Repair(ctx context.Context, id int, req *models.StatusChangeRequest) (*models.Asset, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var current models.Asset
	checkQuery := "SELECT id, asset_code, status, assigned_to, location FROM assets WHERE id = $1 FOR UPDATE"
	err = tx.QueryRowContext(ctx, checkQuery, id).Scan(&current.ID, &current.AssetCode, &current.Status, &current.AssignedTo, &current.Location)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if current.Status == models.StatusRetired {
		return nil, fmt.Errorf("%w: cannot send retired asset to repair", ErrInvalidStatus)
	}

	newLocation := current.Location
	if strings.TrimSpace(req.Location) != "" {
		newLocation = strings.TrimSpace(req.Location)
	}

	updateQuery := `
		UPDATE assets
		SET status = 'UNDER_REPAIR', assigned_to = NULL, location = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
		RETURNING id, asset_code, asset_name, asset_type, serial_number,
		          to_char(purchase_date, 'YYYY-MM-DD') AS purchase_date,
		          status, assigned_to, location, created_at, updated_at
	`
	var updated models.Asset
	err = tx.QueryRowContext(ctx, updateQuery, newLocation, id).Scan(
		&updated.ID, &updated.AssetCode, &updated.AssetName, &updated.AssetType, &updated.SerialNumber,
		&updated.PurchaseDate, &updated.Status, &updated.AssignedTo, &updated.Location,
		&updated.CreatedAt, &updated.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to mark asset under repair: %w", err)
	}

	performedBy := req.PerformedBy
	if performedBy == "" {
		performedBy = "Admin"
	}
	historyDetail := fmt.Sprintf("Marked status as UNDER_REPAIR (Location: %s)", newLocation)
	if req.Notes != "" {
		historyDetail += fmt.Sprintf(" - Reason: %s", req.Notes)
	}

	historyQuery := `
		INSERT INTO asset_history (asset_id, action, performed_by, details)
		VALUES ($1, 'REPAIR', $2, $3)
	`
	_, err = tx.ExecContext(ctx, historyQuery, id, performedBy, historyDetail)
	if err != nil {
		return nil, fmt.Errorf("failed to log repair history: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit repair transaction: %w", err)
	}

	return &updated, nil
}

func (r *PostgresAssetRepository) Retire(ctx context.Context, id int, req *models.StatusChangeRequest) (*models.Asset, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var current models.Asset
	checkQuery := "SELECT id, asset_code, status, assigned_to, location FROM assets WHERE id = $1 FOR UPDATE"
	err = tx.QueryRowContext(ctx, checkQuery, id).Scan(&current.ID, &current.AssetCode, &current.Status, &current.AssignedTo, &current.Location)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if current.Status == models.StatusRetired {
		return nil, fmt.Errorf("%w: asset is already retired", ErrInvalidStatus)
	}

	newLocation := "Disposal / Decommissioned Storage"
	if strings.TrimSpace(req.Location) != "" {
		newLocation = strings.TrimSpace(req.Location)
	}

	updateQuery := `
		UPDATE assets
		SET status = 'RETIRED', assigned_to = NULL, location = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
		RETURNING id, asset_code, asset_name, asset_type, serial_number,
		          to_char(purchase_date, 'YYYY-MM-DD') AS purchase_date,
		          status, assigned_to, location, created_at, updated_at
	`
	var updated models.Asset
	err = tx.QueryRowContext(ctx, updateQuery, newLocation, id).Scan(
		&updated.ID, &updated.AssetCode, &updated.AssetName, &updated.AssetType, &updated.SerialNumber,
		&updated.PurchaseDate, &updated.Status, &updated.AssignedTo, &updated.Location,
		&updated.CreatedAt, &updated.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to retire asset: %w", err)
	}

	performedBy := req.PerformedBy
	if performedBy == "" {
		performedBy = "Admin"
	}
	historyDetail := fmt.Sprintf("Decommissioned and marked RETIRED (Location: %s)", newLocation)
	if req.Notes != "" {
		historyDetail += fmt.Sprintf(" - Reason: %s", req.Notes)
	}

	historyQuery := `
		INSERT INTO asset_history (asset_id, action, performed_by, details)
		VALUES ($1, 'RETIRED', $2, $3)
	`
	_, err = tx.ExecContext(ctx, historyQuery, id, performedBy, historyDetail)
	if err != nil {
		return nil, fmt.Errorf("failed to log retirement history: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit retirement transaction: %w", err)
	}

	return &updated, nil
}

func (r *PostgresAssetRepository) GetHistory(ctx context.Context, assetID int) ([]models.AssetHistory, error) {
	// First check if asset exists
	var exists bool
	err := r.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM assets WHERE id = $1)", assetID).Scan(&exists)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrNotFound
	}

	query := `
		SELECT id, asset_id, action, performed_by, details, created_at
		FROM asset_history
		WHERE asset_id = $1
		ORDER BY created_at DESC, id DESC
	`
	rows, err := r.db.QueryContext(ctx, query, assetID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch asset history: %w", err)
	}
	defer rows.Close()

	var list []models.AssetHistory
	for rows.Next() {
		var h models.AssetHistory
		if err := rows.Scan(&h.ID, &h.AssetID, &h.Action, &h.PerformedBy, &h.Details, &h.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan history row: %w", err)
		}
		list = append(list, h)
	}

	if list == nil {
		list = []models.AssetHistory{}
	}

	return list, nil
}
