package models

// Version For Database Migration
type Version struct {
	Version int   `json:"version" gorm:"column:version;not null"`
	WriteAt int64 `json:"writeAt" gorm:"not null;type:BIGINT"`
}
