package entities

import "time"

type PRID string

type PRStatus string

const (
	PRStatusOpen   PRStatus = "OPEN"
	PRStatusMerged PRStatus = "MERGED"
)

type PullRequest struct {
	ID        PRID
	Title     string
	AuthorID  UserID
	Status    PRStatus
	Reviewers []UserID
	CreatedAt *time.Time
	MergedAt  *time.Time
}
