package file

import (
	"github.com/serhiirubets/rubeticket/internal/pkg/db"
)

type Repository struct {
	Db db.IDb
}

func NewRepository(Db db.IDb) IFileRepository {
	return &Repository{
		Db: Db,
	}
}

func (repo *Repository) Create(file *File) (*File, error) {
	createdFile := repo.Db.Create(file)
	if createdFile.Error != nil {
		return nil, createdFile.Error
	}
	return file, nil
}

func (repo *Repository) GetById(id string) (*File, error) {
	var file File
	result := repo.Db.First(&file, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &file, nil
}

func (repo *Repository) CreateWithStorage(file *File) (*File, error) {
	if err := repo.Db.Create(file).Error; err != nil {
		return nil, err
	}
	return file, nil
}
