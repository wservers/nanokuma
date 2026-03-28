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

package core

import (
	"fmt"
	"os"
	"slices"
	"sync"
)

type KumaModule interface {
	GetName() string

	Load() error
	Unload() error
}

type NanoKuma struct {
	sync.RWMutex
	orders  []string
	modules map[string]KumaModule

	Initialized bool
}

func NewNanoKuma() *NanoKuma {
	return &NanoKuma{
		orders:  make([]string, 0),
		modules: make(map[string]KumaModule, 0),

		Initialized: false,
	}
}

func (n *NanoKuma) AddModule(module KumaModule) {
	if n.Initialized {
		_, _ = fmt.Fprintf(os.Stderr, "your module is ignored by this system, because service already loaded.\n")
		return
	}

	if module == nil {
		return
	}

	n.orders = append(n.orders, module.GetName())
	n.modules[module.GetName()] = module
}

func (n *NanoKuma) Init() error {
	var err error

	if n.Initialized {
		err = fmt.Errorf("[nanokuma] mininaru core is already loaded")
		goto out_handle_error
	}

	fmt.Printf("Starting nanokuma...\n")
	n.Lock()

	for _, name := range n.orders {
		var module = n.modules[name]
		fmt.Printf("[nanokuma]: loading kuma module: %s\n", name)

		err = module.Load()
		if err != nil {
			err = fmt.Errorf("[nanokuma]: failed to load module %s: %v", name, err)
			goto out_unlock_mutex
		}
	}

	n.Unlock()
	n.Initialized = true

	return nil

out_unlock_mutex:
	n.Unlock()
out_handle_error:
	return err
}

func (n *NanoKuma) Destroy() error {
	var err error
	if !n.Initialized {
		return fmt.Errorf("[mininaru]: nanokuma core's state is already dead")
	}

	fmt.Printf("[nanokuma]: Shutting down nanokuma...\n")
	slices.Reverse(n.orders)
	n.Lock()

	for _, order := range n.orders {
		fmt.Printf("[nanokuma]: unloading kuma module: %s\n", order)
		var module, ok = n.modules[order]
		if !ok {
			continue
		}

		err = module.Unload()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "[nanokuma]: failed to unload module %s: %v\n", order, err)
		}
	}

	n.Unlock()
	n.orders = nil

	n.Initialized = false

	return nil
}
