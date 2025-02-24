package users

type IUserRepository interface {
	Create(user *User) (*User, error)
	FindByEmail(email string) (*User, error)
}
