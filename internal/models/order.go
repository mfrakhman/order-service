package models

import "time"

type Order struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductID  string    `gorm:"type:uuid;not null" json:"productId"`
	TotalPrice int       `gorm:"not null" json:"totalPrice"`
	Status     string    `gorm:"type:varchar(20);not null;default:'PENDING'" json:"status"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"createdAt"`
}
