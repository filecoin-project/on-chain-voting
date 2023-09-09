package model

import "time"

type Dict struct {
	Id         int64     `json:"id" gorm:"not null"`
	Name       string    `json:"name" gorm:"unique;not null"`
	Value      string    `json:"value" gorm:"not null"`
	CreateTime time.Time `json:"create_time" gorm:"not null;"`
	UpdateTime time.Time `json:"update_time" gorm:"not null"`
}
