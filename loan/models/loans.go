package models

import "time"

type LoanType string
type LoanStatus string

const (
	Personal LoanType = "PERSONAL"
	Home     LoanType = "HOME"
	Auto     LoanType = "AUTO"
	Business LoanType = "BUSINESS"
)

const (
	Applied          LoanStatus = "APPLIED"
	Processing       LoanStatus = "PROCESSING"
	ApprovedBySystem LoanStatus = "APPROVED_BY_SYSTEM"
	RejectedBySystem LoanStatus = "REJECTED_BY_SYSTEM"
	UnderReview      LoanStatus = "UNDER_REVIEW"
	ApprovedByAgent  LoanStatus = "APPROVED_BY_AGENT"
	RejectedByAgent  LoanStatus = "REJECTED_BY_AGENT"
)

type Loan struct {
	ID                int        `gorm:"primaryKey" json:"loan_id"`
	CustomerID        int        `gorm:"not null;index;constraint:OnDelete:CASCADE" json:"customer_id"`
	LoanAmount        float64    `gorm:"not null" json:"loan_amount"`
	LoanType          LoanType   `gorm:"type:varchar(20);not null" json:"loan_type"`
	ApplicationStatus LoanStatus `gorm:"type:varchar(30);not null" json:"application_status"`
	CreatedAt         time.Time  `gorm:"autoCreateTime" json:"created_at"`
	AssignedAgentID   *int       `gorm:"index;constraint:OnDelete:SET NULL" json:"assigned_agent_id,omitempty"`
}

type SubmitLoanRequest struct {
	CustomerName  string   `json:"customer_name" binding:"required"`
	CustomerPhone string   `json:"customer_phone" binding:"required"`
	LoanAmount    float64  `json:"loan_amount" binding:"required,gt=0"`
	LoanType      LoanType `json:"loan_type" binding:"required"`
}

type StatusCountResponse struct {
	Status string `json:"status"`
	Count  int    `json:"count"`
}

type TopCustomerResponse struct {
	CustomerName  string `json:"customer_name"`
	ApprovedLoans int    `json:"approved_loans"`
}

type LoanAssignment struct {
	ID         int       `gorm:"primaryKey" json:"id"`
	LoanID     int       `gorm:"not null;index;constraint:OnDelete:CASCADE" json:"loan_id"`
	AgentID    int       `gorm:"not null;index;constraint:OnDelete:CASCADE" json:"agent_id"`
	AssignedAt time.Time `gorm:"autoCreateTime" json:"assigned_at"`
}
