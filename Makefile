# ANSI color codes
COLOR_RESET=\033[0m
COLOR_BOLD=\033[1m
COLOR_GREEN=\033[32m
COLOR_YELLOW=\033[33m
COLOR_BLUE=\033[34m

# Stack name
STACK_NAME=maintenance

# Environment configuration
# Default to local environment, can be overridden with ENV=deploy
ENV ?= local
ENV_FILE = .env.$(ENV)

help:
	@echo ""
	@echo "  $(COLOR_YELLOW)Available targets:$(COLOR_RESET)"
	@echo ""
	@echo "  $(COLOR_BLUE)Local Development:$(COLOR_RESET)"
	@echo "  $(COLOR_GREEN)install$(COLOR_RESET)		- Install Go dependencies"
	@echo "  $(COLOR_GREEN)database$(COLOR_RESET)		- Start Postgres container"
	@echo "  $(COLOR_GREEN)run$(COLOR_RESET)			- Run development server"
	@echo "  $(COLOR_GREEN)dev$(COLOR_RESET)			- Full local setup (install + database + run)"
	@echo "  $(COLOR_GREEN)db-stop$(COLOR_RESET)		- Stop database container"
	@echo "  $(COLOR_GREEN)db-clean$(COLOR_RESET)		- Clean database data"
	@echo ""
	@echo "  $(COLOR_BLUE)Deployment:$(COLOR_RESET)"
	@echo "  $(COLOR_GREEN)deploy$(COLOR_RESET)		- Full deployment (swarm + build + database)"
	@echo "  $(COLOR_GREEN)swarm-init$(COLOR_RESET)		- Initialize Docker Swarm"
	@echo "  $(COLOR_GREEN)build$(COLOR_RESET)			- Build Docker image"
	@echo "  $(COLOR_GREEN)stack-up$(COLOR_RESET)		- Deploy stack to swarm"
	@echo "  $(COLOR_GREEN)stack-down$(COLOR_RESET)		- Remove stack from swarm"
	@echo ""
	@echo "  $(COLOR_BLUE)Monitoring:$(COLOR_RESET)"
	@echo "  $(COLOR_GREEN)status$(COLOR_RESET)		- Show services status"
	@echo "  $(COLOR_GREEN)logs$(COLOR_RESET)			- Show logs in real-time"
	@echo "  $(COLOR_GREEN)ps$(COLOR_RESET)			- Show stack tasks"
	@echo "  $(COLOR_GREEN)scale$(COLOR_RESET)			- Scale service (use REPLICAS=N)"
	@echo ""
	@echo "  $(COLOR_YELLOW)Environment Control:$(COLOR_RESET)"
	@echo "  ENV=local make run		- Run with local config (.env.local)"
	@echo "  ENV=deploy make deploy	- Deploy with production config (.env.deploy)"
	@echo ""
	@echo "  $(COLOR_YELLOW)Examples:$(COLOR_RESET)"
	@echo "  make dev			- Start local development (uses .env.local)"
	@echo "  make deploy			- Deploy to production (uses .env.deploy)"
	@echo "  make scale REPLICAS=3		- Scale API service"
	@echo ""

# LOCAL DEVELOPMENT TARGETS
install:
	@echo "$(COLOR_YELLOW)Installing Go dependencies...$(COLOR_RESET)"
	go mod download
	go mod tidy
	@echo "$(COLOR_GREEN)✅ Dependencies installed successfully!$(COLOR_RESET)"

database:
	@echo "$(COLOR_YELLOW)Starting Postgres container for $(ENV) environment...$(COLOR_RESET)"
	@if [ ! -f $(ENV_FILE) ]; then echo "$(COLOR_YELLOW)⚠️  Environment file $(ENV_FILE) not found, creating from template...$(COLOR_RESET)"; cp .env.sample $(ENV_FILE); fi
	docker compose --env-file $(ENV_FILE) up -d postgres_database
	@echo "$(COLOR_GREEN)✅ Database container started!$(COLOR_RESET)"
	@echo "$(COLOR_BLUE)Using environment: $(ENV) ($(ENV_FILE))$(COLOR_RESET)"

run:
	@echo "$(COLOR_YELLOW)Starting development server with $(ENV) environment...$(COLOR_RESET)"
	@if [ ! -f $(ENV_FILE) ]; then echo "$(COLOR_YELLOW)⚠️  Environment file $(ENV_FILE) not found, creating from template...$(COLOR_RESET)"; cp .env.sample $(ENV_FILE); fi
	@echo "$(COLOR_BLUE)Loading environment from: $(ENV_FILE)$(COLOR_RESET)"
	@cp $(ENV_FILE) .env
	go run ./cmd/main.go

dev: install database
	@echo "$(COLOR_YELLOW)Waiting for database to be ready...$(COLOR_RESET)"
	@sleep 3
	@echo "$(COLOR_GREEN)✅ Local development environment ready!$(COLOR_RESET)"
	@echo "$(COLOR_BLUE)Environment: $(ENV) ($(ENV_FILE))$(COLOR_RESET)"
	@echo "$(COLOR_BLUE)Run 'make run' to start the server$(COLOR_RESET)"

db-stop:
	@echo "$(COLOR_YELLOW)Stopping database container...$(COLOR_RESET)"
	docker compose stop postgres_database
	@echo "$(COLOR_GREEN)✅ Database stopped!$(COLOR_RESET)"

db-clean:
	@echo "$(COLOR_YELLOW)Cleaning database data...$(COLOR_RESET)"
	docker compose down postgres_database
	docker volume rm maintenance-v2_postgres 2>/dev/null || true
	@echo "$(COLOR_GREEN)✅ Database data cleaned!$(COLOR_RESET)"


# DEPLOYMENT TARGETS
# Force deploy environment for deployment commands
deploy: ENV=deploy
deploy: swarm-init build stack-up
	@echo "$(COLOR_GREEN)✅ Deployment completed successfully!$(COLOR_RESET)"
	@echo "$(COLOR_BLUE)Environment: $(ENV) ($(ENV_FILE))$(COLOR_RESET)"
	@echo "$(COLOR_BLUE)API available at: http://localhost:9999$(COLOR_RESET)"

swarm-init:
	@echo "$(COLOR_YELLOW)Initializing Docker Swarm...$(COLOR_RESET)"
	@docker swarm init 2>/dev/null || echo "$(COLOR_GREEN)✅ Swarm already initialized$(COLOR_RESET)"

build:
	@echo "$(COLOR_YELLOW)Building Docker image...$(COLOR_RESET)"
	docker build -t $(STACK_NAME):latest .
	@echo "$(COLOR_GREEN)✅ Image built successfully!$(COLOR_RESET)"

stack-up:
	@echo "$(COLOR_YELLOW)Deploying stack to swarm with $(ENV) environment...$(COLOR_RESET)"
	@if [ ! -f $(ENV_FILE) ]; then echo "$(COLOR_YELLOW)⚠️  Environment file $(ENV_FILE) not found, creating from template...$(COLOR_RESET)"; cp .env.sample $(ENV_FILE); fi
	@echo "$(COLOR_BLUE)Loading environment from: $(ENV_FILE)$(COLOR_RESET)"
	docker stack deploy -c docker-compose.yml --env-file $(ENV_FILE) $(STACK_NAME)
	@echo "$(COLOR_GREEN)✅ Stack deployed successfully!$(COLOR_RESET)"
	@echo "$(COLOR_BLUE)Use 'make status' to check deployment status$(COLOR_RESET)"

stack-down:
	@echo "$(COLOR_YELLOW)Removing stack from swarm...$(COLOR_RESET)"
	docker stack rm $(STACK_NAME)
	@echo "$(COLOR_GREEN)✅ Stack removed successfully!$(COLOR_RESET)"

# MONITORING TARGETS
status:
	@echo "$(COLOR_YELLOW)Services Status:$(COLOR_RESET)"
	docker stack services $(STACK_NAME)

logs:
	@echo "$(COLOR_YELLOW)Following logs...$(COLOR_RESET)"
	docker service logs -f $(STACK_NAME)_api

ps:
	@echo "$(COLOR_YELLOW)Stack Tasks:$(COLOR_RESET)"
	docker stack ps $(STACK_NAME)

scale:
	@echo "$(COLOR_YELLOW)Scaling API service to $(REPLICAS) replicas...$(COLOR_RESET)"
	docker service scale $(STACK_NAME)_api=$(REPLICAS)
	@echo "$(COLOR_GREEN)✅ Service scaled successfully!$(COLOR_RESET)"

# ENVIRONMENT MANAGEMENT
env-status:
	@echo "$(COLOR_YELLOW)Current Environment Configuration:$(COLOR_RESET)"
	@echo "$(COLOR_BLUE)Active Environment: $(ENV)$(COLOR_RESET)"
	@echo "$(COLOR_BLUE)Environment File: $(ENV_FILE)$(COLOR_RESET)"
	@if [ -f $(ENV_FILE) ]; then echo "$(COLOR_GREEN)✅ Environment file exists$(COLOR_RESET)"; else echo "$(COLOR_YELLOW)⚠️  Environment file missing$(COLOR_RESET)"; fi
	@echo ""

env-create:
	@echo "$(COLOR_YELLOW)Creating environment files...$(COLOR_RESET)"
	@if [ ! -f .env.local ]; then cp .env-sample .env.local && echo "$(COLOR_GREEN)✅ Created .env.local$(COLOR_RESET)"; else echo "$(COLOR_BLUE).env.local already exists$(COLOR_RESET)"; fi
	@if [ ! -f .env.deploy ]; then cp .env-sample .env.deploy && echo "$(COLOR_GREEN)✅ Created .env.deploy$(COLOR_RESET)"; else echo "$(COLOR_BLUE).env.deploy already exists$(COLOR_RESET)"; fi
	@echo "$(COLOR_YELLOW)Please edit the environment files as needed$(COLOR_RESET)"

env-switch:
	@echo "$(COLOR_YELLOW)Switching to $(ENV) environment...$(COLOR_RESET)"
	@if [ ! -f $(ENV_FILE) ]; then echo "$(COLOR_YELLOW)⚠️  Environment file $(ENV_FILE) not found, creating from template...$(COLOR_RESET)"; cp .env-sample $(ENV_FILE); fi
	@cp $(ENV_FILE) .env
	@echo "$(COLOR_GREEN)✅ Switched to $(ENV) environment$(COLOR_RESET)"
	@echo "$(COLOR_BLUE)Active configuration: $(ENV_FILE) → .env$(COLOR_RESET)"
