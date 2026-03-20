// SPDX-License-Identifier: GPL-2.0-or-later
/*
 * nanokuma
 * Copyright (C) 2022-2026 WSERVER
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 2 of the License, or
 * (at your option) any later version.
 */

package repo

import (
	"database/sql"
	"fmt"

	"git.wh64.net/wserver/nanokuma/server/config"
	"git.wh64.net/wserver/nanokuma/shared/job"
	"github.com/google/uuid"
)

func (m *RepoModule) CreateJob(payload *job.JobPayload) (string, error) {
	var err error
	var id string
	var query string
	var conf = config.Get.Database

	id = uuid.NewString()

	query = fmt.Sprintf("INSERT INTO %sjobs (id, repo_url, branch, command, state) VALUES (?, ?, ?, ?, 'queued');", conf.Prefix)

	_, err = m.DB.Exec(query, id, payload.RepoUrl, payload.Branch, payload.Command)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (m *RepoModule) GetJob(id string) (*job.Job, error) {
	var err error
	var query string
	var rows *sql.Rows
	var conf = config.Get.Database

	var data job.Job

	query = fmt.Sprintf(`SELECT id, agent_id, project_id, branch, command, %s, created_at, started_at, finished_at, lease_until
		FROM %sjobs
		WHERE id = ?;`,
		`state`,
		conf.Prefix,
	)

	rows, err = m.DB.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		err = fmt.Errorf("project id '%s' is not exists", id)
		return nil, err
	}

	err = rows.Scan(
		&data.Id, &data.AgentId, &data.ProjectId, &data.Branch,
		&data.Command, &data.State, &data.CreatedAt,
		&data.StartedAt, &data.FinishedAt, &data.LeaseUntil,
	)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (m *RepoModule) GetJobs(projectId string) ([]job.Job, error) {
	var err error
	var query string
	var rows *sql.Rows
	var conf = config.Get.Database

	var cur job.Job
	var jobs []job.Job = make([]job.Job, 0)

	query = fmt.Sprintf("SELECT id, agent_id, project_id, branch, command, %s, created_at, started_at, finished_at, lease_until FROM %sjobs WHERE id = ?;",
		`state`,
		conf.Prefix,
	)

	rows, err = m.DB.Query(query, projectId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		cur = job.Job{}

		err = rows.Scan(
			&cur.Id, &cur.AgentId, &cur.ProjectId, &cur.Branch,
			&cur.Command, &cur.State, &cur.CreatedAt,
			&cur.StartedAt, &cur.FinishedAt, &cur.LeaseUntil,
		)
		if err != nil {
			return nil, err
		}

		jobs = append(jobs, cur)
	}

	return jobs, nil
}

func (m *RepoModule) UpdateJobState(id string, state job.JobState) error {
	var err error
	var query string
	var conf = config.Get.Database

	query = fmt.Sprintf("UPDATE %sjobs SET `state` = ? WHERE id = ?;", conf.Prefix)

	_, err = m.DB.Exec(query, state, id)
	if err != nil {
		return err
	}

	return nil
}

func (m *RepoModule) PollingJob(agentId string, timeout int) (*job.Job, error) {
	var err error

	var tx *sql.Tx
	var row *sql.Rows
	var query string
	var conf = config.Get.Database

	var job job.Job

	tx, err = m.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query = fmt.Sprintf(
		"SELECT id, agent_id, project_id, branch, command, `state`, created_at, started_at, finished_at, lease_until FROM %sjobs WHERE state = 'queued' ORDER BY created_at LIMIT 1 FOR UPDATE",
		conf.Prefix,
	)

	row, err = tx.Query(query)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}

		return nil, nil
	}
	defer row.Close()

	if !row.Next() {
		return nil, nil
	}

	err = row.Scan(
		&job.Id, &job.AgentId, &job.ProjectId,
		&job.Branch, &job.Command, &job.State, &job.CreatedAt,
		&job.StartedAt, &job.FinishedAt, &job.LeaseUntil,
	)
	if err != nil {
		return nil, err
	}

	query = fmt.Sprintf(
		"UPDATE %sjobs SET `state` = 'running', agent_id = ?, started_at = NOW(), lease_until = DATE_ADD(NOW(), INTERVAL ? MINUTE) WHERE id = ?;",
		conf.Prefix,
	)

	_, err = tx.Exec(query, agentId, timeout, job.Id)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &job, nil
}

func (m *RepoModule) DeleteJob(id string) error {
	var err error
	var query string
	var conf = config.Get.Database

	query = fmt.Sprintf("DELETE FROM %sjobs WHERE id = ?;", conf.Prefix)

	_, err = m.DB.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}
