package venues

import (
	"errors"
)

type VenueService struct {
	repository IVenueRepository
}

func NewVenueService(repository IVenueRepository) *VenueService {
	return &VenueService{repository: repository}
}

func (s *VenueService) Create(payload *CreateVenueRequest) (*VenueResponse, error) {
	venue := &Venue{
		Name:        payload.Name,
		Description: payload.Description,
		Address:     payload.Address,
		Phone:       payload.Phone,
		Email:       payload.Email,
	}

	created, err := s.repository.Create(venue)
	if err != nil {
		return nil, err
	}

	return &VenueResponse{
		ID:          created.ID,
		Name:        created.Name,
		Description: created.Description,
		Address:     created.Address,
		Phone:       created.Phone,
		Email:       created.Email,
		CreatedAt:   created.CreatedAt,
		UpdatedAt:   created.UpdatedAt,
	}, nil
}

func (s *VenueService) Update(id uint, payload *UpdateVenueRequest) (*VenueResponse, error) {
	venue, err := s.repository.GetByID(id)
	if err != nil {
		return nil, errors.New("venue not found")
	}

	if payload.Name != nil {
		venue.Name = *payload.Name
	}
	if payload.Description != nil {
		venue.Description = *payload.Description
	}
	if payload.Address != nil {
		venue.Address = *payload.Address
	}
	if payload.Phone != nil {
		venue.Phone = *payload.Phone
	}
	if payload.Email != nil {
		venue.Email = *payload.Email
	}

	err = s.repository.Update(venue)
	if err != nil {
		return nil, err
	}

	return &VenueResponse{
		ID:          venue.ID,
		Name:        venue.Name,
		Description: venue.Description,
		Address:     venue.Address,
		Phone:       venue.Phone,
		Email:       venue.Email,
		CreatedAt:   venue.CreatedAt,
		UpdatedAt:   venue.UpdatedAt,
	}, nil
}

func (s *VenueService) Delete(id uint) error {
	return s.repository.Delete(id)
}

func (s *VenueService) GetByID(id uint) (*VenueResponse, error) {
	venue, err := s.repository.GetByID(id)
	if err != nil {
		return nil, err
	}

	return &VenueResponse{
		ID:          venue.ID,
		Name:        venue.Name,
		Description: venue.Description,
		Address:     venue.Address,
		Phone:       venue.Phone,
		Email:       venue.Email,
		CreatedAt:   venue.CreatedAt,
		UpdatedAt:   venue.UpdatedAt,
	}, nil
}

func (s *VenueService) List(page, pageSize int) (*ListVenuesResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	venues, err := s.repository.List(page, pageSize)
	if err != nil {
		return nil, err
	}

	response := &ListVenuesResponse{
		Items: make([]VenueResponse, len(venues)),
	}

	for i, venue := range venues {
		response.Items[i] = VenueResponse{
			ID:          venue.ID,
			Name:        venue.Name,
			Description: venue.Description,
			Address:     venue.Address,
			Phone:       venue.Phone,
			Email:       venue.Email,
			CreatedAt:   venue.CreatedAt,
			UpdatedAt:   venue.UpdatedAt,
		}
	}

	return response, nil
}
