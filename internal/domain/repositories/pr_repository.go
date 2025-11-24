package repositories

import "github.com/Vitaliy-X/avito-internship-autumn-2025/internal/domain/entities"

type PullRequestRepository interface {
	CreatePR(pr *entities.PullRequest) error
	GetPRByID(id string) (*entities.PullRequest, error)
	UpdatePR(pr *entities.PullRequest) error
	GetPRsWhereReviewer(userID string) ([]*entities.PullRequest, error)
	AssignReviewers(prID string) ([]entities.UserID, error)
	ReassignReviewer(prID string, oldUserID string) (string, error)
}
