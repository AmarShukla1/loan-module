package repository

import (
	"loan-module/customer/models"
	loanModels "loan-module/loan/models"
	"loan-module/repository"
)

type CustomerRepository struct {
	db *database.Database
}

func NewCustomerRepository(db *database.Database) *CustomerRepository {
	return &CustomerRepository{db: db}
}

func (r *CustomerRepository) AddCustomer(customer *models.Customer) *models.Customer {
	r.db.DB.Create(customer)
	return customer
}

func (r *CustomerRepository) GetCustomerByID(id int) (*models.Customer, bool) {
	var customer models.Customer
	result := r.db.DB.First(&customer, id)
	return &customer, result.Error == nil
}

func (r *CustomerRepository) GetCustomerByPhone(phone string) (*models.Customer, bool) {
	var customer models.Customer
	result := r.db.DB.Where("phone = ?", phone).First(&customer)
	return &customer, result.Error == nil
}

func (r *CustomerRepository) GetAllCustomers() []*models.Customer {
	var customers []*models.Customer
	r.db.DB.Find(&customers)
	return customers
}

func (r *CustomerRepository) UpdateCustomer(customer *models.Customer) {
	r.db.DB.Save(customer)
}

func (r *CustomerRepository) GetTopCustomers() []loanModels.TopCustomerResponse {
	var results []loanModels.TopCustomerResponse
	r.db.DB.Raw(`
		SELECT c.name as customer_name, COUNT(*) as approved_loans
		FROM loans l
		JOIN customers c ON l.customer_id = c.id
		WHERE l.application_status IN ('APPROVED_BY_SYSTEM', 'APPROVED_BY_AGENT')
		GROUP BY c.name
		ORDER BY approved_loans DESC
		LIMIT 3
	`).Scan(&results)
	return results
}
