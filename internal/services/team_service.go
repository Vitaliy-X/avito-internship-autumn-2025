package services

import (
	"errors"

	"github.com/Vitaliy-X/avito-internship-autumn-2025/internal/domain/entities"
	"github.com/Vitaliy-X/avito-internship-autumn-2025/internal/domain/repositories"
)

type TeamService struct {
	teamRepo repositories.TeamRepository
	userRepo repositories.UserRepository
}

func NewTeamService(teamRepo repositories.TeamRepository, userRepo repositories.UserRepository) *TeamService {
	return &TeamService{
		teamRepo: teamRepo,
		userRepo: userRepo,
	}
}

func (s *TeamService) CreateTeam(team *entities.Team, members []*entities.User) error {
	exists, err := s.teamRepo.TeamExists(team.Name)
	if err != nil {
		return err
	}

	if !exists {
		if err := s.teamRepo.CreateTeam(team); err != nil {
			return err
		}
	}

	for _, u := range members {
		u.TeamName = team.Name
		if err := s.userRepo.CreateOrUpdateUser(u); err != nil {
			return err
		}
	}

	return nil
}

func (s *TeamService) GetTeam(name string) (*entities.Team, error) {
	team, err := s.teamRepo.GetTeamByName(name)
	if err != nil || team == nil {
		return nil, errors.New("NOT_FOUND")
	}
	return team, nil
}

func (s *TeamService) GetTeamMembers(teamName string) ([]*entities.User, error) {
	return s.userRepo.GetUsersByTeam(teamName)
}
