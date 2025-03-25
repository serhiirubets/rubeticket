package concerts

import (
	"github.com/serhiirubets/rubeticket/internal/pkg/db"
)

type IConcertRepository interface {
	Create(concert *Concert) (*Concert, error)
	Update(concert *Concert) error
	Delete(id uint) error
	GetByID(id uint) (*Concert, error)
	List(page, pageSize int) ([]Concert, error)
}

type ConcertRepository struct {
	Db db.IDb
}

func NewConcertRepository(Db db.IDb) IConcertRepository {
	return &ConcertRepository{Db: Db}
}

func (r *ConcertRepository) Create(concert *Concert) (*Concert, error) {
	if err := r.Db.Create(concert).Error; err != nil {
		return nil, err
	}
	return concert, nil
}

func (r *ConcertRepository) Update(concert *Concert) error {
	return r.Db.Save(concert).Error
}

func (r *ConcertRepository) Delete(id uint) error {
	return r.Db.Delete(&Concert{}, id).Error
}

func (r *ConcertRepository) GetByID(id uint) (*Concert, error) {
	var concert Concert
	if err := r.Db.First(&concert, id).Error; err != nil {
		return nil, err
	}
	return &concert, nil
}

func (r *ConcertRepository) List(page, pageSize int) ([]Concert, error) {
	var concerts []Concert
	offset := (page - 1) * pageSize
	if err := r.Db.Offset(offset).Limit(pageSize).Find(&concerts).Error; err != nil {
		return nil, err
	}
	return concerts, nil
}
