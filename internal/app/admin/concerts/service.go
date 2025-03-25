package concerts

import (
	"errors"

	"github.com/serhiirubets/rubeticket/internal/app/admin/bands"
	"github.com/serhiirubets/rubeticket/internal/app/admin/venues"
)

type ConcertService struct {
	repository IConcertRepository
	venueRepo  venues.IVenueRepository
	bandRepo   bands.IBandRepository
}

func NewConcertService(repository IConcertRepository, venueRepo venues.IVenueRepository, bandRepo bands.IBandRepository) *ConcertService {
	return &ConcertService{
		repository: repository,
		venueRepo:  venueRepo,
		bandRepo:   bandRepo,
	}
}

func (s *ConcertService) Create(payload *CreateConcertRequest) (*ConcertResponse, error) {
	// check if venue exists
	venue, err := s.venueRepo.GetByID(payload.VenueID)
	if err != nil {
		return nil, errors.New("venue not found")
	}

	// check if all bands exist
	var bandsList []bands.Band
	for _, bandID := range payload.BandIDs {
		band, err := s.bandRepo.GetByID(bandID)
		if err != nil {
			return nil, errors.New("band not found")
		}
		bandsList = append(bandsList, *band)
	}

	concert := &Concert{
		Title:       payload.Title,
		Description: payload.Description,
		PosterURL:   payload.PosterURL,
		Date:        payload.Date,
		VenueID:     payload.VenueID,
		Venue:       *venue,
		Bands:       bandsList,
	}

	createdConcert, err := s.repository.Create(concert)
	if err != nil {
		return nil, err
	}

	bandResponses := make([]bands.BandResponse, len(createdConcert.Bands))
	for i, band := range createdConcert.Bands {
		bandResponse := bands.ToBandResponse(&band)
		bandResponses[i] = *bandResponse
	}

	response := &ConcertResponse{
		ID:          createdConcert.Model.ID,
		Title:       createdConcert.Title,
		Description: createdConcert.Description,
		PosterURL:   createdConcert.PosterURL,
		Date:        createdConcert.Date,
		VenueID:     createdConcert.VenueID,
		Venue:       *venues.ToVenueResponse(&createdConcert.Venue),
		Bands:       bandResponses,
		CreatedAt:   createdConcert.Model.CreatedAt,
		UpdatedAt:   createdConcert.Model.UpdatedAt,
	}

	return response, nil
}

func (s *ConcertService) Update(id uint, payload *UpdateConcertRequest) (*ConcertResponse, error) {
	concert, err := s.repository.GetByID(id)
	if err != nil {
		return nil, errors.New("concert not found")
	}

	if payload.Title != nil {
		concert.Title = *payload.Title
	}
	if payload.Description != nil {
		concert.Description = *payload.Description
	}
	if payload.PosterURL != nil {
		concert.PosterURL = *payload.PosterURL
	}
	if payload.Date != nil {
		concert.Date = *payload.Date
	}
	if payload.VenueID != nil {
		venue, err := s.venueRepo.GetByID(*payload.VenueID)
		if err != nil {
			return nil, errors.New("venue not found")
		}
		concert.VenueID = *payload.VenueID
		concert.Venue = *venue
	}
	if payload.BandIDs != nil {
		var bandsList []bands.Band
		for _, bandID := range payload.BandIDs {
			band, err := s.bandRepo.GetByID(bandID)
			if err != nil {
				return nil, errors.New("band not found")
			}
			bandsList = append(bandsList, *band)
		}
		concert.Bands = bandsList
	}

	err = s.repository.Update(concert)
	if err != nil {
		return nil, err
	}

	bandResponses := make([]bands.BandResponse, len(concert.Bands))
	for i, band := range concert.Bands {
		bandResponse := bands.ToBandResponse(&band)
		bandResponses[i] = *bandResponse
	}

	response := &ConcertResponse{
		ID:          concert.Model.ID,
		Title:       concert.Title,
		Description: concert.Description,
		PosterURL:   concert.PosterURL,
		Date:        concert.Date,
		VenueID:     concert.VenueID,
		Venue:       *venues.ToVenueResponse(&concert.Venue),
		Bands:       bandResponses,
		CreatedAt:   concert.Model.CreatedAt,
		UpdatedAt:   concert.Model.UpdatedAt,
	}

	return response, nil
}

func (s *ConcertService) Delete(id uint) error {
	return s.repository.Delete(id)
}

func (s *ConcertService) GetByID(id uint) (*ConcertResponse, error) {
	concert, err := s.repository.GetByID(id)
	if err != nil {
		return nil, err
	}

	bandResponses := make([]bands.BandResponse, len(concert.Bands))
	for i, band := range concert.Bands {
		bandResponse := bands.ToBandResponse(&band)
		bandResponses[i] = *bandResponse
	}

	response := &ConcertResponse{
		ID:          concert.Model.ID,
		Title:       concert.Title,
		Description: concert.Description,
		PosterURL:   concert.PosterURL,
		Date:        concert.Date,
		VenueID:     concert.VenueID,
		Venue:       *venues.ToVenueResponse(&concert.Venue),
		Bands:       bandResponses,
		CreatedAt:   concert.Model.CreatedAt,
		UpdatedAt:   concert.Model.UpdatedAt,
	}

	return response, nil
}

func (s *ConcertService) List(page, pageSize int) (*ListConcertsResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	concerts, err := s.repository.List(page, pageSize)
	if err != nil {
		return nil, err
	}

	response := &ListConcertsResponse{
		Items: make([]ConcertResponse, len(concerts)),
	}

	for i, concert := range concerts {
		bandResponses := make([]bands.BandResponse, len(concert.Bands))
		for j, band := range concert.Bands {
			bandResponse := bands.ToBandResponse(&band)
			bandResponses[j] = *bandResponse
		}

		response.Items[i] = ConcertResponse{
			ID:          concert.Model.ID,
			Title:       concert.Title,
			Description: concert.Description,
			PosterURL:   concert.PosterURL,
			Date:        concert.Date,
			VenueID:     concert.VenueID,
			Venue:       *venues.ToVenueResponse(&concert.Venue),
			Bands:       bandResponses,
			CreatedAt:   concert.Model.CreatedAt,
			UpdatedAt:   concert.Model.UpdatedAt,
		}
	}

	return response, nil
}
