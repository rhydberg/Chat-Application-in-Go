package models

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	Type    int
	Sender  string
	Content string
}