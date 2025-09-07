package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"loan-module/providers"
	"log"
	"net/http"

	agentHandler "loan-module/agent/handler"
	agentRepo "loan-module/agent/repository"
	agentService "loan-module/agent/service"

	customerHandler "loan-module/customer/handler"
	customerRepo "loan-module/customer/repository"
	customerService "loan-module/customer/service"

	loanHandler "loan-module/loan/handler"
	loanRepo "loan-module/loan/repository"
	loanService "loan-module/loan/service"

	agentModels "loan-module/agent/models"
	"loan-module/notification"
	database "loan-module/repository"
)

func main() {
	// Create a root context with cancellation
	rootCtx, rootCancel := context.WithCancel(context.Background())
	defer rootCancel()

	// Load configuration from YAML file
	config, err := providers.GetConfig("loan-module-configuration.yaml")
	if err != nil {
		log.Fatal("Failed to load configuration: ", err)
	}

	// Initialize database with configuration
	db := database.NewDatabaseWithConfig(config)

	// Initialize repositories
	customerRepository := customerRepo.NewCustomerRepository(db)
	agentRepository := agentRepo.NewAgentRepository(db)
	loanRepository := loanRepo.NewLoanRepository(db)

	// Initialize notification
	notificationService := notification.NewNotificationService()
	customerService := customerService.NewCustomerService(customerRepository)
	loanService := loanService.NewLoanService(loanRepository, agentRepository, customerRepository, notificationService)
	agentService := agentService.NewAgentService(agentRepository, loanRepository, customerRepository, notificationService)

	// Initialize handlers
	customerHandler := customerHandler.NewCustomerHandler(customerService)
	loanHandler := loanHandler.NewLoanHandler(loanService)
	agentHandler := agentHandler.NewAgentHandler(agentService)

	// Initialize sample data
	initSampleData(agentRepository)

	// Start loan processor with context
	go loanService.StartLoanProcessor(rootCtx)

	// Setup router
	router := gin.Default()
	v1 := router.Group("/api/v1")
	{
		// Customer endpoints
		v1.POST("/customers", customerHandler.CreateCustomer)
		v1.GET("/customers/:id", customerHandler.GetCustomerByID)
		v1.GET("/customers", customerHandler.GetAllCustomers)
		v1.GET("/customers/top", customerHandler.GetTopCustomers)

		// Loan endpoints
		v1.POST("/loans", loanHandler.SubmitLoan)
		v1.GET("/loans/status-count", loanHandler.GetStatusCount)
		v1.GET("/loans", loanHandler.GetLoansByStatus)
		v1.GET("/loans/:id", loanHandler.GetLoanByID)

		// Agent endpoints
		v1.POST("/agents", agentHandler.CreateAgent)
		v1.PUT("/agents/:agent_id/loans/:loan_id/decision", agentHandler.MakeDecision)
	}

	fmt.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func initSampleData(agentRepo *agentRepo.AgentRepository) {
	// Add sample agents
	agentRepo.AddAgent(&agentModels.Agent{ID: 1, Name: "John Manager", ManagerID: nil})
	agentRepo.AddAgent(&agentModels.Agent{ID: 2, Name: "Alice Agent", ManagerID: &[]int{1}[0]})
	agentRepo.AddAgent(&agentModels.Agent{ID: 3, Name: "Bob Agent", ManagerID: &[]int{1}[0]})

	// Update the sequence to prevent primary key conflicts
	agentRepo.DB.DB.Exec("SELECT setval('agents_id_seq', (SELECT MAX(id) FROM agents))")
}
