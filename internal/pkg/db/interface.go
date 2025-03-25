package db

import (
	"gorm.io/gorm"
)

type IDb interface {
	Create(value interface{}) *gorm.DB
	First(dest interface{}, conds ...interface{}) *gorm.DB
	Where(query interface{}, args ...interface{}) *gorm.DB
	Save(value interface{}) *gorm.DB
	Model(value interface{}) *gorm.DB
	Offset(offset int) *gorm.DB
	Delete(value interface{}, conds ...interface{}) *gorm.DB
}
