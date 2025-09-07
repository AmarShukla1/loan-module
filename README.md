# Loan Module Application

A loan origination system (LOS) built with Go that handles loan applications, customer management, and agent decision processes.

## Features

- Customer management
- Loan application submission and processing
- Automatic loan approval/rejection based on amount thresholds
- Agent review and decision making for loans
- Notification service
- RESTful API endpoints

## Tech Stack

- **Backend**: Go with Gin framework
- **Database**: PostgreSQL
- **ORM**: GORM

## Prerequisites

- Go 1.23.2 or higher
- PostgreSQL 12 or higher

## Database Setup

### PostgreSQL Setup

1. Install PostgreSQL if you haven't already:

   ```bash
   # macOS (using Homebrew)
   brew install postgresql
   
   # Start PostgreSQL service
   brew services start postgresql
   ```

2. Create a database for the application:

   ```bash
   # Login to PostgreSQL
   psql -U postgres
   
   # Create database
   CREATE DATABASE loandb;
   
   # Connect to the database
   \c loandb
   ```

3. Initialize the database schema:

   ```bash
   # From the project root directory
   psql -U postgres -d loandb -f schema.sql
   ```

## Configuration

The application uses a YAML configuration file (`loan-module-configuration.yaml`) for database connection settings. Update this file with your PostgreSQL credentials:

```yaml
db:
  timeZone: "UTC"
  host: "localhost"
  user: "your_username"
  password: "your_password"
  port: 5432
  name: "loandb"
  timeout: 5
  maxIdleConn: 2
  maxOpenConn: 4
```

## Running the Application

1. Clone the repository:

   ```bash
   git clone <repository-url>
   cd loan-module
   ```

2. Install dependencies:

   ```bash
   go mod download
   ```

3. Build and run the application:

   ```bash
   go build -o loan-app
   ./loan-app
   ```

   Alternatively, you can run directly with:

   ```bash
   go run main.go
   ```

4. The server will start on port 8080. You can access the API at `http://localhost:8080/api/v1/`

## API Endpoints

### Customer Endpoints

- `POST /api/v1/customers` - Create a new customer
- `GET /api/v1/customers/:id` - Get customer by ID
- `GET /api/v1/customers` - Get all customers
- `GET /api/v1/customers/top` - Get top customers with approved loans

### Loan Endpoints

- `POST /api/v1/loans` - Submit a new loan application
- `GET /api/v1/loans/status-count` - Get count of loans by status
- `GET /api/v1/loans` - Get loans by status
- `GET /api/v1/loans/:id` - Get loan by ID

### Agent Endpoints

- `PUT /api/v1/agents/:agent_id/loans/:loan_id/decision` - Make a decision on a loan
