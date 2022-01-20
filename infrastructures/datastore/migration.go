package datastore

import (
	"gorm.io/gorm"

	"github.com/jblim0125/redredis-expire/models"
	"github.com/jblim0125/redredis-expire/tools/util"
)

// Migration 데이터 베이스 마이그레이션을 위한 구조체
type Migration struct {
	Version       int
	StrVersion    string
	MigrationFunc func(db *gorm.DB, m *Migration) error
}

var migration = []Migration{
	//{
	//	Version:       1,
	//	StrVersion:    "1.0.0",
	//	MigrationFunc: V1_0_0,
	//},
}

func writeVersion(db *gorm.DB, m *Migration) error {
	result := db.Create(&models.Version{
		Version: m.Version,
		WriteAt: util.GetMillis(),
	})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// V1_0_0 v1.0.0 마이그레이션
func V1_0_0(db *gorm.DB, m *Migration) error {
	//var err error
	//err = db.AutoMigrate(&models.Version{})
	//if err != nil {
	//	return err
	//}
	//err = db.AutoMigrate(&models.ConfigurationDAO{})
	//if err != nil {
	//	return err
	//}
	//// Insert Default Value
	//defConf := models.ConfigurationDAO{
	//	CreatedAt:     util.GetMillis(),
	//	Configuration: models.DefaultConf,
	//}
	//tx := db.Create(&defConf)
	//if tx.Error != nil {
	//	return tx.Error
	//}
	//// Write Version
	//return writeVersion(db, m)
	return nil
}
