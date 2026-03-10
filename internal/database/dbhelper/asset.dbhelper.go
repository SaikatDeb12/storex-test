package dbhelper

import (
	"errors"
	"time"

	"github.com/SaikatDeb12/storeX/internal/database"
	"github.com/SaikatDeb12/storeX/internal/models"
	"github.com/jmoiron/sqlx"
)

func InsertKeyboardDetails(assetID string, req models.KeyboardRequest) error {
	SQL := `
		INSERT INTO keyboards(asset_id, connectivity, layout)
		VALUES($1, $2, $3)
	`

	args := []interface{}{
		assetID,
		req.Connectivity,
		req.Layout,
	}

	_, err := database.DB.Exec(SQL, args...)
	return err
}

func InsertLaptopDetails(assetID string, req models.LaptopRequest) error {
	SQL := `
		INSERT INTO laptops(asset_id, processor, ram, storage, operating_system, charger, device_password)
		VALUES($1, $2, $3, $4, $5, $6, $7)
	`

	args := []interface{}{
		assetID,
		req.Processor,
		req.RAM,
		req.Storage,
		req.OperatingSystem,
		req.Charger,
		req.DevicePassword,
	}

	_, err := database.DB.Exec(SQL, args...)
	return err
}

func InsertMouseDetails(assetID string, req models.MouseRequest) error {
	SQL := `
		INSERT INTO mice(asset_id, dpi, connectivity)
		VALUES($1, $2, $3)
	`

	args := []interface{}{
		assetID,
		req.DPI,
		req.Connectivity,
	}

	_, err := database.DB.Exec(SQL, args...)
	return err
}

func InsertMobileDetails(assetID string, req models.MobileRequest) error {
	SQL := `
		INSERT INTO mobiles(asset_id, operating_system, ram, storage, charger, device_password)
		VALUES($1, $2, $3, $4, $5, $6)
	`

	args := []interface{}{
		assetID,
		req.OperatingSystem,
		req.RAM,
		req.Storage,
		req.Charger,
		req.DevicePassword,
	}

	_, err := database.DB.Exec(SQL, args...)
	return err
}

func CreateAsset(model models.CreateAssetRequest) (string, error) {
	SQL := `
		INSERT INTO assets(brand, model, serial_number, asset_type, owner_type, warranty_start, warranty_end)
		VALUES($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	args := []interface{}{
		model.Brand,
		model.Model,
		model.SerialNumber,
		model.Type,
		model.Owner,
		model.WarrantyStart,
		model.WarrantyEnd,
	}

	var assetID string
	err := database.DB.Get(&assetID, SQL, args...)
	if err != nil {
		return assetID, err
	}
	return assetID, nil
}

func FetchAssets(brand, model, assetType, serial_number, status, owner string, limit, offset int) ([]models.AssetAssignRequest, error) {
	SQL := `SELECT id, brand, model, asset_type, serial_number, status, owner_type, assigned_by_id, assigned_to_id, assigned_at, warranty_start, warranty_end, service_start, service_end, returned_at, created_at, updated_at 
          FROM assets
          WHERE archived_at IS NULL 
          AND (
              $1= '' or brand LIKE '%'||$1||'%'
          )
          AND(
              $2 ='' or model LIKE '%'||$2||'%'
          )
          AND (
              $3='' or asset_type::text LIKE '%'||$3||'%'
          )
          AND(
              $4 ='' or serial_number LIKE '%'||$4||'%'
          )
          AND(
              $5='' or status::text LIKE '%'||$5||'%'
          )
          AND(
              $6=''or owner_type::text LIKE '%'||$6||'%'
          )
          ORDER BY created_at
		  LIMIT $7 OFFSET $8
          `
	var result []models.AssetAssignRequest
	err := database.DB.Select(&result, SQL, brand, model, assetType, serial_number, status, owner, limit, offset)
	return result, err
}

func GettingAssetsCount() (models.DashboardSummaryRequest, error) {
	SQL := `
		SELECT 
		COUNT(*) AS total,
		COUNT(*) FILTER (WHERE status='available' ) AS available,
		COUNT(*) FILTER (WHERE status='assigned' ) AS assigned,
		COUNT(*) FILTER (WHERE status='under_repair' ) AS waitingForRepair,
		COUNT(*) FILTER (WHERE status='in_service' ) AS inService,
		COUNT(*) FILTER (WHERE status='damaged' ) AS damaged
		FROM assets
		WHERE archived_at IS NULL
	`
	var result models.DashboardSummaryRequest
	err := database.DB.Get(&result, SQL)
	return result, err
}

func AssignedAssets(id, assignedById, assignedTo string) error {
	SQL := `UPDATE assets
          SET assigned_to_id=$3,
			  assigned_by_id=$2,
              assigned_at=NOW(),
              status='assigned',
              updated_at=NOW()
          WHERE id=$1
          AND archived_at IS NULL 
              `
	// _, err := database.DB.Exec(SQL, assignedById, assignedTo, id)
	res, err := database.DB.Exec(SQL, id, assignedById, assignedTo)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("asset not found or already archived")
	}
	return nil
}

func UpdateAsset(tx *sqlx.Tx, assetID, brand, model, serialNo, assetType, owner string, warrantyStart, warrantyEnd time.Time) error {
	query := `UPDATE assets
            set brand = $2, model = $3, serial_number = $4, asset_type=$5, owner_type=$6, warranty_start = $7,warranty_end=$8, updated_at =now()
            where id= $1 and archived_at is null `
	_, err := tx.Exec(query, assetID, brand, model, serialNo, assetType, owner, warrantyStart, warrantyEnd)
	if err != nil {
		return err
	}
	return nil
}

func UpdateLaptop(tx *sqlx.Tx, assetID string, laptop *models.LaptopRequest) error {
	query := `
    UPDATE laptops
    SET
        processor = $2,
        ram = $3,
        storage = $4,
        operating_system = $5,
        charger = $6,
        device_password = $7
    WHERE asset_id = $1
    `

	_, err := tx.Exec(query,
		assetID,
		laptop.Processor,
		laptop.RAM,
		laptop.Storage,
		laptop.OperatingSystem,
		laptop.Charger,
		laptop.DevicePassword,
	)

	return err
}

func UpdateMouse(tx *sqlx.Tx, assetID string, mouse *models.MouseRequest) error {
	query := `
    UPDATE mice
    SET
        dpi = $2,
        connectivity = $3
    WHERE asset_id = $1
    `

	_, err := tx.Exec(query, assetID, mouse.DPI, mouse.Connectivity)
	return err
}

func UpdateKeyboard(tx *sqlx.Tx, assetID string, keyboard *models.KeyboardRequest) error {
	query := `
    UPDATE keyboards
    SET
        layout = $2,
        connectivity = $3
    WHERE asset_id = $1
    `

	_, err := tx.Exec(query, assetID, keyboard.Layout, keyboard.Connectivity)
	return err
}

func UpdateMobile(tx *sqlx.Tx, assetID string, mobile *models.MobileRequest) error {
	query := `
    UPDATE mobiles
    SET
        operating_system = $2,
        ram = $3,
        storage = $4,
        charger = $5,
        device_password = $6
    WHERE asset_id = $1
    `

	_, err := tx.Exec(
		query,
		assetID,
		mobile.OperatingSystem,
		mobile.RAM,
		mobile.Storage,
		mobile.Charger,
		mobile.DevicePassword,
	)

	return err
}

func SentToService(assetId string, serviceStart, serviceEnd time.Time) error {
	query := `update assets set status='in_service',service_start=$2,service_end=$3,updated_at=now()
              where id=$1 and archived_at is NULL and status ='available'`
	_, err := database.DB.Exec(query, assetId, serviceStart, serviceEnd)
	if err != nil {
		return err
	}
	return nil
}
