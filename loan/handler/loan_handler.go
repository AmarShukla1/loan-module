package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"loan-module/loan/models"
	"loan-module/loan/service"
)

type LoanHandler struct {
	loanService *service.LoanService
}

func NewLoanHandler(loanService *service.LoanService) *LoanHandler {
	return &LoanHandler{loanService: loanService}
}

func (h *LoanHandler) SubmitLoan(c *gin.Context) {
	var req models.SubmitLoanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	loan := h.loanService.SubmitLoan(&req)
	c.JSON(http.StatusCreated, loan)
}

func (h *LoanHandler) GetStatusCount(c *gin.Context) {
	counts := h.loanService.GetStatusCount()
	c.JSON(http.StatusOK, counts)
}

func (h *LoanHandler) GetLoansByStatus(c *gin.Context) {
	status := models.LoanStatus(c.Query("status"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 10
	}
	loans := h.loanService.GetLoansByStatus(status, page, size)
	response := gin.H{"loans": loans, "page": page, "size": size, "total": len(loans)}
	c.JSON(http.StatusOK, response)
}

func (h *LoanHandler) GetLoanByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid loan ID"})
		return
	}
	loan, exists := h.loanService.GetLoanByID(id)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Loan not found"})
		return
	}
	c.JSON(http.StatusOK, loan)
}
