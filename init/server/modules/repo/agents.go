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
	"net/http"

	"git.wh64.net/wserver/nanokuma/include/agent"
	"git.wh64.net/wserver/nanokuma/server/config"
)

func (m *RepoModule) UpsertAgent(payload *agent.AgentData) error {
	var err error
	var query string
	var conf = config.Get.Database

	var resp *http.Response
	var status agent.AgentStatus

	query = fmt.Sprintf(
		`INSERT INTO %sagents (id, ip_addr, port, hostname, %s) VALUES (?, ?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE
			ip_addr = ?,
			port = ?,
			hostname = ?,
			%s = ?;`,
		conf.Prefix,
		"`status`",
		"`status`",
	)

	resp, err = http.Get(fmt.Sprintf("http://%s:%d/health", payload.IPAddr, payload.Port))
	if err == nil {
		if resp.StatusCode != 200 {
			status = agent.Offline
		} else {
			status = agent.Online
		}
	} else {
		status = agent.Offline
	}

	_, err = m.DB.Exec(query,
		payload.Id, payload.IPAddr, payload.Port, payload.Hostname, status,
		payload.IPAddr, payload.Port, payload.Hostname, status)
	if err != nil {
		return err
	}

	return nil
}

func (m *RepoModule) GetAgent(id string) (*agent.AgentData, error) {
	var err error
	var query string
	var rows *sql.Rows
	var conf = config.Get.Database

	var data agent.AgentData

	query = fmt.Sprintf("SELECT id, ip_addr, port, hostname, authorized, `status`, last_action_at FROM %sagents WHERE id = ?;", conf.Prefix)

	rows, err = m.DB.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		err = fmt.Errorf("agent id '%s' is not exists", id)
		return nil, err
	}

	err = rows.Scan(&data.Id, &data.IPAddr, &data.Port, &data.Hostname, &data.Authorized, &data.Status, &data.LastActionAt)
	if err != nil {
		return nil, err
	}

	return &data, err
}

func (m *RepoModule) GetAgents() ([]agent.AgentData, error) {
	var err error
	var query string
	var rows *sql.Rows
	var conf = config.Get.Database

	var cur agent.AgentData
	var agents []agent.AgentData = make([]agent.AgentData, 0)

	query = fmt.Sprintf("SELECT id, ip_addr, port, hostname, authorized, `status`, last_action_at FROM %sagents;", conf.Prefix)

	rows, err = m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		cur = agent.AgentData{}
		err = rows.Scan(&cur.Id, &cur.IPAddr, &cur.Port, &cur.Hostname, &cur.Authorized, &cur.Status, &cur.LastActionAt)
		if err != nil {
			return nil, err
		}

		agents = append(agents, cur)
	}

	return agents, nil
}

func (m *RepoModule) AuthorizeAgent(id string) error {
	var err error
	var query string
	var conf = config.Get.Database

	query = fmt.Sprintf("UPDATE %sagents SET authorized = true WHERE id = ?;", conf.Prefix)

	_, err = m.DB.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

func (m *RepoModule) DeleteAgent(id string) error {
	var err error
	var query string
	var conf = config.Get.Database

	query = fmt.Sprintf("DELETE FROM %sagents WHERE id = ?;", conf.Prefix)

	_, err = m.DB.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}
