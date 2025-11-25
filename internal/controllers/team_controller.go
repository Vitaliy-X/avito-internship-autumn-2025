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
		JSONError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid JSON")
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
			JSONError(w, http.StatusBadRequest, "TEAM_EXISTS", "team_name already exists")
			return
		}
		JSONError(w, http.StatusInternalServerError, "INTERNAL", err.Error())
		return
	}

	resp := map[string]any{
		"team": map[string]any{
			"team_name": team.Name,
			"members":   req.Members,
		},
	}
	JSONResponse(w, http.StatusCreated, resp)
}

func (c *TeamController) GetTeam(w http.ResponseWriter, r *http.Request) {
	teamName := r.URL.Query().Get("team_name")
	if teamName == "" {
		JSONError(w, http.StatusBadRequest, "INVALID_REQUEST", "team_name is required")
		return
	}

	team, err := c.teamService.GetTeam(teamName)
	if err != nil {
		JSONError(w, http.StatusNotFound, "NOT_FOUND", "resource not found")
		return
	}

	users, err := c.teamService.GetTeamMembers(teamName)
	if err != nil {
		JSONError(w, http.StatusInternalServerError, "INTERNAL", err.Error())
		return
	}

	type memberJSON struct {
		UserID   string `json:"user_id"`
		Username string `json:"username"`
		IsActive bool   `json:"is_active"`
	}

	resp := map[string]any{
		"team_name": team.Name,
		"members":   []memberJSON{},
	}

	for _, u := range users {
		resp["members"] = append(resp["members"].([]memberJSON), memberJSON{
			UserID:   string(u.ID),
			Username: u.Name,
			IsActive: u.IsActive,
		})
	}

	JSONResponse(w, http.StatusOK, resp)
}
