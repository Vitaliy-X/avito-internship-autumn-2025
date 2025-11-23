package entities

type UserID int64

type User struct {
	ID       UserID
	Name     string
	IsActive bool
}
