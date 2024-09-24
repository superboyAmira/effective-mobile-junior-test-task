package postgresql

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() (*gorm.DB, error) {
	dsn := newDSN();
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = Ping(db)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func Ping(db *gorm.DB) error {
	dbSql, err := db.DB()
	if err != nil {
		return err
	}
	if err = dbSql.Ping(); err != nil {
		return err
	}
	return nil
}

// обертка миграции данных
func Migrate(db *gorm.DB, models ...any) error {
	return db.AutoMigrate(models...)
}

func TxSaveExecutor(db *gorm.DB, fn func(*gorm.DB) error) error {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	dbSql, err := db.DB()
	if err != nil {
		tx.Rollback()
		return err
	}
	if err = dbSql.Ping(); err != nil {
		tx.Rollback()
		return err
	}
	if err = fn(db); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}


func newDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"))
}