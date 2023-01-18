package api

import (
	"io"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"

	"github.com/dhanifudin/ransim-api-demo/pkg/ransim"

	"github.com/onosproject/onos-lib-go/pkg/logging"
)

var log = logging.GetLogger("server")

type Api interface {
	GetUes(ctx *gin.Context)
	GetCells(ctx *gin.Context)
	Status(ctx *gin.Context)
	StartServer()
}

type ApiServer struct {
	ransimHandler  ransim.Handler
	servingAddress string
}

type UeEvent struct {
	Message       chan []ransim.UE
	NewClients    chan chan []ransim.UE
	ClosedClients chan chan []ransim.UE
	TotalClients  map[chan []ransim.UE]bool
}

type CellEvent struct {
	Message       chan ransim.Cell
	NewClients    chan chan ransim.Cell
	ClosedClients chan chan ransim.Cell
	TotalClients  chan chan ransim.Cell
}

type UeClientChan chan []ransim.UE
type CellClientChan chan []ransim.Cell

func NewUeEvent() (event *UeEvent) {
	event = &UeEvent{
		Message:       make(chan []ransim.UE),
		NewClients:    make(chan chan []ransim.UE),
		ClosedClients: make(chan chan []ransim.UE),
		TotalClients:  make(map[chan []ransim.UE]bool),
	}

	go event.listen()
	return
}

func NewOwnApiServer(ransimEndPoint string, servingAddress string) (Api, error) {

	ransimHandler, err := ransim.NewHandler(ransimEndPoint)
	if err != nil {

		log.Fatal("Unable to run RANSIM Manager", err)
		return nil, err

	}

	return &ApiServer{
		ransimHandler:  ransimHandler,
		servingAddress: servingAddress,
	}, nil

}

func (serv *ApiServer) GetUes(c *gin.Context) {
	// ues, err := serv.ransimHandler.GetUEs(c)
	// if err != nil {
	// 	log.Warn("Something's gone wrong when getting the UEs list by API server [GetUEs()].", err)
	// 	c.Status(http.StatusNoContent)
	// } else {
	// 	c.IndentedJSON(http.StatusOK, ues)
	// }
	value, ok := c.Get("ueClientChan")
	if !ok {
		return
	}
	ueClientChan, ok := value.(UeClientChan)
	if !ok {
		return
	}
	c.Stream(func(w io.Writer) bool {
		if msg, ok := <-ueClientChan; ok {
			c.SSEvent("message", msg)
			return true
		}
		return false
	})
}

func (serv *ApiServer) GetCells(c *gin.Context) {
	// cells, err := serv.ransimHandler.GetCells(c)
	// if err != nil {
	// 	log.Warn("Something's gone wrong when getting the Cells list by API server [GetCells()].", err)
	// 	c.Status(http.StatusNoContent)
	// } else {
	// 	c.IndentedJSON(http.StatusOK, cells)
	// }
}

func (serv *ApiServer) Status(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, true)
}

func (stream *UeEvent) listen() {
	for {
		select {
		case client := <-stream.NewClients:
			stream.TotalClients[client] = true
		case client := <-stream.ClosedClients:
			delete(stream.TotalClients, client)
			close(client)

		case eventMsg := <-stream.Message:
			for clientMessageChan := range stream.TotalClients {
				clientMessageChan <- eventMsg
			}
		}
	}
}

func (stream *UeEvent) streamUes() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ueClientChan := make(UeClientChan)
		stream.NewClients <- ueClientChan

		defer func() {
			stream.ClosedClients <- ueClientChan
		}()

		ctx.Set("ueClientChan", ueClientChan)
		ctx.Next()
	}
}

func headersMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Content-Type", "text/event-stream")
		ctx.Writer.Header().Set("Cache-Control", "no-cache")
		ctx.Writer.Header().Set("Connection", "keep-alive")
		ctx.Writer.Header().Set("Transfer-Encoding", "chunked")
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Next()
	}
}

func (serv *ApiServer) StartServer() {
	ueStream := NewUeEvent()

	go func() {
		for {
			time.Sleep(time.Second * 1)
			ues, err := serv.ransimHandler.GetUEs(context.Background())
			if err != nil {
				ueStream.Message <- ues
			}
		}
	}()

	router := gin.New()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin, X-Requested-With, Content-Type, Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	router.RemoveExtraSlash = true
	router.GET("/get/ues", headersMiddleware(), serv.GetUes)
	// router.GET("/get/cells", headersMiddleware(), serv.GetCells)
	router.GET("/status", serv.Status)
	router.Run(serv.servingAddress)
}
