package service

import (
	"loan-module/customer/models"
	"loan-module/customer/repository"
	loanModels "loan-module/loan/models"
)

type CustomerService struct {
	repo *repository.CustomerRepository
}

func NewCustomerService(repo *repository.CustomerRepository) *CustomerService {
	return &CustomerService{repo: repo}
}

func (s *CustomerService) CreateCustomer(req *models.CreateCustomerRequest) *models.Customer {
	customer := &models.Customer{
		Name:  req.Name,
		Phone: req.Phone,
		Email: req.Email,
	}
	return s.repo.AddCustomer(customer)
}

func (s *CustomerService) GetCustomerByID(id int) (*models.Customer, bool) {
	return s.repo.GetCustomerByID(id)
}

func (s *CustomerService) GetAllCustomers() []*models.Customer {
	return s.repo.GetAllCustomers()
}

func (s *CustomerService) GetTopCustomers() []loanModels.TopCustomerResponse {
	return s.repo.GetTopCustomers()
}
