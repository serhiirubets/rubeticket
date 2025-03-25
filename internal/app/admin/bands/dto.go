package bands

import "time"

// @Description Band response model
type BandResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Genre       string    `json:"genre"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// @Description Create band request
type CreateBandRequest struct {
	Name        string `json:"name" validate:"required,max=255"`
	Description string `json:"description"`
	Genre       string `json:"genre" validate:"max=100"`
}

// @Description Update band request
type UpdateBandRequest struct {
	Name        *string `json:"name" validate:"omitempty,max=255"`
	Description *string `json:"description"`
	Genre       *string `json:"genre" validate:"omitempty,max=100"`
}

// @Description List bands response
type ListBandsResponse struct {
	Items []BandResponse `json:"items"`
}

func ToBandResponse(band *Band) *BandResponse {
	return &BandResponse{
		ID:          band.Model.ID,
		Name:        band.Name,
		Description: band.Description,
		Genre:       band.Genre,
		CreatedAt:   band.CreatedAt,
		UpdatedAt:   band.UpdatedAt,
	}
}

func ToBandResponses(bands []Band) []BandResponse {
	responses := make([]BandResponse, len(bands))
	for i, band := range bands {
		responses[i] = *ToBandResponse(&band)
	}
	return responses
}
