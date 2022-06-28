package api

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/RIMEDO-Labs/ransim-api-demo/pkg/ransim"

	"github.com/onosproject/onos-lib-go/pkg/logging"
)

var log = logging.GetLogger("server")


type Api interface {
	Status(ctx *gin.Context)
	StartServer()
}

type ApiServer struct {
	ransimHandler ransim.Handler
	servingAddress string
}


func NewOwnApiServer(ransimEndPoint string, servingAddress string) (Api, error) {

	ransimHandler, err := ransim.NewHandler(ransimEndPoint)
	if err != nil {

		log.Fatal("Unable to run RANSIM Manager", err)
		return nil, err

	}

	return &ApiServer{
		ransimHandler: ransimHandler,
		servingAddress: servingAddress,
	}, nil

}


func (serv *ApiServer) Status(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, true)
}

func (serv *ApiServer) StartServer() {

	router := gin.New()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin, X-Requested-With, Content-Type, Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	router.GET("/status", serv.Status)
	router.Run(serv.servingAddress)
}