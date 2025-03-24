package models

import "gorm.io/gorm"

// For testing purposes
// TODO : Remove this after proper db architecture is implemented

type User struct {
	gorm.Model
	Name  string `json:"firstname"`
	Email string `json:"lastname"`
} 
