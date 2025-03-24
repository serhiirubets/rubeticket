package users

import (
	"github.com/serhiirubets/rubeticket/internal/pkg/db"
)

type UserRepository struct {
	DB db.IDb
}

func NewUserRepository(db db.IDb) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (repo *UserRepository) Create(user *User) (*User, error) {
	createdUser := repo.DB.Create(user)

	if createdUser.Error != nil {
		return nil, createdUser.Error
	}

	return user, nil
}

func (repo *UserRepository) GetByEmail(email string) (*User, error) {
	var user User
	result := repo.DB.First(&user, "email = ?", email)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (repo *UserRepository) GetById(id string) (*User, error) {
	var user User
	result := repo.DB.First(&user, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (repo *UserRepository) Update(user *User, updates map[string]interface{}) error {
	return repo.DB.Model(user).Updates(updates).Error
}
