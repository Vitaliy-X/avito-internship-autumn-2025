package services

import (
	"errors"
	"time"

	"github.com/Vitaliy-X/avito-internship-autumn-2025/internal/domain/entities"
	"github.com/Vitaliy-X/avito-internship-autumn-2025/internal/domain/repositories"
)

type PRService struct {
	prRepo   repositories.PullRequestRepository
	userRepo repositories.UserRepository
}

func NewPRService(prRepo repositories.PullRequestRepository, userRepo repositories.UserRepository) *PRService {
	return &PRService{prRepo: prRepo, userRepo: userRepo}
}

func (s *PRService) CreatePR(prID, title, authorID string) (*entities.PullRequest, error) {
	existing, _ := s.prRepo.GetPRByID(prID)
	if existing != nil {
		return nil, errors.New("PR_EXISTS")
	}

	author, err := s.userRepo.GetUserByID(authorID)
	if err != nil || author == nil {
		return nil, errors.New("NOT_FOUND")
	}

	pr := &entities.PullRequest{
		ID:        entities.PRID(prID),
		Title:     title,
		AuthorID:  author.ID,
		Status:    entities.PRStatusOpen,
		Reviewers: []entities.UserID{},
		CreatedAt: func() *time.Time { t := time.Now().UTC(); return &t }(),
	}

	if err := s.prRepo.CreatePR(pr); err != nil {
		return nil, err
	}

	reviewers, err := s.prRepo.AssignReviewers(prID)
	if err != nil {
		return nil, err
	}
	pr.Reviewers = reviewers

	return pr, nil
}
