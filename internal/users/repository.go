package users

import "github.com/serhiirubets/rubeticket/pkg/db"

type UserRepository struct {
	DB *db.Db
}

func NewUserRepository(db *db.Db) *UserRepository {
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

func (repo *UserRepository) FindByEmail(email string) (*User, error) {
	var user User
	result := repo.DB.First(&user, "email = ?", email)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
