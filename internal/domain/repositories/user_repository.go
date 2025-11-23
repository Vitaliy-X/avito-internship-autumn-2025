package repositories

import "github.com/Vitaliy-X/avito-internship-autumn-2025/internal/domain/entities"

type UserRepository interface {
	CreateOrUpdateUser(user *entities.User) error

	GetUserByID(id string) (*entities.User, error)

	GetUsersByTeam(teamName string) ([]*entities.User, error)

	SetIsActive(id string, active bool) (*entities.User, error)
}
