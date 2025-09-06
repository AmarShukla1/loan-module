package models

import "time"

type Customer struct {
	ID        int       `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Phone     string    `gorm:"not null;uniqueIndex:idx_customers_phone" json:"phone"`
	Email     string    `json:"email,omitempty"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

type CreateCustomerRequest struct {
	Name  string `json:"name" binding:"required"`
	Phone string `json:"phone" binding:"required"`
	Email string `json:"email"`
}

type CustomerResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Email string `json:"email,omitempty"`
}
