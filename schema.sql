-- Loan Origination System (LOS) Database Schema
-- This file contains the complete database schema with foreign key relationships

-- =============================================
-- CUSTOMERS TABLE
-- =============================================
CREATE TABLE customers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(20) NOT NULL,
    email VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create unique index on phone for customer deduplication
CREATE UNIQUE INDEX idx_customers_phone ON customers(phone);

-- =============================================
-- AGENTS TABLE
-- =============================================
CREATE TABLE agents (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    manager_id INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Foreign key constraint: Self-referencing for manager hierarchy
    CONSTRAINT fk_agents_manager 
        FOREIGN KEY (manager_id) 
        REFERENCES agents(id) 
        ON DELETE SET NULL
);

-- Create index on manager_id for hierarchy queries
CREATE INDEX idx_agents_manager_id ON agents(manager_id);

-- =============================================
-- LOANS TABLE
-- =============================================
CREATE TABLE loans (
    id SERIAL PRIMARY KEY,
    customer_id INTEGER NOT NULL,
    loan_amount DECIMAL(15,2) NOT NULL,
    loan_type VARCHAR(20) NOT NULL CHECK (loan_type IN ('PERSONAL', 'HOME', 'AUTO', 'BUSINESS')),
    application_status VARCHAR(30) NOT NULL CHECK (application_status IN (
        'APPLIED', 'PROCESSING', 'APPROVED_BY_SYSTEM', 'REJECTED_BY_SYSTEM', 
        'UNDER_REVIEW', 'APPROVED_BY_AGENT', 'REJECTED_BY_AGENT'
    )),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    assigned_agent_id INTEGER,
    
    -- Foreign key constraints
    CONSTRAINT fk_loans_customer 
        FOREIGN KEY (customer_id) 
        REFERENCES customers(id) 
        ON DELETE CASCADE,
    
    CONSTRAINT fk_loans_assigned_agent 
        FOREIGN KEY (assigned_agent_id) 
        REFERENCES agents(id) 
        ON DELETE SET NULL
);

-- Create indexes for foreign key relationships
CREATE INDEX idx_loans_customer_id ON loans(customer_id);
CREATE INDEX idx_loans_assigned_agent_id ON loans(assigned_agent_id);

-- =============================================
-- LOAN ASSIGNMENTS TABLE
-- =============================================
CREATE TABLE loan_assignments (
    id SERIAL PRIMARY KEY,
    loan_id INTEGER NOT NULL,
    agent_id INTEGER NOT NULL,
    assigned_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Foreign key constraints
    CONSTRAINT fk_loan_assignments_loan 
        FOREIGN KEY (loan_id) 
        REFERENCES loans(id) 
        ON DELETE CASCADE,
    
    CONSTRAINT fk_loan_assignments_agent 
        FOREIGN KEY (agent_id) 
        REFERENCES agents(id) 
        ON DELETE CASCADE
);

-- Create indexes for foreign key relationships
CREATE INDEX idx_loan_assignments_loan_id ON loan_assignments(loan_id);
CREATE INDEX idx_loan_assignments_agent_id ON loan_assignments(agent_id);
