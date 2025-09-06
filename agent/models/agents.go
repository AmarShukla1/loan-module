package models

import "time"

type Agent struct {
	ID        int       `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	ManagerID *int      `gorm:"index;constraint:OnDelete:SET NULL" json:"manager_id,omitempty"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

type AgentDecisionRequest struct {
	Decision string `json:"decision" binding:"required"`
}
