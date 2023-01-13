package ransim

import (
	"context"
	"fmt"
	"crypto/tls"
	"strconv"

	"github.com/onosproject/onos-lib-go/pkg/certs"
	modelAPI "github.com/onosproject/onos-api/go/onos/ransim/model"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var log = logging.GetLogger("ransim")

type UE struct {
	ID              string  `json:"id"`
	Latitude        float64 `json:"lat"`
	Longitude       float64 `json:"lng"`
	ServingCell     string `json:"serving_cell"`
	RxPower 	float64 `json:"rx_power"`
}

type Cell struct {
	ID        string  `json:"id"`
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
}


func NewHandler(endpoint string) (Handler, error) {

	log.SetLevel(logging.DebugLevel)

	cert, err := tls.X509KeyPair([]byte(certs.DefaultClientCrt), []byte(certs.DefaultClientKey))
	if err != nil {
		return nil, err
	}

	dialOpts := []grpc.DialOption{}
	dialOpts = append(dialOpts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	})))

	conn, err := grpc.Dial(endpoint, dialOpts...)
	if err != nil {
		return nil, err
	}

	log.Info("Dialed gRPC to ransim endpoint")

	return &handler{
		ueClient:   modelAPI.NewUEModelClient(conn),
		cellClient: modelAPI.NewCellModelClient(conn),
	}, nil

}

type Handler interface {
	GetUEs(ctx context.Context) ([]UE, error)
	GetCells(ctx context.Context) ([]Cell, error)
}

type handler struct {
	ueClient   modelAPI.UEModelClient
	cellClient modelAPI.CellModelClient
}

func (h *handler) GetUEs(ctx context.Context) ([]UE, error) {

	stream, err := h.ueClient.ListUEs(context.Background(), &modelAPI.ListUEsRequest{})
	if err != nil {
		log.Warn("Something's gone wrong when getting the UEs info list [GetUEs()].", err)
	}

	results := make([]UE, 0)
	for {
		receiver, err := stream.Recv()
		if err != nil {
			break
		}

		ue := receiver.Ue
		log.Debug(ue)

		ueIdStr := fmt.Sprintf("%d", ue.IMSI)
		ueObj := UE{
			ID: ueIdStr,
			Latitude: ue.Position.Lat, 
			Longitude: ue.Position.Lng,
			ServingCell: strconv.FormatUint(uint64(ue.ServingTower), 16),
			RxPower: ue.ServingTowerStrength,
		}
		results = append(results, ueObj)
	}

	return results, nil

}

func (h *handler) GetCells(ctx context.Context) ([]Cell, error) {

	stream, err := h.cellClient.ListCells(context.Background(), &modelAPI.ListCellsRequest{})
	if err != nil {
		log.Warn("Something's gone wrong when getting the cells info list [GetCells()].", err)
	}

	results := make([]Cell, 0)
	for {
		receiver, err := stream.Recv()
		if err != nil {
			break
		}

		cell := receiver.Cell
		log.Debug(cell)
		
		cellIdStr := fmt.Sprintf("%x", cell.NCGI)
		cellObj := Cell{
			ID: cellIdStr,
			Latitude: cell.Location.Lat, 
			Longitude: cell.Location.Lng,
		}
		results = append(results, cellObj)

	}

	return results, nil
}

