package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/Vitaliy-X/avito-internship-autumn-2025/internal/domain/entities"
	"github.com/Vitaliy-X/avito-internship-autumn-2025/internal/services"
)

type TeamController struct {
	teamService *services.TeamService
}

func NewTeamController(service *services.TeamService) *TeamController {
	return &TeamController{teamService: service}
}

type AddTeamRequest struct {
	TeamName string               `json:"team_name"`
	Members  []AddTeamUserRequest `json:"members"`
}

type AddTeamUserRequest struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

func (c *TeamController) AddTeam(w http.ResponseWriter, r *http.Request) {
	var req AddTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	team := &entities.Team{Name: req.TeamName}

	users := make([]*entities.User, 0, len(req.Members))
	for _, m := range req.Members {
		users = append(users, &entities.User{
			ID:       entities.UserID(m.UserID),
			Name:     m.Username,
			IsActive: m.IsActive,
			TeamName: team.Name,
		})
	}

	err := c.teamService.CreateTeam(team, users)
	if err != nil {
		if err.Error() == "TEAM_EXISTS" {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"error":{"code":"TEAM_EXISTS","message":"team_name already exists"}}`))
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := map[string]any{
		"team": map[string]any{
			"team_name": team.Name,
			"members":   req.Members,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
