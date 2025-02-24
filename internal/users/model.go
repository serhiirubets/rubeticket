package users

import (
	"gorm.io/gorm"
	"time"
)

type Status string
type Gender string

const (
	Active  Status = "active"
	Banned  Status = "banned"
	Pending Status = "pending"
	Deleted Status = "deleted"
)

const (
	Male   Gender = "Male"
	Female Gender = "Female"
)

var StatusMap = map[Status]string{
	Active:  "active",
	Banned:  "banned",
	Pending: "pending",
	Deleted: "deleted",
}

type User struct {
	*gorm.Model
	Email        string     `gorm:"type:varchar(255);uniqueIndex;not null"`
	FirstName    string     `gorm:"not null" json:"firstName"`
	LastName     string     `gorm:"not null" json:"lastName"`
	PasswordHash string     `gorm:"not null" json:"passwordHash"`
	Birthday     time.Time  `gorm:"not null" json:"birthday"`
	Gender       Gender     `gorm:"type:varchar(6);not null" json:"gender"` // (male/female)
	ActivatedAt  *time.Time `json:"activatedAt"`
	Status       Status     `gorm:"type:varchar(20);default:'pending'" json:"status"`
}
