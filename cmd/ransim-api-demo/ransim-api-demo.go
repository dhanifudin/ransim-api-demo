package main

import (
	"github.com/dhanifudin/ransim-api-demo/pkg/manager"
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
		ApiServingAddress: "127.0.0.1",
		ApiServingPort:    8888,
	}

	mgr := manager.NewManager(cfg)

	mgr.Run()

	<-ready

}
