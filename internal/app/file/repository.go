package file

import (
	"github.com/serhiirubets/rubeticket/internal/pkg/db"
)

type Repository struct {
	Db *db.Db
}

func NewRepository(Db *db.Db) *Repository {
	return &Repository{
		Db: Db,
	}
}
