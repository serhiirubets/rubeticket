package users

import (
	"time"
)

type GetUserResponse struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Gender    Gender    `json:"gender"`
	Birthday  time.Time `json:"birthday"`
	Status    Status    `json:"status"`
	Role      Role      `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// ToResponse преобразует модель User в GetUserResponse
func (u *User) ToResponse() *GetUserResponse {
	if u == nil {
		return nil
	}

	return &GetUserResponse{
		ID:        u.ID,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Gender:    u.Gender,
		Birthday:  u.Birthday,
		Status:    u.Status,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
