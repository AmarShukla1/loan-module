package service

import (
	"errors"

	"loan-module/agent/repository"
	customerRepo "loan-module/customer/repository"
	loanModels "loan-module/loan/models"
	loanRepo "loan-module/loan/repository"
	"loan-module/notification"
)

type AgentService struct {
	repo                *repository.AgentRepository
	loanRepo            *loanRepo.LoanRepository
	customerRepo        *customerRepo.CustomerRepository
	notificationService *notification.NotificationService
}

func NewAgentService(
	repo *repository.AgentRepository,
	loanRepo *loanRepo.LoanRepository,
	customerRepo *customerRepo.CustomerRepository,
	notificationService *notification.NotificationService,
) *AgentService {
	return &AgentService{
		repo:                repo,
		loanRepo:            loanRepo,
		customerRepo:        customerRepo,
		notificationService: notificationService,
	}
}

func (s *AgentService) MakeDecision(agentID, loanID int, decision string) (*loanModels.Loan, error) {
	loan, exists := s.loanRepo.GetLoanByID(loanID)
	if !exists {
		return nil, errors.New("loan not found")
	}
	if loan.AssignedAgentID == nil || *loan.AssignedAgentID != agentID {
		return nil, errors.New("loan not assigned to this agent")
	}
	if loan.ApplicationStatus != loanModels.UnderReview {
		return nil, errors.New("loan is not under review")
	}
	_, exists = s.repo.GetAgentByID(agentID)
	if !exists {
		return nil, errors.New("agent not found")
	}

	// Get customer phone for notification
	customer, customerExists := s.customerRepo.GetCustomerByID(loan.CustomerID)
	if !customerExists {
		return nil, errors.New("customer not found")
	}

	switch decision {
	case "APPROVE":
		loan.ApplicationStatus = loanModels.ApprovedByAgent
		s.notificationService.SendSMS(customer.Phone, "Your loan has been approved by our agent.")
	case "REJECT":
		loan.ApplicationStatus = loanModels.RejectedByAgent
		s.notificationService.SendSMS(customer.Phone, "loan has been rejected after review.")
	default:
		return nil, errors.New("invalid decision. Must be APPROVE or REJECT")
	}
	s.loanRepo.UpdateLoan(loan)
	return loan, nil
}
