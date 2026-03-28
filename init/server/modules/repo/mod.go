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

	"git.wh64.net/wserver/nanokuma/server/modules/database"
)

type RepoModule struct {
	DB *sql.DB
}

func (*RepoModule) GetName() string {
	return "repo"
}

func (m *RepoModule) Load() error {
	if database.Database == nil {
		return fmt.Errorf("[repo]: database module is not loaded")
	}

	m.DB = database.Database.DB

	return nil
}

func (m *RepoModule) Unload() error {
	if m.DB != nil {
		m.DB = nil
	}

	return nil
}

var Repo *RepoModule = &RepoModule{}
