package repository

import (
	"math"

	"gorm.io/gorm"

	"github.com/jblim0125/redredis-expire/models"
)

// Paginate for page func
func Paginate(value interface{}, pagination *models.Pagination, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	var totalRows int64
	db.Model(value).Count(&totalRows)
	pagination.TotalRows = totalRows
	totalPages := int(math.Ceil(float64(totalRows) / float64(pagination.PageSize)))
	pagination.TotalPages = totalPages
	return func(db *gorm.DB) *gorm.DB {
		if len(pagination.GetSort()) > 0 {
			return db.Offset(pagination.GetOffset()).Limit(pagination.GetPageSize()).Order(pagination.GetSort())
		}
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetPageSize())
	}
}
