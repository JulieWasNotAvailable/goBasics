package models

import "gorm.io/gorm"

type Beats struct{
	ID uint
	Author *string
	Title *string
	LicenseName *string
}

//for authomigration, because in postgres DB is not created authomatically
func MigrateBeats(db *gorm.DB) error {
	err := db.AutoMigrate(&Beats{})
	return err
}