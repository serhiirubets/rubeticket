package venues

import (
	"github.com/serhiirubets/rubeticket/internal/pkg/db"
)

type IVenueRepository interface {
	Create(venue *Venue) (*Venue, error)
	Update(venue *Venue) error
	Delete(id uint) error
	GetByID(id uint) (*Venue, error)
	List(page, pageSize int) ([]Venue, error)
}

type VenueRepository struct {
	Db db.IDb
}

func NewVenueRepository(Db db.IDb) IVenueRepository {
	return &VenueRepository{Db: Db}
}

func (r *VenueRepository) Create(venue *Venue) (*Venue, error) {
	if err := r.Db.Create(venue).Error; err != nil {
		return nil, err
	}
	return venue, nil
}

func (r *VenueRepository) Update(venue *Venue) error {
	return r.Db.Save(venue).Error
}

func (r *VenueRepository) Delete(id uint) error {
	return r.Db.Delete(&Venue{}, id).Error
}

func (r *VenueRepository) GetByID(id uint) (*Venue, error) {
	var venue Venue
	if err := r.Db.First(&venue, id).Error; err != nil {
		return nil, err
	}
	return &venue, nil
}

func (r *VenueRepository) List(page, pageSize int) ([]Venue, error) {
	var venues []Venue
	offset := (page - 1) * pageSize
	if err := r.Db.Offset(offset).Limit(pageSize).Find(&venues).Error; err != nil {
		return nil, err
	}
	return venues, nil
}
