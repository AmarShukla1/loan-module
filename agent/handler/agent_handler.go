package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"loan-module/agent/models"
	"loan-module/agent/service"
)

type AgentHandler struct {
	agentService *service.AgentService
}

func NewAgentHandler(agentService *service.AgentService) *AgentHandler {
	return &AgentHandler{agentService: agentService}
}

func (h *AgentHandler) MakeDecision(c *gin.Context) {
	agentID, err := strconv.Atoi(c.Param("agent_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid agent ID"})
		return
	}
	loanID, err := strconv.Atoi(c.Param("loan_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid loan ID"})
		return
	}
	var req models.AgentDecisionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	loan, err := h.agentService.MakeDecision(agentID, loanID, req.Decision)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Decision recorded successfully", "loan": loan})
}
