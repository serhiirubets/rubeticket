package bands

import "errors"

type BandService struct {
	repository IBandRepository
}

func NewBandService(repository IBandRepository) *BandService {
	return &BandService{repository: repository}
}

func (s *BandService) Create(payload *CreateBandRequest) (*Band, error) {
	band := &Band{
		Name:        payload.Name,
		Description: payload.Description,
	}

	return s.repository.Create(band)
}

func (s *BandService) Update(id uint, payload *UpdateBandRequest) (*Band, error) {
	band, err := s.repository.GetByID(id)
	if err != nil {
		return nil, errors.New("band not found")
	}

	if payload.Name != nil {
		band.Name = *payload.Name
	}
	if payload.Description != nil {
		band.Description = *payload.Description
	}

	err = s.repository.Update(band)
	if err != nil {
		return nil, err
	}

	return band, nil
}

func (s *BandService) Delete(id uint) error {
	return s.repository.Delete(id)
}

func (s *BandService) GetByID(id uint) (*Band, error) {
	return s.repository.GetByID(id)
}

func (s *BandService) List(page, pageSize int) ([]Band, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	return s.repository.List(page, pageSize)
}
