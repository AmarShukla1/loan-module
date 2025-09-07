package models

import "time"

type Agent struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	ManagerID *int      `gorm:"index;constraint:OnDelete:SET NULL" json:"manager_id,omitempty"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

type AgentDecisionRequest struct {
	Decision string `json:"decision" binding:"required"`
}

type CreateAgentRequest struct {
	Name      string `json:"name" binding:"required"`
	ManagerID *int   `json:"manager_id"`
}
