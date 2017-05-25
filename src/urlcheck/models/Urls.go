package models

import (
    "github.com/jinzhu/gorm"
)

// Urls describes the storage model for our data.
type Urls struct {
	gorm.Model

	Hostname string `json:"hostname" gorm:"not null"`
	Path     string `json:"path"     gorm:"not null"`
}

// SQL Support
// TableName allows us to customize the table name.
func (Urls) TableName() string {
  return "urls"
}
