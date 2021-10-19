package models

import (
	"gorm.io/gorm"
)

type AutomationTemplates struct {
	gorm.Model
	ProjectName      string
	Technology       string
	UserID           int
	User             Users `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	TemplateFile     []byte
	TemplateMetaData []byte
}

func (a AutomationTemplates) GetTemplateFile() string {
	return string(a.TemplateFile)
}

type AutomationMetadata struct {
	Name        string `json:"variable"`
	Hint        string `json:"hint"`
	Placeholder string `json:"placeholder"`
}
