package users

type IUserRepository interface {
	Create(user *User) (*User, error)
	GetByEmail(email string) (*User, error)
}
