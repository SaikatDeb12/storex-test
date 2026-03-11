package database

import (
	"errors"
	"fmt"

	"github.com/SaikatDeb12/storeX/internal/utils"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jmoiron/sqlx"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

type SSLMODETYPE string

const sslmode SSLMODETYPE = "disable"

func migrateUp(db *sqlx.DB) error {
	driver, driverErr := postgres.WithInstance(db.DB, &postgres.Config{})
	if driverErr != nil {
		return driverErr
	}

	migrateIns, migrateErr := migrate.NewWithDatabaseInstance("file://internal/database/migrations", "postgres", driver)
	if migrateErr != nil {
		return migrateErr
	}

	if err := migrateIns.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

func Connect() error {
	db_user := utils.GetEnvVariables("DB_USER")
	db_password := utils.GetEnvVariables("DB_PASSWORD")
	db_name := utils.GetEnvVariables("DB_NAME")
	db_host := utils.GetEnvVariables("DB_HOST")
	db_port := utils.GetEnvVariables("DB_PORT")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s ", db_host, db_port, db_user, db_password, db_name, sslmode)
	var err error
	DB, err = sqlx.Connect("postgres", connStr)
	if err != nil {
		return err
	}

	fmt.Println("Database successfully connected")
	return migrateUp(DB)
}

func Tx(fn func(tx *sqlx.Tx) error) error {
	tx, err := DB.Beginx()
	if err != nil {
		return fmt.Errorf("failed to start a transaction: %+v", err)
	}
	defer func() {
		if err != nil {
			if rollBackErr := tx.Rollback(); rollBackErr != nil {
				fmt.Printf("failed to rollback tx: %s", rollBackErr)
			}
			return
		}
		if commitErr := tx.Commit(); commitErr != nil {
			fmt.Printf("failed to commit tx: %s", commitErr)
		}
	}()
	err = fn(tx)
	return err
}
