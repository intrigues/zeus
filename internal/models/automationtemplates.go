package models

import (
	"gorm.io/gorm"
)

type AutomationTemplates struct {
	gorm.Model
	ID          string
	ProjectName string
	Technology  string
	UserID      int
	User        Users `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type AutomationMetadata struct {
	Name        string `json:"variable"`
	Description string `json:"description"`
}
