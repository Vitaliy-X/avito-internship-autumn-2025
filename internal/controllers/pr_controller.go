package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/Vitaliy-X/avito-internship-autumn-2025/internal/services"
)

type PRController struct {
	prService *services.PRService
}

func NewPRController(svc *services.PRService) *PRController {
	return &PRController{prService: svc}
}

type CreatePRRequest struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
}

type MergePRRequest struct {
	PullRequestID string `json:"pull_request_id"`
}

func (c *PRController) CreatePR(w http.ResponseWriter, r *http.Request) {
	var req CreatePRRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":{"code":"INVALID_REQUEST","message":"invalid JSON"}}`, http.StatusBadRequest)
		return
	}

	pr, err := c.prService.CreatePR(req.PullRequestID, req.PullRequestName, req.AuthorID)
	if err != nil {
		code := err.Error()
		status := http.StatusInternalServerError
		if code == "PR_EXISTS" {
			status = http.StatusConflict
		} else if code == "NOT_FOUND" {
			status = http.StatusNotFound
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(map[string]any{
			"error": map[string]string{"code": code, "message": code},
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"pr": map[string]any{
			"pull_request_id":    string(pr.ID),
			"pull_request_name":  pr.Title,
			"author_id":          pr.AuthorID,
			"status":             pr.Status,
			"assigned_reviewers": pr.Reviewers,
			"created_at":         pr.CreatedAt,
		},
	})
}

func (c *PRController) MergePR(w http.ResponseWriter, r *http.Request) {
	var req MergePRRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":{"code":"INVALID_REQUEST","message":"invalid JSON"}}`, http.StatusBadRequest)
		return
	}
	if req.PullRequestID == "" {
		http.Error(w, `{"error":{"code":"INVALID_REQUEST","message":"pull_request_id is required"}}`, http.StatusBadRequest)
		return
	}

	pr, err := c.prService.MergePR(req.PullRequestID)
	if err != nil {
		code := err.Error()
		status := http.StatusInternalServerError
		if code == "NOT_FOUND" {
			status = http.StatusNotFound
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(map[string]any{
			"error": map[string]string{"code": code, "message": code},
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"pr": map[string]any{
			"pull_request_id":    string(pr.ID),
			"pull_request_name":  pr.Title,
			"author_id":          pr.AuthorID,
			"status":             pr.Status,
			"assigned_reviewers": pr.Reviewers,
			"created_at":         pr.CreatedAt,
			"merged_at":          pr.MergedAt,
		},
	})
}

type ReassignRequest struct {
	PullRequestID string `json:"pull_request_id"`
	OldUserID     string `json:"old_user_id"`
}

func (c *PRController) ReassignReviewer(w http.ResponseWriter, r *http.Request) {
	var req ReassignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":{"code":"INVALID_REQUEST","message":"invalid JSON"}}`, http.StatusBadRequest)
		return
	}
	if req.PullRequestID == "" || req.OldUserID == "" {
		http.Error(w, `{"error":{"code":"INVALID_REQUEST","message":"pull_request_id and old_user_id are required"}}`, http.StatusBadRequest)
		return
	}

	pr, replacedBy, err := c.prService.ReassignReviewer(req.PullRequestID, req.OldUserID)
	if err != nil {
		code := err.Error()
		status := http.StatusInternalServerError

		switch code {
		case "PR_MERGED":
			status = http.StatusConflict // 409
		case "NOT_ASSIGNED":
			status = http.StatusConflict // 409
		case "NO_CANDIDATE":
			status = http.StatusConflict // 409
		case "NOT_FOUND":
			status = http.StatusNotFound // 404
		default:
			status = http.StatusInternalServerError
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(map[string]any{
			"error": map[string]string{
				"code":    code,
				"message": code,
			},
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	resp := map[string]any{
		"pr": map[string]any{
			"pull_request_id":    string(pr.ID),
			"pull_request_name":  pr.Title,
			"author_id":          pr.AuthorID,
			"status":             pr.Status,
			"assigned_reviewers": pr.Reviewers,
			"created_at":         pr.CreatedAt,
			"merged_at":          pr.MergedAt,
		},
		"replaced_by": replacedBy,
	}
	json.NewEncoder(w).Encode(resp)
}
