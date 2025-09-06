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

func (r *LoanRepository) AddLoan(loan *models.Loan) (*models.Loan, error) {
	tx := r.db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	loan.ApplicationStatus = models.Applied
	if err := tx.Create(loan).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return loan, nil
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

func (r *LoanRepository) UpdateLoan(loan *models.Loan) error {
	tx := r.db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Save(loan).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
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

func (r *LoanRepository) AssignLoanToAgent(loan *models.Loan, agentID int) error {
	tx := r.db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Update loan with agent ID and status
	if err := tx.Model(loan).Updates(map[string]interface{}{
		"assigned_agent_id": agentID,
		"application_status": models.UnderReview,
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Create assignment record
	assignment := models.LoanAssignment{
		LoanID:     loan.ID,
		AgentID:    agentID,
		AssignedAt: time.Now(),
	}
	if err := tx.Create(&assignment).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
