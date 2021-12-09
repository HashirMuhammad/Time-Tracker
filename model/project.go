package model

type Project struct {
	Id          int64  `json:"id"`
	ClientName  string `json:"client_name"`
	StartedBy   string `json:"started_by"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
