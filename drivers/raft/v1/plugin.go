// SPDX-FileCopyrightText: 2023-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/vpascoalr/atomix/drivers/raft/v1/driver"
	"github.com/vpascoalr/atomix/runtime/pkg/network"
)

var Plugin = driver.New(network.NewDefaultDriver())
