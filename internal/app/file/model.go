package file

import (
	"gorm.io/gorm"
)

type File struct {
	*gorm.Model
	UUID     string `gorm:"unique;not null"`
	UserID   uint   `gorm:"not null"`
	FilePath string `gorm:"not null"`
	Purpose  string
}
