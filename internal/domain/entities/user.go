package entities

type UserID string

type User struct {
	ID       UserID
	Name     string
	IsActive bool
	TeamName string
}
