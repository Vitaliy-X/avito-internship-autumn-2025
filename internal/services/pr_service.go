package services

import (
	"errors"
	"math/rand"
	"strconv"
	"time"

	"github.com/Vitaliy-X/avito-internship-autumn-2025/internal/domain/entities"
	"github.com/Vitaliy-X/avito-internship-autumn-2025/internal/domain/repositories"
)

type PRService struct {
	prRepo  repositories.PullRequestRepository
	userSvc *UserService
}

func NewPRService(prRepo repositories.PullRequestRepository, userSvc *UserService) *PRService {
	return &PRService{
		prRepo:  prRepo,
		userSvc: userSvc,
	}
}

func (s *PRService) CreatePR(prID entities.PRID, title string, authorID entities.UserID) (*entities.PullRequest, error) {
	author, err := s.userSvc.userRepo.GetUserByID(string(authorID))
	if err != nil || author == nil {
		return nil, errors.New("NOT_FOUND")
	}

	pr := &entities.PullRequest{
		ID:        prID,
		Title:     title,
		AuthorID:  authorID,
		Status:    entities.PRStatusOpen,
		Reviewers: []entities.UserID{},
	}

	candidates, _ := s.userSvc.GetActiveUsersByTeam(author.Name, []entities.UserID{authorID})
	pr.Reviewers = pickRandomReviewers(candidates, 2)

	if err := s.prRepo.CreatePR(pr); err != nil {
		return nil, errors.New("PR_EXISTS")
	}

	return pr, nil
}

func (s *PRService) ReassignReviewer(prID entities.PRID, oldReviewerID entities.UserID) (entities.UserID, error) {
	pr, err := s.prRepo.GetPRByID(string(prID))
	if err != nil || pr == nil {
		return entities.UserID(strconv.Itoa(0)), errors.New("NOT_FOUND")
	}
	if pr.Status == entities.PRStatusMerged {
		return entities.UserID(strconv.Itoa(0)), errors.New("PR_MERGED")
	}

	if !contains(pr.Reviewers, oldReviewerID) {
		return entities.UserID(strconv.Itoa(0)), errors.New("NOT_ASSIGNED")
	}

	oldUser, _ := s.userSvc.userRepo.GetUserByID(string(oldReviewerID))
	candidates, _ := s.userSvc.GetActiveUsersByTeam(oldUser.Name, pr.Reviewers)
	if len(candidates) == 0 {
		return entities.UserID(strconv.Itoa(0)), errors.New("NO_CANDIDATE")
	}

	newReviewer := candidates[0].ID
	for i, r := range pr.Reviewers {
		if r == oldReviewerID {
			pr.Reviewers[i] = newReviewer
			break
		}
	}

	_ = s.prRepo.UpdatePR(pr)
	return newReviewer, nil
}

func (s *PRService) GetPRsForReviewer(userID entities.UserID) ([]*entities.PullRequest, error) {
	return s.prRepo.GetPRsWhereReviewer(string(userID))
}

func (s *PRService) MergePR(prID entities.PRID) (*entities.PullRequest, error) {
	pr, err := s.prRepo.GetPRByID(string(prID))
	if err != nil || pr == nil {
		return nil, errors.New("NOT_FOUND")
	}
	if pr.Status != entities.PRStatusMerged {
		now := time.Now()
		pr.Status = entities.PRStatusMerged
		pr.MergedAt = &now
		_ = s.prRepo.UpdatePR(pr)
	}
	return pr, nil
}

func pickRandomReviewers(users []*entities.User, max int) []entities.UserID {
	if len(users) == 0 {
		return nil
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(users), func(i, j int) { users[i], users[j] = users[j], users[i] })

	n := max
	if len(users) < max {
		n = len(users)
	}
	res := make([]entities.UserID, n)
	for i := 0; i < n; i++ {
		res[i] = users[i].ID
	}
	return res
}

func contains(slice []entities.UserID, val entities.UserID) bool {
	for _, s := range slice {
		if s == val {
			return true
		}
	}
	return false
}
