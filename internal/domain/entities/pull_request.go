package entities

type PRID int64

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
}
