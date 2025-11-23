package postgres

import (
	"database/sql"

	"github.com/Vitaliy-X/avito-internship-autumn-2025/internal/domain/entities"
	"github.com/Vitaliy-X/avito-internship-autumn-2025/internal/domain/repositories"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) repositories.UserRepository {
	return &UserRepo{db: db}
}

func (r *UserRepo) CreateOrUpdateUser(u *entities.User) error {
	_, err := r.db.Exec(`
		INSERT INTO users (user_id, username, is_active, team_name)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id) DO UPDATE
		SET username = EXCLUDED.username,
		    is_active = EXCLUDED.is_active,
		    team_name = EXCLUDED.team_name
	`, u.ID, u.Name, u.IsActive, u.TeamName)
	return err
}

func (r *UserRepo) GetUserByID(id string) (*entities.User, error) {
	row := r.db.QueryRow(`
		SELECT user_id, username, is_active
		FROM users WHERE user_id = $1
	`, id)

	var u entities.User
	err := row.Scan(&u.ID, &u.Name, &u.IsActive)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *UserRepo) GetUsersByTeam(teamName string) ([]*entities.User, error) {
	rows, err := r.db.Query(`
		SELECT user_id, username, is_active
		FROM users
		WHERE team_name = $1
	`, teamName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*entities.User
	for rows.Next() {
		var u entities.User
		rows.Scan(&u.ID, &u.Name, &u.IsActive)
		list = append(list, &u)
	}

	return list, nil
}

func (r *UserRepo) SetIsActive(id string, active bool) (*entities.User, error) {
	_, err := r.db.Exec(`
		UPDATE users SET is_active = $1 WHERE user_id = $2
	`, active, id)

	if err != nil {
		return nil, err
	}

	return r.GetUserByID(id)
}
