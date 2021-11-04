package model

type user struct {
	ID				int64
	FirstName		string
	LastName		string
	Email			string
	Password		string
	ImageUrl		string
	CreatedAt		Time
	UpdatedAt		Time
}