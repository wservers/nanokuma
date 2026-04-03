package project

import "time"

type Project struct {
	ID        string    `json:"id"`
	RepoURL   string    `json:"repo_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProjectPayload struct {
	RepoURL string `json:"repo_url" binding:"required"`
}
