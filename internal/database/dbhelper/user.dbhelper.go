package dbhelper

import (
	"errors"

	"github.com/SaikatDeb12/storeX/internal/database"
	"github.com/SaikatDeb12/storeX/internal/models"
	"github.com/jmoiron/sqlx"
)

func CheckUserExistsByEmail(email string) (bool, error) {
	SQL := `
		SELECT COUNT(*)
		FROM users
		WHERE email=TRIM(LOWER($1)) AND archived_at IS NULL
	`
	var count int
	err := database.DB.Get(&count, SQL, email)
	if err != nil {
		return false, err
	}

	return count > 0, err
}

func CreateUser(tx *sqlx.Tx, name, email, phoneNumber, role, employment, hashedPassword string) (string, error) {
	SQL := `
		INSERT INTO users(name, email, phone_number, role, employment, password)
		VALUES($1, TRIM(LOWER($2)), $3, $4, $5, $6)
		RETURNING id
	`
	var userID string
	err := tx.Get(&userID, SQL, name, email, phoneNumber, role, employment, hashedPassword)
	if err != nil {
		return "", err
	}
	return string(userID), nil
}

func CreateSessionOnRegister(tx *sqlx.Tx, userID string) (string, error) {
	SQL := `
		INSERT INTO user_sessions(user_id)
		VALUES($1)
		RETURNING id
	`
	var sessionID string
	err := tx.Get(&sessionID, SQL, userID)
	if err != nil {
		return "", err
	}
	return string(sessionID), nil
}

func CreateSessionOnLogin(userID string) (string, error) {
	SQL := `
		INSERT INTO user_sessions(user_id)
		VALUES($1)
		RETURNING id
	`
	var sessionID string
	err := database.DB.Get(&sessionID, SQL, userID)
	if err != nil {
		return "", err
	}
	return string(sessionID), nil
}

func GetUserAuthByEmail(email string) (models.User, error) {
	SQL := `
		SELECT id, email, password  
		FROM users 
		WHERE email=TRIM(LOWER($1)) AND archived_at IS NULL
	`
	var user models.User
	err := database.DB.Get(&user, SQL, email)
	if err != nil {
		return user, err
	}
	return user, nil
}

func FetchAssetInfo(userID string) ([]models.AssetInfoRequest, error) {
	SQL := `
		SELECT id, brand, model, status, asset_type
		FROM assets
		WHERE archived_at IS NULL AND assigned_to_id=$1
	`
	assetDetails := make([]models.AssetInfoRequest, 0)
	err := database.DB.Select(&assetDetails, SQL, userID)
	return assetDetails, err
}

func FetchUsers(name, role, employment, assetStatus string) ([]models.UserInfoRequest, error) {
	SQL := `
		SELECT id, name, email, phone_number, role, employment, created_at
		FROM users
		WHERE ($1 = '' OR name LIKE '%' || $1 || '%')
		AND ($2 = '' OR role::TEXT=$2)
		AND ($3 = '' OR employment::TEXT=$3)
		AND archived_at IS NOT NULL
	`
	users := make([]models.UserInfoRequest, 0)
	err := database.DB.Select(&users, SQL, name, role, employment)
	if err != nil {
		return users, err
	}

	// to make change on the original slice:
	filteredUsers := make([]models.UserInfoRequest, 0)
	for _, user := range users {
		assetDetails, err := FetchAssetsInfo(user.ID, assetStatus)
		if err != nil {
			return users, err
		}
		if assetStatus == "available" {
			if len(assetDetails) > 0 {
				continue
			} else {
				if len(assetDetails) == 0 {
					continue
				}
			}
		}

		user.AssetDetails = assetDetails
		filteredUsers = append(filteredUsers, user)
	}
	return filteredUsers, nil

	// change in the copy not the original
	// for _, user := range users {
	// 	userDetails, err := GetAssetInfo(user.ID)
	// 	if err != nil {
	// 		return users, err
	// 	}
	// 	fmt.Println(userDetails)
	// 	user.AssetDetails = userDetails
	// }
}

func FetchUserByID(userID string) (models.UserInfoRequest, error) {
	SQL := `
		SELECT id, name, email, phone_number, role, employment, created_at
		FROM users
		WHERE archived_at IS NULL AND id=$1
	`
	var user models.UserInfoRequest
	err := database.DB.Get(&user, SQL, userID)
	if err != nil {
		return user, err
	}

	assets, err := FetchAssetInfo(userID)
	if err != nil {
		return user, err
	}
	if len(assets) == 0 {
		return user, nil
	}
	user.AssetDetails = assets
	return user, err
}

func ValidateUserSession(sessionID string) (bool, error) {
	SQL := `
		SELECT COUNT(*) 
		FROM user_sessions
		WHERE id=$1 AND archived_at IS NULL
	`
	var count int
	err := database.DB.Get(&count, SQL, sessionID)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func UpdateUserSession(sessionID string) error {
	SQL := `
		UPDATE user_sessions
		SET archived_at=NOW()
		WHERE id=$1 AND archived_at IS NULL
	`
	result, err := database.DB.Exec(SQL, sessionID)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("session not found")
	}
	return nil
}

func DeleteUser(tx *sqlx.Tx, userID string) error {
	SQL := `
		UPDATE users
		SET archived_at=NOW()
		WHERE id=$1 AND archived_at IS NULL
	`
	result, err := tx.Exec(SQL, userID)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("session not found")
	}
	return nil
}

func DeleteUserSession(tx *sqlx.Tx, userID string) error {
	SQL := `
		UPDATE user_sessions
		SET archived_at=NOW()
		WHERE user_id=$1 AND  archived_at IS NULL
	`

	result, err := tx.Exec(SQL, userID)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("session not found")
	}
	return nil
}
