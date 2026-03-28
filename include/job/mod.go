// SPDX-License-Identifier: MIT
// Copyright (C) 2022-2026 WSERVER

package job

import "time"

type JobState string

var (
	Queued  JobState = "queued"
	Running JobState = "running"
	Success JobState = "success"
	Failed  JobState = "failed"
)

type Job struct {
	Id         string    `json:"id"`
	AgentId    *string   `json:"agent_id"`
	ProjectId  string    `json:"repo_url"`
	Branch     string    `json:"branch"`
	Command    string    `json:"command"`
	State      JobState  `json:"state"`
	CreatedAt  time.Time `json:"created_at"`
	StartedAt  time.Time `json:"started_at"`
	FinishedAt time.Time `json:"finished_at"`
	LeaseUntil time.Time `json:"lease_until"`
}

type JobPayload struct {
	RepoUrl string `json:"repo_url"`
	Branch  string `json:"branch"`
	Command string `json:"command"`
}
