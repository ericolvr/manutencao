-- Database: maintenance_v2
-- Complete schema with all entities

-- Branchs table  
CREATE TABLE branchs (
    id SERIAL PRIMARY KEY,
    client VARCHAR(255),
    name VARCHAR(255) NOT NULL,
    uniorg VARCHAR(50),
    zipcode VARCHAR(10),
    state VARCHAR(50),
    city VARCHAR(100),
    neighborhood VARCHAR(100),
    address VARCHAR(255),
    complement VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert Branchs
INSERT INTO branchs (client, name, uniorg, zipcode, state, city, neighborhood, address, complement) VALUES 
('Basa', 'Icoaraci', '001-0001', '68810-100', 'PA', 'Belém', 'Centro', 'Rua Manoel Barata, 660', ''),
('Basa', 'Imperatriz', '002-0002', '65900-120', 'MA', 'Imperatriz', 'Beira io', 'Av. getúlio Vargas, 404', ''),
('Correios', 'Vanderlei', '003-0009', '05011-001', 'SP', 'São Paulo', 'Pompéia', 'Rua Vanderlei, 832', '');


-- Clients table
CREATE TABLE clients (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert Clients
INSERT INTO clients (name) VALUES
('Basa'),
('Banco do Nordeste'),
('Correios'),
('Santander');


-- Costs table  
CREATE TABLE costs (
    id SERIAL PRIMARY KEY,
    value_per_km DECIMAL(10,2) NOT NULL,
    initial_value DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- Insert Costs
INSERT INTO costs (value_per_km, initial_value) VALUES
(1.00, 100.00);


-- Distances table
CREATE TABLE IF NOT EXISTS distances (
    id SERIAL PRIMARY KEY,
    distance DECIMAL(10,2) NOT NULL,
    ticket_number VARCHAR(50) NOT NULL,
    provider_id INTEGER NOT NULL,
    provider_name VARCHAR(255) NOT NULL,
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

-- Insert Providers
INSERT INTO providers (name, mobile, zipcode, state, city, neighborhood, address, complement) VALUES
('Julio Mesquita', '11987654321', '04531-080', 'SP', 'São Paulo', 'Itaim Bibi', 'Rua Prof. Carlos de Carvalho, 74', ''),
('Manoel Santos', '11987654321', '02553-050', 'SP', 'São Paulo', 'Casa Verde', 'Rua António João 600', ''),
('Mario Lago', '31965432109', '05011-001', 'SP', 'São Paulo', 'Pompéia', 'Rua Vanderlei, 832', '');


-- Users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    mobile VARCHAR(20) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role INTEGER NOT NULL DEFAULT 0,
    status BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert Users
INSERT INTO users (name, mobile, password, role, status) VALUES
('Administrador', '11999000001', '$2a$10$N9qo8uLOickgx2ZMRZoMye', 1, true),
('João Silva', '11999000002', '$2a$10$N9qo8uLOickgx2ZMRZoMye', 3, true),
('Maria Santos', '11999000003', '$2a$10$N9qo8uLOickgx2ZMRZoMye', 5, true);


-- Ticket table
CREATE TABLE IF NOT EXISTS tickets (
    id SERIAL PRIMARY KEY,
    number VARCHAR(50) NOT NULL UNIQUE,
    status INTEGER NOT NULL DEFAULT 1,
    priority VARCHAR(50) NOT NULL,
    description TEXT NOT NULL,
    open_date TIMESTAMP NOT NULL,
    close_date TIMESTAMP NULL,
    branch_id INTEGER NOT NULL,
    provider_id INTEGER NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Problem table
CREATE TABLE IF NOT EXISTS problems (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- Insert Problems
INSERT INTO problems (name, description) VALUES
('Câmera não exibe imagem', 'Câmera não exibe imagem'),
('Falha de Bateria', 'Gerador com Falha de Bateria'),
('Sirene', 'Problemas na sirene');


-- Solution table
CREATE TABLE IF NOT EXISTS solutions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    unit_price DECIMAL(10,2) NOT NULL,
    problem_id INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- Insert Solutions
INSERT INTO solutions (name, description, unit_price, problem_id) VALUES
('Configuração de Câmera', 'Configuração de Câmera', 80.00, 1),
('Troca de Câmera', 'Troca de Câmera', 2100.00, 1),
('Configuração de DVR', 'Configuração de DVR', 350.00, 1),
('Verificação de Bateria', 'Verificação de Bateria', 100.00, 2),
('Troca de Bateria', 'Troca de Bateria', 3890.00, 2),
('Configuração de Sirene', 'Configuração de Sirene', 120.00, 3),
('Troca de Sirene', 'Troca de Sirene', 876.00, 3);

-- TicketProblem table
CREATE TABLE IF NOT EXISTS ticket_problems (
    id SERIAL PRIMARY KEY,
    ticket_id INTEGER NOT NULL,
    problem_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(ticket_id, problem_id)
);

-- TicketCost table
CREATE TABLE IF NOT EXISTS ticket_costs (
    id SERIAL PRIMARY KEY,
    ticket_id INTEGER NOT NULL,
    problem_id INTEGER NOT NULL,
    problem_name VARCHAR(255) NOT NULL DEFAULT '',
    solution_id INTEGER NOT NULL,
    solution_name VARCHAR(255) NOT NULL DEFAULT '',
    quantity INTEGER NOT NULL DEFAULT 1,
    unit_price DECIMAL(10,2) NOT NULL,  -- Preço no momento da aplicação
    subtotal DECIMAL(10,2) NOT NULL,    -- quantity * unit_price
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_tickets_status ON tickets(status);
CREATE INDEX IF NOT EXISTS idx_tickets_branch_id ON tickets(branch_id);
CREATE INDEX IF NOT EXISTS idx_tickets_provider_id ON tickets(provider_id);
CREATE INDEX IF NOT EXISTS idx_tickets_open_date ON tickets(open_date);
CREATE INDEX IF NOT EXISTS idx_tickets_number ON tickets(number);
CREATE INDEX IF NOT EXISTS idx_solutions_problem_id ON solutions(problem_id);
CREATE INDEX IF NOT EXISTS idx_ticket_problems_ticket_id ON ticket_problems(ticket_id);
CREATE INDEX IF NOT EXISTS idx_ticket_problems_problem_id ON ticket_problems(problem_id);
CREATE INDEX IF NOT EXISTS idx_ticket_costs_ticket_id ON ticket_costs(ticket_id);
CREATE INDEX IF NOT EXISTS idx_ticket_costs_problem_id ON ticket_costs(problem_id);
CREATE INDEX IF NOT EXISTS idx_ticket_costs_solution_id ON ticket_costs(solution_id);
CREATE INDEX IF NOT EXISTS idx_distances_ticket_number ON distances(ticket_number);
CREATE INDEX IF NOT EXISTS idx_distances_provider_id ON distances(provider_id);
