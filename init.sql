-- Tabela de clientes
CREATE TABLE IF NOT EXISTS clients (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    mobile VARCHAR(20),
    zipcode VARCHAR(10),
    state VARCHAR(50),
    city VARCHAR(100),
    neighborhood VARCHAR(100),
    address VARCHAR(255),
    complement VARCHAR(255)
);

-- Tabela de filiais
CREATE TABLE IF NOT EXISTS branches (
    id SERIAL PRIMARY KEY,
    client INTEGER REFERENCES clients(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    uniorg VARCHAR(50),
    zipcode VARCHAR(10),
    state VARCHAR(50),
    city VARCHAR(100),
    neighborhood VARCHAR(100),
    address VARCHAR(255),
    complement VARCHAR(255)
);

-- Tabela de prestadores
CREATE TABLE IF NOT EXISTS providers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    mobile VARCHAR(20),
    zipcode VARCHAR(10),
    state VARCHAR(50),
    city VARCHAR(100),
    neighborhood VARCHAR(100),
    address VARCHAR(255),
    complement VARCHAR(255)
);

-- Tabela de usuários
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

-- Tabela de tickets SIMPLIFICADA (sem categoria direta, usa associação com problemas)
CREATE TABLE IF NOT EXISTS tickets (
    id SERIAL PRIMARY KEY,
    number VARCHAR(50) NOT NULL UNIQUE,
    status INTEGER NOT NULL DEFAULT 1,
    priority VARCHAR(50) NOT NULL,
    description TEXT NOT NULL,
    open_date TIMESTAMP NOT NULL,
    close_date TIMESTAMP NULL,
    branch_id INTEGER NOT NULL REFERENCES branches(id) ON DELETE CASCADE,
    provider_id INTEGER NULL REFERENCES providers(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabela de problemas específicos
CREATE TABLE IF NOT EXISTS problems (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabela de catálogo de soluções (referencia problemas específicos)
CREATE TABLE IF NOT EXISTS solutions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    unit_price DECIMAL(10,2) NOT NULL,
    problem_id INTEGER REFERENCES problems(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabela de associação many-to-many entre tickets e problemas
CREATE TABLE IF NOT EXISTS ticket_problems (
    id SERIAL PRIMARY KEY,
    ticket_id INTEGER NOT NULL REFERENCES tickets(id) ON DELETE CASCADE,
    problem_id INTEGER NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(ticket_id, problem_id)
);

-- Tabela de custos dos tickets (histórico de soluções aplicadas)
CREATE TABLE IF NOT EXISTS ticket_costs (
    id SERIAL PRIMARY KEY,
    ticket_id INTEGER NOT NULL REFERENCES tickets(id) ON DELETE CASCADE,
    solution_id INTEGER NOT NULL REFERENCES solutions(id),
    solution_name VARCHAR(255) NOT NULL DEFAULT '',
    quantity INTEGER NOT NULL DEFAULT 1,
    unit_price DECIMAL(10,2) NOT NULL,  -- Preço no momento da aplicação
    subtotal DECIMAL(10,2) NOT NULL,    -- quantity * unit_price
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Índices para performance
CREATE INDEX IF NOT EXISTS idx_tickets_status ON tickets(status);
CREATE INDEX IF NOT EXISTS idx_tickets_branch_id ON tickets(branch_id);
CREATE INDEX IF NOT EXISTS idx_tickets_provider_id ON tickets(provider_id);
CREATE INDEX IF NOT EXISTS idx_tickets_open_date ON tickets(open_date);
CREATE INDEX IF NOT EXISTS idx_tickets_number ON tickets(number);
CREATE INDEX IF NOT EXISTS idx_solutions_problem_id ON solutions(problem_id);
CREATE INDEX IF NOT EXISTS idx_ticket_problems_ticket_id ON ticket_problems(ticket_id);
CREATE INDEX IF NOT EXISTS idx_ticket_problems_problem_id ON ticket_problems(problem_id);
CREATE INDEX IF NOT EXISTS idx_ticket_costs_ticket_id ON ticket_costs(ticket_id);
CREATE INDEX IF NOT EXISTS idx_ticket_costs_solution_id ON ticket_costs(solution_id);

-- Dados de exemplo para clientes
INSERT INTO clients (name, mobile, zipcode, state, city, neighborhood, address, complement) VALUES
('Empresa ABC Ltda', '11999887766', '01310-100', 'SP', 'São Paulo', 'Bela Vista', 'Av. Paulista, 1000', 'Conjunto 101'),
('Comércio XYZ S/A', '21988776655', '20040-020', 'RJ', 'Rio de Janeiro', 'Centro', 'Rua da Carioca, 50', '2º andar'),
('Indústria DEF ME', '31977665544', '30112-000', 'MG', 'Belo Horizonte', 'Centro', 'Rua dos Carijós, 200', 'Sala 301');

-- Dados de exemplo para filiais
INSERT INTO branches (client, name, uniorg, zipcode, state, city, neighborhood, address, complement) VALUES
(1, 'Filial São Paulo Centro', 'SP001', '01310-100', 'SP', 'São Paulo', 'Bela Vista', 'Av. Paulista, 1000', 'Conjunto 101'),
(1, 'Filial São Paulo Norte', 'SP002', '02071-000', 'SP', 'São Paulo', 'Santana', 'Av. Cruzeiro do Sul, 1100', 'Bloco A'),
(2, 'Filial Rio Centro', 'RJ001', '20040-020', 'RJ', 'Rio de Janeiro', 'Centro', 'Rua da Carioca, 50', '2º andar'),
(3, 'Filial BH Principal', 'MG001', '30112-000', 'MG', 'Belo Horizonte', 'Centro', 'Rua dos Carijós, 200', 'Sala 301');

-- Dados de exemplo para prestadores
INSERT INTO providers (name, mobile, zipcode, state, city, neighborhood, address, complement) VALUES
('TechService Manutenção', '11987654321', '04567-890', 'SP', 'São Paulo', 'Vila Olímpia', 'Rua Funchal, 500', 'Conjunto 1001'),
('RapidFix Assistência', '21976543210', '22071-900', 'RJ', 'Rio de Janeiro', 'Copacabana', 'Av. Atlântica, 1702', 'Cobertura'),
('MegaTech Soluções', '31965432109', '30140-071', 'MG', 'Belo Horizonte', 'Funcionários', 'Av. do Contorno, 6061', '12º andar');

-- Dados de exemplo para usuários
INSERT INTO users (name, mobile, password, role, status) VALUES
('Administrador', '11999000001', '$2a$10$N9qo8uLOickgx2ZMRZoMye', 1, true),
('João Silva', '11999000002', '$2a$10$N9qo8uLOickgx2ZMRZoMye', 3, true),
('Maria Santos', '11999000003', '$2a$10$N9qo8uLOickgx2ZMRZoMye', 5, true);

-- Dados de exemplo para problemas (sem categorias)
INSERT INTO problems (name, description) VALUES
('Ar condicionado não gela', 'Sistema de refrigeração não está funcionando adequadamente'),
('Ar condicionado faz barulho excessivo', 'Ruídos anômalos durante funcionamento'),
('Ar condicionado vaza água', 'Vazamento de água do equipamento'),
('Ar condicionado não liga', 'Equipamento não responde ao comando'),
('Lâmpada queimada', 'Iluminação não funciona por lâmpada defeituosa'),
('Tomada sem energia', 'Ponto de energia não fornece eletricidade'),
('Disjuntor desarma frequentemente', 'Proteção elétrica atua constantemente'),
('Internet lenta', 'Conexão com velocidade abaixo do esperado'),
('Computador não liga', 'Equipamento não inicializa'),
('Impressora não funciona', 'Equipamento de impressão com defeito'),
('Torneira pingando', 'Vazamento constante na torneira'),
('Vaso sanitário entupido', 'Obstrução no sistema de esgoto');

-- Dados de exemplo para soluções do catálogo
INSERT INTO solutions (name, description, unit_price, problem_id) VALUES
-- Soluções para Ar Condicionado
('Limpeza de filtros', 'Limpeza e higienização dos filtros do ar condicionado', 80.00, 1),
('Recarga de gás R22', 'Recarga do gás refrigerante R22', 150.00, 1),
('Recarga de gás R410A', 'Recarga do gás refrigerante R410A', 180.00, 1),
('Mão de obra especializada', 'Serviço técnico especializado em climatização', 140.00, 1),
('Lubrificação de componentes', 'Lubrificação de partes móveis do equipamento', 60.00, 2),
('Limpeza do dreno', 'Desobstrução e limpeza do sistema de drenagem', 90.00, 3),
('Verificação elétrica', 'Diagnóstico e reparo de componentes elétricos', 120.00, 4),
-- Soluções para Elétrica
('Troca de lâmpada LED', 'Substituição por lâmpada LED equivalente', 25.00, 5),
('Troca de lâmpada fluorescente', 'Substituição por lâmpada fluorescente', 15.00, 5),
('Verificação de fiação', 'Inspeção e reparo da instalação elétrica', 100.00, 6),
('Troca de disjuntor', 'Substituição de disjuntor defeituoso', 80.00, 7),
-- Soluções para Rede/TI
('Configuração de rede', 'Otimização e configuração da rede', 120.00, 8),
('Diagnóstico de hardware', 'Verificação de componentes do computador', 100.00, 9),
('Manutenção de impressora', 'Limpeza e calibração da impressora', 80.00, 10),
-- Soluções para Hidráulica
('Troca de reparo da torneira', 'Substituição de componentes internos', 45.00, 11),
('Desentupimento', 'Desobstrução do sistema de esgoto', 120.00, 12);

-- Dados de exemplo para tickets
INSERT INTO tickets (number, status, priority, description, open_date, branch_id, provider_id) VALUES
('TK-2024-001', 1, 'Alta', 'Ar condicionado da sala de reuniões não está gelando adequadamente', '2024-01-15 09:00:00', 1, 1),
('TK-2024-002', 2, 'Média', 'Lâmpadas do corredor principal queimadas', '2024-01-16 14:30:00', 2, 2),
('TK-2024-003', 1, 'Baixa', 'Internet lenta no setor administrativo', '2024-01-17 11:15:00', 3, 3);

-- Associações de problemas aos tickets
INSERT INTO ticket_problems (ticket_id, problem_id) VALUES
(1, 1), -- TK-2024-001 → Ar condicionado não gela
(2, 5), -- TK-2024-002 → Lâmpada queimada
(3, 8); -- TK-2024-003 → Internet lenta

-- Custos aplicados aos tickets
INSERT INTO ticket_costs (ticket_id, solution_id, solution_name, quantity, unit_price, subtotal) VALUES
-- TK-2024-001: Ar condicionado não gela
(1, 1, 'Limpeza de filtros', 2, 80.00, 160.00),   -- Limpeza de filtros x2
(1, 2, 'Recarga de gás R22', 1, 150.00, 150.00),  -- Recarga de gás R22 x1
(1, 4, 'Mão de obra especializada', 1, 140.00, 140.00),  -- Mão de obra especializada x1
-- TK-2024-002: Lâmpadas queimadas
(2, 9, 'Troca de lâmpada LED', 3, 25.00, 75.00),    -- Troca de lâmpada LED x3
-- TK-2024-003: Internet lenta
(3, 13, 'Configuração de rede', 1, 120.00, 120.00); -- Configuração de rede x1




