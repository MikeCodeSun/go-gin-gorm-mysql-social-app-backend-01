package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name string `json:"name" validate:"required,min=3"`
	Email string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	User_type string `json:"user_type" gorm:"default:user" validate:"eq=user|eq=admin"`
	Posts []Post `gorm:"foreignKey:UserName;references:Name"`
}

type Post struct {
	gorm.Model
	Body string `json:"body" validate:"required"`
	UserName string `json:"userName" validate:"required"`
	
}

