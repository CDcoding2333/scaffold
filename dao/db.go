package dao

import (
	"CDcoding2333/scaffold/conf"
	"time"

	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

// Model ...
type Model struct {
	CreatedAt time.Time  `json:"created_at,omitempty"`
	UpdatedAt time.Time  `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

var db *gorm.DB

// InitDB ...
func InitDB(config *conf.DbConfig) error {
	var err error
	db, err = gorm.Open(config.Driver, config.Source)
	if err != nil {
		return err
	}

	db.LogMode(config.LogMode)
	db.DB().SetMaxOpenConns(100)
	db.DB().SetMaxIdleConns(100)
	db.DB().SetConnMaxLifetime(10 * time.Second)

	return autoMigrateTable()
}

func autoMigrateTable() error {
	if db.HasTable(&User{}) {
		if err := db.AutoMigrate(&User{}).Error; err != nil {
			log.Errorf("migrate table err %s\n", err.Error())
			return err
		}
	} else {
		if err := db.CreateTable(&User{}).Error; err != nil {
			log.Errorf("create table err %s\n", err.Error())
			return err
		}
	}

	return nil
}
