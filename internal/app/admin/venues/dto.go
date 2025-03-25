package venues

import "time"

// @Description Venue response model
type VenueResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Address     string    `json:"address"`
	Phone       string    `json:"phone"`
	Email       string    `json:"email"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// @Description Create venue request
type CreateVenueRequest struct {
	Name        string `json:"name" validate:"required,max=100"`
	Description string `json:"description" validate:"max=300"`
	Address     string `json:"address" validate:"required,max=50"`
	Phone       string `json:"phone" validate:"max=20"`
	Email       string `json:"email" validate:"max=50,email"`
}

// @Description Update venue request
type UpdateVenueRequest struct {
	Name        *string `json:"name" validate:"omitempty,max=100"`
	Description *string `json:"description" validate:"omitempty,max=300"`
	Address     *string `json:"address" validate:"omitempty,max=50"`
	Phone       *string `json:"phone" validate:"omitempty,max=20"`
	Email       *string `json:"email" validate:"omitempty,max=50,email"`
}

// @Description List venues response
type ListVenuesResponse struct {
	Items []VenueResponse `json:"items"`
}

// ToVenueResponse converts from Venue to VenueResponse
func ToVenueResponse(venue *Venue) *VenueResponse {
	return &VenueResponse{
		ID:          venue.Model.ID,
		Name:        venue.Name,
		Description: venue.Description,
		Address:     venue.Address,
		Phone:       venue.Phone,
		Email:       venue.Email,
		CreatedAt:   venue.CreatedAt,
		UpdatedAt:   venue.UpdatedAt,
	}
}

func ToVenuesResponse(venues []Venue) []VenueResponse {
	responses := make([]VenueResponse, len(venues))
	for i, venue := range venues {
		responses[i] = *ToVenueResponse(&venue)
	}
	return responses
}
