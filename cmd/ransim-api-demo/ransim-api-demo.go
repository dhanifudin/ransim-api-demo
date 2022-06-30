package main

import (
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/RIMEDO-Labs/ransim-api-demo/pkg/manager"
)

var log = logging.GetLogger("main")

func main() {
	log.SetLevel(logging.DebugLevel)

	ready := make(chan bool)

	log.Info("Starting ransim-api-demo")

	cfg := manager.Config{
		AppID:            "ransim-api-demo",
		RansimAddress:    "ran-simulator",
		RansimPort:       5150,
		ApiServingAddress:    "",
		ApiServingPort:       8888,
	}

	mgr := manager.NewManager(cfg)

	mgr.Run()

	<-ready

}