package service

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	agent "loan-module/agent/repository"
	"loan-module/customer/models"
	customer "loan-module/customer/repository"
	loanModels "loan-module/loan/models"
	"loan-module/loan/repository"
	"loan-module/notification"
)

// Worker pool config
const workerCount = 5

type LoanService struct {
	repo                *repository.LoanRepository
	agentRepo           *agent.AgentRepository
	customerRepo        *customer.CustomerRepository
	notificationService *notification.NotificationService
}

func NewLoanService(
	repo *repository.LoanRepository,
	agentRepo *agent.AgentRepository,
	customerRepo *customer.CustomerRepository,
	notificationService *notification.NotificationService,
) *LoanService {
	return &LoanService{
		repo:                repo,
		agentRepo:           agentRepo,
		customerRepo:        customerRepo,
		notificationService: notificationService,
	}
}

func (s *LoanService) SubmitLoan(req *loanModels.SubmitLoanRequest) *loanModels.Loan {
	// Check if customer exists by phone number
	customer, exists := s.customerRepo.GetCustomerByPhone(req.CustomerPhone)

	if !exists {
		// Create new customer
		newCustomer := &models.Customer{
			Name:  req.CustomerName,
			Phone: req.CustomerPhone,
		}
		customer = s.customerRepo.AddCustomer(newCustomer)
	}

	loan := &loanModels.Loan{
		CustomerID: customer.ID,
		LoanAmount: req.LoanAmount,
		LoanType:   req.LoanType,
	}
	return s.repo.AddLoan(loan)
}

func (s *LoanService) StartLoanProcessor() {
	log.Println("Starting loan processor with worker pool...")

	jobChan := make(chan *loanModels.Loan, 100)

	// Start workers
	for i := 0; i < workerCount; i++ {
		go func(workerID int) {
			for loan := range jobChan {
				log.Printf("[Worker %d] Processing loan %d", workerID, loan.ID)
				s.processLoan(loan)
			}
		}(i + 1)
	}

	// Feed jobs periodically
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C
		loans := s.repo.GetLoansByStatus(loanModels.Applied)
		for _, loan := range loans {
			loan.ApplicationStatus = loanModels.Processing
			s.repo.UpdateLoan(loan)
			jobChan <- loan
		}
	}
}

func (s *LoanService) processLoan(loan *loanModels.Loan) {
	delay := time.Duration(rand.Intn(20)+5) * time.Second
	log.Printf("Processing loan %d, waiting %v seconds...", loan.ID, delay.Seconds())
	time.Sleep(delay)

	var newStatus loanModels.LoanStatus

	// Get customer for notification
	customer, exists := s.customerRepo.GetCustomerByID(loan.CustomerID)
	if !exists {
		log.Printf("Customer not found for loan %d", loan.ID)
		return
	}

	switch {
	case loan.LoanAmount < 10000:
		newStatus = loanModels.ApprovedBySystem
		s.notificationService.SendSMS(customer.Phone, "Your loan has been approved by system.")
	case loan.LoanAmount > 500000:
		newStatus = loanModels.RejectedBySystem
		s.notificationService.SendSMS(customer.Phone, "loan application has been rejected by system.")
	default:
		newStatus = loanModels.UnderReview
		s.assignToAgent(loan, customer)
	}

	loan.ApplicationStatus = newStatus
	s.repo.UpdateLoan(loan)
}

func (s *LoanService) assignToAgent(loan *loanModels.Loan, customer *models.Customer) {
	agent := s.agentRepo.GetAvailableAgent()
	if agent == nil {
		log.Printf("No available agent for loan %d", loan.ID)
		return
	}

	// Assign the loan to agent first
	loan.AssignedAgentID = &agent.ID
	s.repo.UpdateLoan(loan)

	// Create assignment record
	s.repo.AssignLoanToAgent(loan, agent.ID)

	// Send notifications
	s.notificationService.SendPushNotification(agent.ID,
		fmt.Sprintf("New loan #%d assigned to you for review", loan.ID))

	if agent.ManagerID != nil {
		s.notificationService.SendPushNotification(*agent.ManagerID,
			fmt.Sprintf("Loan #%d assigned to your team member %s", loan.ID, agent.Name))
	}

	log.Printf("Loan %d assigned to agent %d (%s)", loan.ID, agent.ID, agent.Name)
}

func (s *LoanService) GetStatusCount() []loanModels.StatusCountResponse {
	counts := s.repo.GetStatusCount()
	var result []loanModels.StatusCountResponse
	allStatuses := []loanModels.LoanStatus{
		loanModels.Applied, loanModels.Processing, loanModels.ApprovedBySystem, loanModels.RejectedBySystem,
		loanModels.UnderReview, loanModels.ApprovedByAgent, loanModels.RejectedByAgent,
	}
	for _, status := range allStatuses {
		result = append(result, loanModels.StatusCountResponse{Status: string(status), Count: counts[status]})
	}
	return result
}

func (s *LoanService) GetLoansByStatus(status loanModels.LoanStatus, page, size int) []*loanModels.Loan {
	var loans []*loanModels.Loan
	if status == "" {
		loans = s.repo.GetAllLoans()
	} else {
		loans = s.repo.GetLoansByStatus(status)
	}
	start := (page - 1) * size
	end := start + size
	if start >= len(loans) {
		return []*loanModels.Loan{}
	}
	if end > len(loans) {
		end = len(loans)
	}
	return loans[start:end]
}

func (s *LoanService) GetLoanByID(id int) (*loanModels.Loan, bool) {
	return s.repo.GetLoanByID(id)
}

func (s *LoanService) UpdateLoan(loan *loanModels.Loan) {
	s.repo.UpdateLoan(loan)
}

func (s *LoanService) GetTopCustomers() []loanModels.TopCustomerResponse {
	return s.customerRepo.GetTopCustomers()
}
