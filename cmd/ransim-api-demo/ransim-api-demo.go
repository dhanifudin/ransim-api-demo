// SPDX-FileCopyrightText: 2022-present Intel Corporation
// SPDX-FileCopyrightText: 2019-present Open Networking Foundation <info@opennetworking.org>
// SPDX-FileCopyrightText: 2019-present Rimedo Labs
//
// SPDX-License-Identifier: Apache-2.0
// Created by RIMEDO-Labs team
// Based on work of Open Networking Foundation team

package main

import (
	"github.com/RIMEDO-Labs/ransim-api-demo/pkg/manager"
	"github.com/onosproject/onos-lib-go/pkg/logging"
)

var log = logging.GetLogger("main")

func main() {
	log.SetLevel(logging.DebugLevel)

	ready := make(chan bool)

	log.Info("Starting ransim-api-demo")

	cfg := manager.Config{
		AppID:             "ransim-api-demo",
		RansimAddress:     "ran-simulator",
		RansimPort:        5150,
		ApiServingAddress: "",
		ApiServingPort:    8888,
	}

	mgr := manager.NewManager(cfg)

	mgr.Run()

	<-ready

}
