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

package database

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"slices"
	"strings"

	"git.wh64.net/wserver/nanokuma/server/config"
	"github.com/go-sql-driver/mysql"
)

type DatabaseModule struct {
	DB *sql.DB
}

//go:embed migrations
var migrations embed.FS

var schema = `CREATE TABLE IF NOT EXISTS migrations (
	version		VARCHAR(255) NOT NULL,
	applied_at	DATETIME DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY(version)
);`

func (m *DatabaseModule) initializeMigrate() error {
	fmt.Printf("[database]: creating migrations schema table...\n")
	var _, err = m.DB.Exec(schema)
	if err != nil {
		return err
	}

	return nil
}

func (m *DatabaseModule) queryMigrations() ([]string, error) {
	var err error
	var rows *sql.Rows
	var versions []string = make([]string, 0)

	fmt.Printf("[database]: get migrations from database file...\n")
	rows, err = m.DB.Query("SELECT version FROM migrations;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var version string
		err = rows.Scan(&version)
		if err != nil {
			return nil, err
		}

		fmt.Printf("[database]: loaded migration version data: %s\n", version)

		versions = append(versions, version)
	}

	return versions, nil
}

func (m *DatabaseModule) applyMigration(tx *sql.Tx, path, version string) error {
	var buf []byte
	var err error

	var conf = config.Get
	var prefix string = ""

	var sql string

	if conf.Database.Prefix != "" {
		prefix = conf.Database.Prefix
	}

	fmt.Printf("[database]: loading embedded migration file: %s\n", path)
	buf, err = fs.ReadFile(migrations, path)
	if err != nil {
		return err
	}

	sql = strings.ReplaceAll(string(buf), "<prefix>", prefix)

	fmt.Printf("[database]: executing migration:\n%s\n", sql)
	_, err = tx.Exec(sql)
	if err != nil {
		return err
	}

	fmt.Printf("[database]: upload version data to migration table: %s\n", version)
	_, err = tx.Exec("INSERT INTO migrations (version) VALUES (?);", version)
	if err != nil {
		return err
	}

	fmt.Printf("[database]: applied migration: %s\n", version)
	return nil
}

func (m *DatabaseModule) migrations() error {
	var versions []string
	var glob []string
	var err error
	var tx *sql.Tx

	err = m.initializeMigrate()
	if err != nil {
		return err
	}

	versions, err = m.queryMigrations()
	if err != nil {
		return err
	}

	glob, err = fs.Glob(migrations, "migrations/*.sql")
	if err != nil {
		return err
	}
	if len(glob) == 0 {
		return nil
	}

	tx, err = m.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, path := range glob {
		var version = strings.ReplaceAll(path, "migrations/", "")
		if slices.Contains(versions, version) {
			log.Printf("[database]: already applied migration: %s\n", version)
			continue
		}

		err = m.applyMigration(tx, path, version)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (m *DatabaseModule) GetName() string {
	return "database-module"
}

func (m *DatabaseModule) Load() error {
	var err error
	var dsn string
	var raw *mysql.Config
	var conf = config.Get.Database

	raw = mysql.NewConfig()
	raw.Net = "tcp"
	raw.User = conf.Username
	raw.Passwd = conf.Password
	raw.Addr = fmt.Sprintf("%s:%d", conf.Host, conf.Port)
	raw.DBName = conf.Name
	raw.Params = map[string]string{
		"charset":         "utf8mb4",
		"parseTime":       "true",
		"loc":             "Asia/Seoul",
		"multiStatements": "true",
	}

	dsn = raw.FormatDSN()

	m.DB, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	fmt.Printf("[database]: connected database to mysql://%s:%d/%s\n", conf.Host, conf.Port, conf.Name)

	err = m.migrations()
	if err != nil {
		return err
	}

	return nil
}

func (m *DatabaseModule) Unload() error {
	_ = m.DB.Close()
	return nil
}

var Database *DatabaseModule = &DatabaseModule{}
