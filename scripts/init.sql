-- Database: maintenance_v2
-- Complete schema with all entities

-- Branchs table  
CREATE TABLE branchs (
    id SERIAL PRIMARY KEY,
    client VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    uniorg VARCHAR(50),
    zipcode VARCHAR(20),
    state VARCHAR(100),
    city VARCHAR(100),
    neighborhood VARCHAR(100),
    address TEXT,
    complement VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Clients table
CREATE TABLE clients (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


-- Providers table
CREATE TABLE providers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    mobile VARCHAR(20),
    zipcode VARCHAR(20),
    state VARCHAR(100),
    city VARCHAR(100),
    neighborhood VARCHAR(100),
    address TEXT,
    complement VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Problems table
CREATE TABLE problems (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Solutions table
CREATE TABLE solutions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    unit_price DECIMAL(10,2) DEFAULT 0.00,
    problem_id INTEGER REFERENCES problems(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tickets table
CREATE TABLE tickets (
    id SERIAL PRIMARY KEY,
    number VARCHAR(50) UNIQUE NOT NULL,
    status INTEGER DEFAULT 12, -- 12 = Novo
    priority VARCHAR(20) DEFAULT 'medium',
    description TEXT,
    open_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    close_date TIMESTAMP NULL,
    branch_id INTEGER REFERENCES branchs(id) ON DELETE CASCADE,
    provider_id INTEGER REFERENCES providers(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Ticket Costs table
CREATE TABLE ticket_costs (
    id SERIAL PRIMARY KEY,
    ticket_id INTEGER REFERENCES tickets(id) ON DELETE CASCADE,
    solution_id INTEGER REFERENCES solutions(id) ON DELETE SET NULL,
    quantity INTEGER DEFAULT 1,
    unit_price DECIMAL(10,2) DEFAULT 0.00,
    subtotal DECIMAL(10,2) DEFAULT 0.00,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Ticket Problems table (many-to-many relationship)
CREATE TABLE ticket_problems (
    id SERIAL PRIMARY KEY,
    ticket_id INTEGER REFERENCES tickets(id) ON DELETE CASCADE,
    problem_id INTEGER REFERENCES problems(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(ticket_id, problem_id)
);

-- Indexes for better performance
CREATE INDEX idx_branchs_client ON branchs(client);
CREATE INDEX idx_providers_name ON providers(name);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_solutions_problem ON solutions(problem_id);
CREATE INDEX idx_tickets_branch ON tickets(branch_id);
CREATE INDEX idx_tickets_provider ON tickets(provider_id);
CREATE INDEX idx_tickets_status ON tickets(status);
CREATE INDEX idx_tickets_number ON tickets(number);
CREATE INDEX idx_ticket_costs_ticket ON ticket_costs(ticket_id);
CREATE INDEX idx_ticket_costs_solution ON ticket_costs(solution_id);
CREATE INDEX idx_ticket_problems_ticket ON ticket_problems(ticket_id);
CREATE INDEX idx_ticket_problems_problem ON ticket_problems(problem_id);

