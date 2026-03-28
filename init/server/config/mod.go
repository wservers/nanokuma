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

package config

type SSLConfig struct {
	Enable   bool   `toml:"enable"`
	KeyFile  string `toml:"key_file"`
	CertFile string `toml:"cert_file"`
}

type DatabaseConfig struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Name     string `toml:"name"`
	Username string `toml:"username"`
	Password string `toml:"password"`
	Prefix   string `toml:"prefix"`
}

type ServerConfig struct {
	Host     string         `toml:"host"`
	Port     int            `toml:"port"`
	SSL      SSLConfig      `toml:"ssl"`
	Database DatabaseConfig `toml:"database"`
}

const CONFIG_PATH string = "config.toml"

var (
	Get           ServerConfig = ServerConfig{}
	DefaultConfig ServerConfig = ServerConfig{
		Host: "localhost",
		Port: 8080,
		SSL: SSLConfig{
			Enable: false,
		},
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     3306,
			Name:     "nanokuma",
			Username: "root",
			Password: "",
			Prefix:   "nk_",
		},
	}
)
