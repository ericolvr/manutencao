# Sistema de Gestão de Tickets de Manutenção

Sistema backend em Golang para gerenciamento de tickets de manutenção, desenvolvido com arquitetura limpa e suporte a múltiplos ambientes.

## Arquitetura Limpa (Clean Architecture)

O projeto segue os princípios da Clean Architecture com separação clara de responsabilidades:

```
Routes → Handlers → Services → Repositories → Database
```

### Fluxo de Dados
1. **Routes** (`internal/routes/`): Define endpoints HTTP e conecta às funções dos handlers
2. **Handlers** (`internal/handlers/`): Recebe requisições HTTP, valida dados de entrada e retorna respostas
3. **Services** (`internal/service/`): Contém a lógica de negócio e regras da aplicação
4. **Repositories** (`internal/repository/`): Abstrai o acesso aos dados e operações de banco
5. **Domain** (`internal/domain/`): Define as entidades e modelos de negócio

### Responsabilidades por Camada
- **Routes**: Mapeamento de URLs para handlers
- **Handlers**: Validação de entrada, serialização JSON, códigos HTTP
- **Services**: Regras de negócio, validações complexas, orquestração
- **Repositories**: Queries SQL, mapeamento de dados, transações
- **DTOs** (`internal/dto/`): Objetos de transferência entre camadas

## Arquitetura

- **Backend**: Go com Gin Framework
- **Banco de Dados**: PostgreSQL 16
- **Containerização**: Docker + Docker Compose
- **Deploy**: Docker Swarm
- **Configuração**: Múltiplos ambientes (.env.local, .env.deploy)

## Principais Rotas

Todas as rotas seguem o padrão REST e estão organizadas sob `/api/v1/`.

### Tickets
| Método | Endpoint | Descrição |
|--------|----------|----------|
| `POST` | `/api/v1/tickets` | Criar novo ticket |
| `GET` | `/api/v1/tickets` | Listar tickets |
| `GET` | `/api/v1/tickets/:id` | Buscar ticket por ID |
| `PUT` | `/api/v1/tickets/:id` | Atualizar ticket |
| `DELETE` | `/api/v1/tickets/:id` | Excluir ticket |
| `GET` | `/api/v1/tickets/number` | Obter número do próximo ticket |
| `POST` | `/api/v1/tickets/:id/providers` | Associar fornecedor ao ticket |
| `GET` | `/api/v1/tickets/:id/providers` | Listar fornecedores do ticket |
| `DELETE` | `/api/v1/tickets/:id/providers` | Remover fornecedor do ticket |
| `POST` | `/api/v1/tickets/:id/problems` | Associar problema ao ticket |
| `GET` | `/api/v1/tickets/:id/problems` | Listar problemas do ticket |
| `DELETE` | `/api/v1/tickets/:id/problems/:problem_id` | Remover problema do ticket |
| `POST` | `/api/v1/tickets/:id/solutions` | Associar solução ao ticket |
| `GET` | `/api/v1/tickets/:id/solutions` | Listar soluções do ticket |
| `DELETE` | `/api/v1/tickets/:id/solutions/:solution_id` | Remover solução do ticket |

### Users (Usuários)
| Método | Endpoint | Descrição |
|--------|----------|----------|
| `POST` | `/api/v1/users` | Criar novo usuário |
| `GET` | `/api/v1/users` | Listar todos os usuários |
| `GET` | `/api/v1/users/:id` | Buscar usuário por ID |
| `PUT` | `/api/v1/users/:id` | Atualizar usuário |
| `DELETE` | `/api/v1/users/:id` | Excluir usuário |
| `POST` | `/api/v1/users/auth` | **Autenticar usuário** |

### Providers (Técnicos/Fornecedores)
| Método | Endpoint | Descrição |
|--------|----------|----------|
| `POST` | `/api/v1/providers` | Criar novo técnico |
| `GET` | `/api/v1/providers` | Listar todos os técnicos |
| `GET` | `/api/v1/providers/:id` | Buscar técnico por ID |
| `GET` | `/api/v1/providers/name/:name` | Buscar técnico por nome |
| `PUT` | `/api/v1/providers/:id` | Atualizar técnico |
| `DELETE` | `/api/v1/providers/:id` | Excluir técnico |

### Branchs (Agências)
| Método | Endpoint | Descrição |
|--------|----------|----------|
| `POST` | `/api/v1/branchs` | Criar nova agência |
| `GET` | `/api/v1/branchs` | Listar todas as agências |
| `GET` | `/api/v1/branchs/:id` | Buscar agência por ID |
| `GET` | `/api/v1/branchs/client/:client` | Buscar agências por cliente |
| `PUT` | `/api/v1/branchs/:id` | Atualizar agência |
| `DELETE` | `/api/v1/branchs/:id` | Excluir agência |

### Clients (Clientes)
| Método | Endpoint | Descrição |
|--------|----------|----------|
| `POST` | `/api/v1/clients` | Criar novo cliente |
| `GET` | `/api/v1/clients` | Listar todos os clientes |
| `GET` | `/api/v1/clients/:id` | Buscar cliente por ID |
| `PUT` | `/api/v1/clients/:id` | Atualizar cliente |
| `DELETE` | `/api/v1/clients/:id` | Excluir cliente |

### Problems (Problemas)
| Método | Endpoint | Descrição |
|--------|----------|----------|
| `POST` | `/api/v1/problems` | Criar novo problema |
| `GET` | `/api/v1/problems` | Listar todos os problemas |
| `GET` | `/api/v1/problems/:id` | Buscar problema por ID |
| `PUT` | `/api/v1/problems/:id` | Atualizar problema |
| `DELETE` | `/api/v1/problems/:id` | Excluir problema |

### Solutions (Soluções)
| Método | Endpoint | Descrição |
|--------|----------|----------|
| `POST` | `/api/v1/solutions` | Criar nova solução |
| `GET` | `/api/v1/solutions` | Listar todas as soluções |
| `GET` | `/api/v1/solutions/:id` | Buscar solução por ID |
| `PUT` | `/api/v1/solutions/:id` | Atualizar solução |
| `DELETE` | `/api/v1/solutions/:id` | Excluir solução |

## Collection de Exemplo

Uma collection completa do Postman com exemplos de todas as rotas está disponível em:

**`postman_collection_v2.json`**

A collection inclui:
- Todos os endpoints documentados
- Exemplos de request/response
- Variáveis de ambiente configuradas
- Casos de uso práticos

### Como usar:
1. Importe o arquivo `postman_collection_v2.json` no Postman
2. Configure as variáveis de ambiente (`base_url`, etc.)
3. Execute as requisições de exemplo

## Início Rápido

### Pré-requisitos
- Go 1.24+
- Docker & Docker Compose
- Make

### Desenvolvimento Local

```bash
# Configuração completa (instala dependências + banco + servidor)
make dev

# Ou passo a passo:
make install    # Instala dependências Go
make database   # Inicia container PostgreSQL
make run        # Executa servidor de desenvolvimento
```

### Deploy em Produção

```bash
# Deploy completo
make deploy

# Ou passo a passo:
make swarm-init  # Inicializa Docker Swarm
make build       # Constrói imagem Docker
make stack-up    # Deploy do stack
```

## Comandos Disponíveis

### Desenvolvimento Local
| Comando | Descrição |
|---------|-----------|
| `make install` | Instala dependências Go |
| `make database` | Inicia container PostgreSQL |
| `make run` | Executa servidor de desenvolvimento |
| `make dev` | Setup completo (install + database + ready) |
| `make db-stop` | Para container do banco |
| `make db-clean` | Limpa dados do banco |

### Deploy e Produção
| Comando | Descrição |
|---------|-----------|
| `make deploy` | Deploy completo (swarm + build + stack) |
| `make swarm-init` | Inicializa Docker Swarm |
| `make build` | Constrói imagem Docker |
| `make stack-up` | Deploy do stack no swarm |
| `make stack-down` | Remove stack do swarm |

### Monitoramento
| Comando | Descrição |
|---------|-----------|
| `make status` | Mostra status dos serviços |
| `make logs` | Exibe logs em tempo real |
| `make ps` | Lista tasks do stack |
| `make scale REPLICAS=N` | Escala serviço da API |

### Gerenciamento de Ambientes
| Comando | Descrição |
|---------|-----------|
| `make env-status` | Mostra configuração atual |
| `make env-create` | Cria arquivos de ambiente |
| `make env-switch` | Troca para ambiente específico |

## Controle de Ambientes

O sistema suporta múltiplos ambientes com switching automático:

### Uso Automático
```bash
make dev        # Usa automaticamente .env.local
make deploy     # Usa automaticamente .env.deploy
```

### Uso Manual
```bash
ENV=local make run      # Força ambiente local
ENV=deploy make deploy  # Força ambiente de produção
```

### Arquivos de Ambiente
- **`.env-sample`**: Template com documentação
- **`.env.local`**: Configuração para desenvolvimento local
- **`.env.deploy`**: Configuração para produção
- **`.env`**: Arquivo ativo (copiado automaticamente)

## Configuração do Banco de Dados

### Local (desenvolvimento)
```
postgresql://maintenance:maintenance@127.0.0.1:5432/maintenance
```

### Deploy (produção)
```
postgresql://maintenance:maintenance@postgres_database:5432/maintenance
```

## Usuários Padrão

O sistema inclui um usuário administrador padrão:
- **Usuário**: `admin`
- **Senha**: `admin123`
- **Papel**: Administrador (0)

### Níveis de Acesso
- **0**: Admin
- **1**: Financeiro
- **2**: Suporte
- **3**: Estoque
- **4**: Técnicos
- **5**: Pagamentos

## API Endpoints

O servidor roda por padrão na porta **9999**:
- **Local**: http://localhost:9999
- **Deploy**: http://localhost:9999 (via Docker Swarm)

## Estrutura do Projeto

```
maintenance-v2/
├── cmd/main.go              # Ponto de entrada
├── internal/
│   ├── domain/              # Modelos de domínio
│   ├── dto/                 # Data Transfer Objects
│   ├── handlers/            # Controladores HTTP
│   ├── middleware/          # Middlewares HTTP
│   ├── repository/          # Camada de dados
│   ├── routes/              # Definição de rotas
│   └── service/             # Lógica de negócio
├── config/                  # Configurações
├── scripts/                 # Scripts SQL
├── uploads/                 # Arquivos enviados
├── docker-compose.yml       # Orquestração Docker
├── Dockerfile              # Imagem da aplicação
├── Makefile                # Automação de comandos
├── postman_collection_v2.json # Collection Postman
└── .env.*                  # Arquivos de ambiente
```

## Exemplos de Uso

### Desenvolvimento
```bash
# Iniciar desenvolvimento
make dev

# Verificar status do ambiente
make env-status

# Trocar para ambiente local
ENV=local make env-switch
```

### Produção
```bash
# Deploy completo
make deploy

# Escalar para 5 réplicas
make scale REPLICAS=5

# Verificar logs
make logs
```

## Associações de Tickets

O sistema permite associar diferentes entidades aos tickets:

### Associar Fornecedor a um Ticket
```bash
POST /api/v1/tickets/1/providers
{
  "provider_id": 1
}
```

### Associar Problema a um Ticket
```bash
POST /api/v1/tickets/1/problems
{
  "problem_id": 1,
  "description": "Descrição do problema específico"
}
```

### Associar Solução a um Ticket
```bash
POST /api/v1/tickets/1/solutions
{
  "solution_id": 1,
  "cost": 150.00,
  "description": "Descrição da solução aplicada"
}
```
