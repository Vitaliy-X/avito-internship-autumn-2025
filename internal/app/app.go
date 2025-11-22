package app

import (
	"log"
	"net/http"

	"github.com/Vitaliy-X/avito-internship-autumn-2025/internal/config"
	"github.com/Vitaliy-X/avito-internship-autumn-2025/internal/db"
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
	_, err := db.Connect(a.cfg)
	if err != nil {
		return err
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("PR service is running"))
	})

	log.Printf("Server started on :%s\n", a.cfg.AppPort)
	return http.ListenAndServe(":"+a.cfg.AppPort, nil)
}
