package concerts

import (
	"time"

	"github.com/serhiirubets/rubeticket/internal/app/admin/bands"
	"github.com/serhiirubets/rubeticket/internal/app/admin/venues"
)

// @Description Concert response model
type ConcertResponse struct {
	ID          uint                 `json:"id"`
	Title       string               `json:"title"`
	Description string               `json:"description"`
	PosterURL   string               `json:"posterUrl"`
	Date        time.Time            `json:"date"`
	VenueID     uint                 `json:"venueId"`
	Venue       venues.VenueResponse `json:"venue"`
	Bands       []bands.BandResponse `json:"bands"`
	CreatedAt   time.Time            `json:"createdAt"`
	UpdatedAt   time.Time            `json:"updatedAt"`
}

// @Description Create concert request
type CreateConcertRequest struct {
	Title       string    `json:"title" validate:"required,max=100"`
	Description string    `json:"description" validate:"max=300"`
	PosterURL   string    `json:"posterUrl" validate:"max=100"`
	Date        time.Time `json:"date" validate:"required"`
	VenueID     uint      `json:"venueId" validate:"required"`
	BandIDs     []uint    `json:"bandIds" validate:"required,min=1"`
}

// @Description Update concert request
type UpdateConcertRequest struct {
	Title       *string    `json:"title" validate:"omitempty,max=100"`
	Description *string    `json:"description" validate:"omitempty,max=300"`
	PosterURL   *string    `json:"posterUrl" validate:"omitempty,max=100"`
	Date        *time.Time `json:"date"`
	VenueID     *uint      `json:"venueId"`
	BandIDs     []uint     `json:"bandIds" validate:"omitempty,min=1"`
}

// @Description List concerts response
type ListConcertsResponse struct {
	Items []ConcertResponse `json:"items"`
}
