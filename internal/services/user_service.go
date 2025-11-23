package services

import (
	"github.com/Vitaliy-X/avito-internship-autumn-2025/internal/domain/entities"
	"github.com/Vitaliy-X/avito-internship-autumn-2025/internal/domain/repositories"
)

type UserService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) GetActiveUsersByTeam(teamName string, exclude []entities.UserID) ([]*entities.User, error) {
	users, err := s.userRepo.GetUsersByTeam(teamName)
	if err != nil {
		return nil, err
	}

	excludeMap := make(map[entities.UserID]struct{}, len(exclude))
	for _, id := range exclude {
		excludeMap[id] = struct{}{}
	}

	var active []*entities.User
	for _, u := range users {
		if u.IsActive {
			if _, ok := excludeMap[u.ID]; !ok {
				active = append(active, u)
			}
		}
	}

	return active, nil
}

func (s *UserService) SetIsActive(userID string, active bool) (*entities.User, error) {
	user, err := s.userRepo.SetIsActive(userID, active)
	if err != nil {
		return nil, err
	}
	return user, nil
}
