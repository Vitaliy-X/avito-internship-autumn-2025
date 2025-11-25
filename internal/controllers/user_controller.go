package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/Vitaliy-X/avito-internship-autumn-2025/internal/services"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController(service *services.UserService) *UserController {
	return &UserController{userService: service}
}

type SetIsActiveRequest struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

func (c *UserController) SetIsActive(w http.ResponseWriter, r *http.Request) {
	var req SetIsActiveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JSONError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid JSON")
		return
	}

	user, err := c.userService.SetIsActive(req.UserID, req.IsActive)
	if err != nil {
		JSONError(w, http.StatusNotFound, "NOT_FOUND", "user not found")
		return
	}

	resp := map[string]any{
		"user": map[string]any{
			"user_id":   string(user.ID),
			"username":  user.Name,
			"team_name": user.TeamName,
			"is_active": user.IsActive,
		},
	}

	JSONResponse(w, http.StatusOK, resp)
}
