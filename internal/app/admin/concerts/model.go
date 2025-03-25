package concerts

import (
	"time"

	"github.com/serhiirubets/rubeticket/internal/app/admin/bands"
	"github.com/serhiirubets/rubeticket/internal/app/admin/venues"
	"gorm.io/gorm"
)

// @Description Concert model
type Concert struct {
	*gorm.Model
	Title       string       `json:"title" gorm:"type:varchar(100);not null"`
	Description string       `json:"description" gorm:"type:varchar(300)"`
	PosterURL   string       `json:"posterUrl" gorm:"type:varchar(100)"`
	Date        time.Time    `json:"date" gorm:"not null;index:idx_concert_date"`
	VenueID     uint         `json:"venueId" gorm:"not null;index:idx_concert_venue_id"`
	Venue       venues.Venue `json:"venue"`
	Bands       []bands.Band `json:"bands" gorm:"many2many:concert_bands;"`
}

// ConcertBands many to many relation model
// @Description Concert-Band relation model
type ConcertBands struct {
	ConcertID uint `json:"concertId" gorm:"primaryKey;index:idx_concert_bands"`
	BandID    uint `json:"bandId" gorm:"primaryKey;index:idx_concert_bands"`
}
