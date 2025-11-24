package app

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/Vitaliy-X/avito-internship-autumn-2025/internal/config"
	"github.com/Vitaliy-X/avito-internship-autumn-2025/internal/controllers"
	"github.com/Vitaliy-X/avito-internship-autumn-2025/internal/db"
	"github.com/Vitaliy-X/avito-internship-autumn-2025/internal/domain/repositories/postgres"
	"github.com/Vitaliy-X/avito-internship-autumn-2025/internal/services"
)

type App struct {
	cfg *config.Config
}

func New() *App {
	return &App{
		cfg: config.Load(),
	}
}

func (a *App) Run() error {
	dbConn, err := db.Connect(a.cfg)
	if err != nil {
		return err
	}

	teamRepo := postgres.NewTeamRepository(dbConn)
	userRepo := postgres.NewUserRepository(dbConn)
	prRepo := postgres.NewPRRepository(dbConn)

	teamService := services.NewTeamService(teamRepo, userRepo)
	teamController := controllers.NewTeamController(teamService)

	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	prService := services.NewPRService(prRepo, userRepo)
	prController := controllers.NewPRController(prService)

	r := chi.NewRouter()

	r.Post("/team/add", teamController.AddTeam)
	r.Get("/team/get", teamController.GetTeam)

	r.Post("/users/setIsActive", userController.SetIsActive)

	r.Post("/pullRequest/create", prController.CreatePR)
	r.Post("/pullRequest/merge", prController.MergePR)

	log.Printf("Server started on :%s\n", a.cfg.AppPort)
	return http.ListenAndServe(":"+a.cfg.AppPort, r)
}
