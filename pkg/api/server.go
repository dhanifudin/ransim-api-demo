// SPDX-FileCopyrightText: 2022-present Intel Corporation
// SPDX-FileCopyrightText: 2019-present Open Networking Foundation <info@opennetworking.org>
// SPDX-FileCopyrightText: 2019-present Rimedo Labs
//
// SPDX-License-Identifier: Apache-2.0
// Created by RIMEDO-Labs team

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
	GetUes(ctx *gin.Context)
	GetCells(ctx *gin.Context)
	Status(ctx *gin.Context)
	StartServer()
}

type ApiServer struct {
	ransimHandler  ransim.Handler
	servingAddress string
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
	ues, err := serv.ransimHandler.GetUEs(c)
	if err != nil {
		log.Warn("Something's gone wrong when getting the UEs list by API server [GetUEs()].", err)
		c.Status(http.StatusNoContent)
	} else {
		c.IndentedJSON(http.StatusOK, ues)
	}
}

func (serv *ApiServer) GetCells(c *gin.Context) {
	cells, err := serv.ransimHandler.GetCells(c)
	if err != nil {
		log.Warn("Something's gone wrong when getting the Cells list by API server [GetCells()].", err)
		c.Status(http.StatusNoContent)
	} else {
		c.IndentedJSON(http.StatusOK, cells)
	}
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
	router.GET("/get/ues", serv.GetUes)
	router.GET("/get/cells", serv.GetCells)
	router.GET("/status", serv.Status)
	router.Run(serv.servingAddress)
}
