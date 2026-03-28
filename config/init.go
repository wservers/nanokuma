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

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
)

func Load[T any](path string, ret *T, defaults T) error {
	var err error
	var buf []byte
	var raw []byte

	buf, err = os.ReadFile(path)
	if err != nil {
		fmt.Printf("[nanokuma-config]: %s is not exists, creating new config file...\n", path)
		raw, err = toml.Marshal(&defaults)
		if err != nil {
			return err
		}

		err = os.WriteFile(path, raw, 0600)
		if err != nil {
			return err
		}
	}

	err = toml.Unmarshal(buf, ret)
	if err != nil {
		return err
	}

	return nil
}

func Unload[T any](ret *T) error {
	if ret != nil {
		ret = nil
	}

	return nil
}
