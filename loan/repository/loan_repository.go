package repository

import (
	"time"

	"loan-module/loan/models"
	"loan-module/repository"
)

type LoanRepository struct {
	db *database.Database
}

func NewLoanRepository(db *database.Database) *LoanRepository {
	return &LoanRepository{db: db}
}

func (r *LoanRepository) AddLoan(loan *models.Loan) *models.Loan {
	loan.ApplicationStatus = models.Applied
	r.db.DB.Create(loan)
	return loan
}

func (r *LoanRepository) GetLoanByID(id int) (*models.Loan, bool) {
	var loan models.Loan
	result := r.db.DB.First(&loan, id)
	return &loan, result.Error == nil
}

func (r *LoanRepository) GetLoansByStatus(status models.LoanStatus) []*models.Loan {
	var loans []*models.Loan
	r.db.DB.Where("application_status = ?", status).Find(&loans)
	return loans
}

func (r *LoanRepository) GetAllLoans() []*models.Loan {
	var loans []*models.Loan
	r.db.DB.Find(&loans)
	return loans
}

func (r *LoanRepository) UpdateLoan(loan *models.Loan) {
	r.db.DB.Save(loan)
}

func (r *LoanRepository) GetStatusCount() map[models.LoanStatus]int {
	var loans []models.Loan
	r.db.DB.Find(&loans)

	counts := make(map[models.LoanStatus]int)
	for _, loan := range loans {
		counts[loan.ApplicationStatus]++
	}
	return counts
}

func (r *LoanRepository) AssignLoanToAgent(loan *models.Loan, agentID int) {
	assignment := models.LoanAssignment{
		LoanID:     loan.ID,
		AgentID:    agentID,
		AssignedAt: time.Now(),
	}
	r.db.DB.Create(&assignment)
}
