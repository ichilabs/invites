package models

import "gorm.io/gorm"

// For testing purposes
// TODO : Remove this after proper db architecture is implemented

type User struct {
	gorm.Model
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
} 
