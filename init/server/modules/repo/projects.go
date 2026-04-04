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

	"git.wh64.net/wserver/nanokuma/include/project"
	"git.wh64.net/wserver/nanokuma/server/config"
	"github.com/google/uuid"
)

func (m *RepoModule) CreateProject(payload project.ProjectPayload) (string, error) {
	var err error
	var id string
	var query string
	var conf = config.Get.Database

	id = uuid.NewString()

	query = fmt.Sprintf("INSERT INTO %sprojects (id, repo_url) VALUES (?, ?);", conf.Prefix)

	_, err = m.DB.Exec(query, id, payload.RepoURL)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (m *RepoModule) GetProject(id string) (*project.Project, error) {
	var err error
	var query string
	var rows *sql.Rows
	var data project.Project
	var conf = config.Get.Database

	query = fmt.Sprintf("SELECT id, repo_url, created_at, updated_at FROM %sprojects WHERE id = ?;", conf.Prefix)

	rows, err = m.DB.Query(query, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if !rows.Next() {
		err = fmt.Errorf("project id '%s' is not exists", id)
		return nil, err
	}

	err = rows.Scan(&data.ID, &data.RepoURL, &data.CreatedAt, &data.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (m *RepoModule) GetProjectByRepoURL(repoURL string) (*project.Project, error) {
	var err error
	var query string
	var rows *sql.Rows
	var data project.Project
	var conf = config.Get.Database

	query = fmt.Sprintf("SELECT id, repo_url, created_at, updated_at FROM %sprojects WHERE repo_url = ?;", conf.Prefix)

	rows, err = m.DB.Query(query, repoURL)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if !rows.Next() {
		err = fmt.Errorf("project repo_url '%s' is not exists", repoURL)
		return nil, err
	}

	err = rows.Scan(&data.ID, &data.RepoURL, &data.CreatedAt, &data.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (m *RepoModule) UpdateProjectRepoURL(id string, repoURL string) error {
	var err error
	var query string
	var conf = config.Get.Database

	query = fmt.Sprintf("UPDATE %sprojects SET repo_url = ?, updated_at = NOW() WHERE id = ?;", conf.Prefix)

	_, err = m.DB.Exec(query, repoURL, id)
	if err != nil {
		return err
	}

	return nil
}

func (m *RepoModule) DeleteProject(id string) error {
	var err error
	var query string
	var conf = config.Get.Database

	query = fmt.Sprintf("DELETE FROM %sprojects WHERE id = ?;", conf.Prefix)

	_, err = m.DB.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}
