package manager

import (
	"fmt"
	"strconv"

	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/RIMEDO-Labs/ransim-api-demo/pkg/api"
)

var log = logging.GetLogger("manager")


type Config struct {
	AppID             string
	RansimAddress     string
	RansimPort        int
	ApiServingAddress string
	ApiServingPort    int
}

func NewManager(config Config) *Manager {
	log.Infof("Creating manager for: %v", config.AppID)

	apiServer, err := api.NewOwnApiServer(
		config.RansimAddress+":"+strconv.Itoa(config.RansimPort),
		config.ApiServingAddress+":"+strconv.Itoa(config.ApiServingPort),
	)

	if err != nil {
		log.Fatal(err)
	}

	manager := &Manager{
		apiServer: apiServer,
	}

	return manager
}

type Manager struct {
	apiServer api.Api
}

func (m *Manager) Run() {
	log.Info("Running Manager")

	if err := m.start(); err != nil {
		log.Fatal("Unable to run Manager", err)
	}

}

func (m *Manager) start() error {

	if m.apiServer != nil {
		go m.apiServer.StartServer()
	} else {
		return fmt.Errorf("API server is not specified!/n")
	}

	return nil
}