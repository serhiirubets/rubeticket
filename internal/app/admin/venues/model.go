package venues

import (
	"gorm.io/gorm"
)

// @Description Venue model
type Venue struct {
	*gorm.Model
	Name        string `json:"name" gorm:"type:varchar(100);not null;index:idx_venue_name"`
	Description string `json:"description" gorm:"type:varchar(300)"`
	Address     string `json:"address" gorm:"type:varchar(50);not null;index:idx_venue_address"`
	Phone       string `json:"phone" gorm:"type:varchar(20)"`
	Email       string `json:"email" gorm:"type:varchar(50)"`
}
