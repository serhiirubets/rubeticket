package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/serhiirubets/rubeticket/internal/app/users"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepository users.IUserRepository
}

func NewAuthService(userRepository users.IUserRepository) *AuthService {
	return &AuthService{UserRepository: userRepository}
}

func (service *AuthService) Register(payload *RegisterRequest) (uint, error) {
	existedUser, _ := service.UserRepository.GetByEmail(payload.Email)

	if existedUser != nil {
		return 0, errors.New(ErrUserExists)
	}

	fromPassword, fromPassErr := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if fromPassErr != nil {
		return 0, fmt.Errorf("hashing password error: %w", fromPassErr)
	}

	birthday, birthErr := time.Parse("2006-01-02", payload.Birthday)
	if birthErr != nil {
		return 0, fmt.Errorf("wrong birthday format: %w", birthErr)
	}

	if payload.Gender != users.Male && payload.Gender != users.Female {
		return 0, errors.New("wrong gender value, expected male or female, but got: " + string(payload.Gender))
	}

	user := &users.User{
		Email:        payload.Email,
		FirstName:    payload.FirstName,
		LastName:     payload.LastName,
		Birthday:     birthday,
		PasswordHash: string(fromPassword),
		Gender:       payload.Gender,
		Role:         users.UserRole,
	}

	createdUser, dbErr := service.UserRepository.Create(user)
	if dbErr != nil {
		return 0, fmt.Errorf("creating user error: %w", dbErr)
	}
	return createdUser.ID, nil
}

func (service *AuthService) Login(email, password string) (*LoginResponseDto, error) {
	existedUser, _ := service.UserRepository.GetByEmail(email)
	if existedUser == nil {
		return nil, errors.New(ErrWrongCredentials)
	}

	err := bcrypt.CompareHashAndPassword([]byte(existedUser.PasswordHash), []byte(password))
	if err != nil {
		return nil, errors.New(ErrWrongCredentials)
	}
	return &LoginResponseDto{
		Id:    existedUser.ID,
		Email: email,
		Role:  existedUser.Role,
	}, nil
}
