package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"loan-module/customer/models"
	"loan-module/customer/service"
)

type CustomerHandler struct {
	customerService *service.CustomerService
}

func NewCustomerHandler(customerService *service.CustomerService) *CustomerHandler {
	return &CustomerHandler{customerService: customerService}
}

func (h *CustomerHandler) CreateCustomer(c *gin.Context) {
	var req models.CreateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	customer := h.customerService.CreateCustomer(&req)
	c.JSON(http.StatusCreated, customer)
}

func (h *CustomerHandler) GetCustomerByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID"})
		return
	}
	customer, exists := h.customerService.GetCustomerByID(id)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}
	c.JSON(http.StatusOK, customer)
}

func (h *CustomerHandler) GetAllCustomers(c *gin.Context) {
	customers := h.customerService.GetAllCustomers()
	c.JSON(http.StatusOK, gin.H{"customers": customers})
}

func (h *CustomerHandler) GetTopCustomers(c *gin.Context) {
	customers := h.customerService.GetTopCustomers()
	c.JSON(http.StatusOK, gin.H{"top_customers": customers})
}
