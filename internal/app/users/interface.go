package users

type IUserRepository interface {
	Create(user *User) (*User, error)
	GetByEmail(email string) (*User, error)
	Update(user *User, updates map[string]interface{}) error
}
