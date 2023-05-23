package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email    string `gorm:"unique"`
	Password string `gorm:"column:passowrd"`   //gambiarra para o gorm mapear corretamente pois jumenta escreveu errado :P
}