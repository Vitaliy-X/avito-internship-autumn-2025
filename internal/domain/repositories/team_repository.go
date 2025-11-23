package repositories

import "github.com/Vitaliy-X/avito-internship-autumn-2025/internal/domain/entities"

type TeamRepository interface {
	CreateTeam(team *entities.Team) error

	GetTeamByName(name string) (*entities.Team, error)

	TeamExists(name string) (bool, error)
}
