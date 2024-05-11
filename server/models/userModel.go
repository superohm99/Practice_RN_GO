package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Id   uint `gorm:"primaryKey"`
	Name string
}
