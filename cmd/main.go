package main

import (
	"log"

	"github.com/ericolvr/maintenance-v2/config"
	"github.com/ericolvr/maintenance-v2/internal/handlers"
	"github.com/ericolvr/maintenance-v2/internal/repository"
	"github.com/ericolvr/maintenance-v2/internal/routes"
	"github.com/ericolvr/maintenance-v2/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()
	db := config.GetDB()

	// Repositories
	branchRepo := repository.NewBranchRepository(db)
	clientRepo := repository.NewClientRepository(db)
	costRepo := repository.NewCostRepository(db)
	distanceRepo := repository.NewDistanceRepository(db)
	providerRepo := repository.NewProviderRepository(db)
	ticketRepo := repository.NewTicketRepository(db)
	userRepo := repository.NewUserRepository(db)
	problemRepo := repository.NewProblemRepository(db)
	solutionRepo := repository.NewSolutionRepository(db)
	slaRepo := repository.NewSlaRepository(db)

	// Services
	branchService := service.NewBranchService(branchRepo)
	clientService := service.NewClientService(clientRepo)
	costService := service.NewCostService(costRepo)
	distanceService := service.NewDistanceService(distanceRepo)
	providerService := service.NewProviderService(providerRepo)
	ticketService := service.NewTicketService(ticketRepo, branchRepo, providerRepo, problemRepo, solutionRepo, distanceService)
	userService := service.NewUserService(userRepo, []byte(cfg.JWTSecret))
	problemService := service.NewProblemService(problemRepo)
	solutionService := service.NewSolutionService(solutionRepo, problemRepo)
	slaService := service.NewSlaService(slaRepo)

	router := gin.Default()

	// CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Routes
	routes.BranchRoutes(router, handlers.NewBranchHandler(branchService))
	routes.ClientRoutes(router, handlers.NewClientHandler(clientService))
	routes.CostRoutes(router, handlers.NewCostHandler(costService))
	routes.DistanceRoutes(router, handlers.NewDistanceHandler(distanceService))
	routes.ProviderRoutes(router, handlers.NewProviderHandler(providerService))
	routes.TicketRoutes(router, handlers.NewTicketHandler(ticketService), handlers.NewTicketProblemHandler(ticketService), handlers.NewTicketSolutionHandler(ticketService))
	routes.UserRoutes(router, handlers.NewUserHandler(userService))
	routes.ProblemRoutes(router, handlers.NewProblemHandler(problemService))
	routes.SolutionRoutes(router, handlers.NewSolutionHandler(solutionService))
	routes.SlaRoutes(router, handlers.NewSlaHandler(slaService))

	log.Printf(
		"Server is running on port %s", cfg.ServerPort,
	)

	err := router.Run(":" + cfg.ServerPort)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
