package main

import (
	"log"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
	"github.com/msoap/raphanus"
)

type server struct {
	cfg      config
	raphanus raphanus.DB
}

func newAPI(cfg config) server {
	return server{
		cfg:      cfg,
		raphanus: raphanus.New(),
	}
}

func (app *server) run() {
	echoServer := echo.New()

	echoServer.Use(
		middleware.LoggerWithConfig(middleware.LoggerConfig{Format: `${time_rfc3339} ${remote_ip} ${method} ${path} ${status} ${bytes_out} "${user_agent}"` + "\n"}),
		middleware.Recover(),
	)

	// setup handlers
	v1API := echoServer.Group("/v1")
	v1API.GET("/keys", app.handlerKeys)
	v1API.GET("/length", app.handlerLength)
	v1API.DELETE("/remove/:key", app.handlerRemoveKey)
	// integer methods
	v1API.GET("/int/:key", app.getInt)
	v1API.POST("/int/:key", app.setInt)
	v1API.PUT("/int/:key", app.updateInt)
	v1API.POST("/int/incr/:key", app.incrInt)
	v1API.POST("/int/decr/:key", app.decrInt)

	log.Printf("Server run on %s", defaultAddress)
	echoServer.Run(standard.New(defaultAddress))
}
