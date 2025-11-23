package postgres

import (
	"database/sql"

	"github.com/Vitaliy-X/avito-internship-autumn-2025/internal/domain/entities"
	"github.com/Vitaliy-X/avito-internship-autumn-2025/internal/domain/repositories"
)

type TeamRepo struct {
	db *sql.DB
}

func NewTeamRepository(db *sql.DB) repositories.TeamRepository {
	return &TeamRepo{db: db}
}

func (r *TeamRepo) CreateTeam(team *entities.Team) error {
	_, err := r.db.Exec(`
		INSERT INTO teams (team_name) VALUES ($1)
	`, team.Name)
	return err
}

func (r *TeamRepo) GetTeamByName(name string) (*entities.Team, error) {
	row := r.db.QueryRow(`
		SELECT team_name FROM teams WHERE team_name = $1
	`, name)

	var t entities.Team
	err := row.Scan(&t.Name)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TeamRepo) TeamExists(name string) (bool, error) {
	row := r.db.QueryRow(`
		SELECT 1 FROM teams WHERE team_name = $1
	`, name)

	var exists int
	if err := row.Scan(&exists); err != nil {
		return false, nil
	}
	return true, nil
}
