package models

import "gorm.io/gorm"

type Users struct {
	gorm.Model
	Username          string `gorm:"unique"`
	Email             string `gorm:"unique"`
	Password          string
	IncorrectPassword int
	Status            int
	Role              string
}
