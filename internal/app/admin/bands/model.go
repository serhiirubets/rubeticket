package bands

import (
	"gorm.io/gorm"
)

// @Description Band model
type Band struct {
	*gorm.Model
	Name        string `json:"name" gorm:"type:varchar(255);not null;index:idx_band_name"`
	Description string `json:"description" gorm:"type:text"`
	Genre       string `json:"genre" gorm:"type:varchar(100);index:idx_band_genre"`
}

type ConcertBands struct {
	ConcertID uint `json:"concertId" gorm:"primaryKey;index:idx_concert_bands"`
	BandID    uint `json:"bandId" gorm:"primaryKey;index:idx_concert_bands"`
}
