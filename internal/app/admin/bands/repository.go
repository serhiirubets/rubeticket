package bands

import (
	"github.com/serhiirubets/rubeticket/internal/pkg/db"
)

type IBandRepository interface {
	Create(band *Band) (*Band, error)
	Update(band *Band) error
	Delete(id uint) error
	GetByID(id uint) (*Band, error)
	List(page, pageSize int) ([]Band, error)
}

type BandRepository struct {
	Db db.IDb
}

func NewBandRepository(Db db.IDb) IBandRepository {
	return &BandRepository{Db: Db}
}

func (r *BandRepository) Create(band *Band) (*Band, error) {
	if err := r.Db.Create(band).Error; err != nil {
		return nil, err
	}
	return band, nil
}

func (r *BandRepository) Update(band *Band) error {
	return r.Db.Save(band).Error
}

func (r *BandRepository) Delete(id uint) error {
	return r.Db.Delete(&Band{}, id).Error
}

func (r *BandRepository) GetByID(id uint) (*Band, error) {
	var band Band
	if err := r.Db.First(&band, id).Error; err != nil {
		return nil, err
	}
	return &band, nil
}

func (r *BandRepository) List(page, pageSize int) ([]Band, error) {
	var bands []Band
	offset := (page - 1) * pageSize
	if err := r.Db.Offset(offset).Limit(pageSize).Find(&bands).Error; err != nil {
		return nil, err
	}
	return bands, nil
}
