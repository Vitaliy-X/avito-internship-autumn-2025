package postgres

import (
	"database/sql"

	"github.com/Vitaliy-X/avito-internship-autumn-2025/internal/domain/entities"
	"github.com/Vitaliy-X/avito-internship-autumn-2025/internal/domain/repositories"
	"github.com/lib/pq"
)

type PRRepo struct {
	db *sql.DB
}

func NewPRRepository(db *sql.DB) repositories.PullRequestRepository {
	return &PRRepo{db: db}
}

func (r *PRRepo) CreatePR(pr *entities.PullRequest) error {
	_, err := r.db.Exec(`
		INSERT INTO pull_requests (pull_request_id, pull_request_name, author_id, status, assigned_reviewers, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, pr.ID, pr.Title, pr.AuthorID, pr.Status, pq.Array(pr.Reviewers), pr.CreatedAt)
	return err
}

func (r *PRRepo) GetPRByID(id string) (*entities.PullRequest, error) {
	row := r.db.QueryRow(`
		SELECT pull_request_id, pull_request_name, author_id, status, assigned_reviewers, created_at, merged_at
		FROM pull_requests
		WHERE pull_request_id = $1
	`, id)

	var pr entities.PullRequest
	var reviewers []string
	err := row.Scan(&pr.ID, &pr.Title, &pr.AuthorID, &pr.Status, pq.Array(&reviewers), &pr.CreatedAt, &pr.MergedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	for _, r := range reviewers {
		pr.Reviewers = append(pr.Reviewers, entities.UserID(r))
	}

	return &pr, nil
}

func (r *PRRepo) UpdatePR(pr *entities.PullRequest) error {
	_, err := r.db.Exec(`
		UPDATE pull_requests
		SET pull_request_name = $1, author_id = $2, status = $3, assigned_reviewers = $4, created_at = $5, merged_at = $6
		WHERE pull_request_id = $7
	`, pr.Title, pr.AuthorID, pr.Status, pq.Array(pr.Reviewers), pr.CreatedAt, pr.MergedAt, pr.ID)
	return err
}

func (r *PRRepo) GetPRsWhereReviewer(userID string) ([]*entities.PullRequest, error) {
	rows, err := r.db.Query(`
		SELECT pull_request_id, pull_request_name, author_id, status, assigned_reviewers, created_at, merged_at
		FROM pull_requests
		WHERE $1 = ANY(assigned_reviewers)
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prs []*entities.PullRequest
	for rows.Next() {
		var pr entities.PullRequest
		var reviewers []string
		if err := rows.Scan(&pr.ID, &pr.Title, &pr.AuthorID, &pr.Status, pq.Array(&reviewers), &pr.CreatedAt, &pr.MergedAt); err != nil {
			return nil, err
		}
		for _, r := range reviewers {
			pr.Reviewers = append(pr.Reviewers, entities.UserID(r))
		}
		prs = append(prs, &pr)
	}

	return prs, nil
}

func (r *PRRepo) AssignReviewers(prID string) ([]entities.UserID, error) {
	rows, err := r.db.Query(`
		SELECT unnest(assign_reviewers_for_pr($1))
	`, prID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviewers []entities.UserID
	for rows.Next() {
		var uid string
		if err := rows.Scan(&uid); err != nil {
			return nil, err
		}
		reviewers = append(reviewers, entities.UserID(uid))
	}

	return reviewers, nil
}
