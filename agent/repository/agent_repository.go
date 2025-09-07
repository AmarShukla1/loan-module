package repository

import (
	"loan-module/agent/models"
	"loan-module/repository"
)

type AgentRepository struct {
	DB *database.Database
}

func NewAgentRepository(db *database.Database) *AgentRepository {
	return &AgentRepository{DB: db}
}

func (r *AgentRepository) AddAgent(agent *models.Agent) (*models.Agent, error) {
	result := r.DB.DB.Create(agent)
	if result.Error != nil {
		return nil, result.Error
	}
	return agent, nil
}

func (r *AgentRepository) GetAgentByID(id int) (*models.Agent, bool) {
	var agent models.Agent
	result := r.DB.DB.First(&agent, id)
	return &agent, result.Error == nil
}

func (r *AgentRepository) GetAvailableAgent() *models.Agent {
	type AgentLoad struct {
		ID    int
		Count int
	}

	var loads []AgentLoad
	r.DB.DB.Raw(`
        SELECT a.id, COUNT(l.id) as count
        FROM agents a
        LEFT JOIN loans l ON l.assigned_agent_id = a.id 
                           AND l.application_status IN ('PROCESSING', 'UNDER_REVIEW')
        WHERE a.manager_id IS NOT NULL
        GROUP BY a.id
        ORDER BY count ASC, a.id ASC
        LIMIT 1
    `).Scan(&loads)

	if len(loads) == 0 {
		return nil
	}

	var agent models.Agent
	r.DB.DB.First(&agent, loads[0].ID)
	return &agent
}
