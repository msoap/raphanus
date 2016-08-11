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
		raphanus: raphanus.New(cfg.filename, cfg.syncTime),
	}
}

func (app *server) run() {
	echoServer := echo.New()

	echoServer.Use(
		middleware.LoggerWithConfig(middleware.LoggerConfig{Format: `${time_rfc3339} ${remote_ip} ${method} ${path} ${status} ${bytes_out} "${user_agent}"` + "\n"}),
		middleware.Recover(),
	)

	if len(app.cfg.user) > 0 {
		echoServer.Use(middleware.BasicAuth(func(username, password string) bool {
			if username == app.cfg.user && password == app.cfg.password {
				return true
			}
			return false
		}))
	}

	// setup handlers
	v1API := echoServer.Group("/v1")
	v1API.GET("/stat", app.handlerStat)
	v1API.GET("/keys", app.handlerKeys)
	v1API.GET("/length", app.handlerLength)
	v1API.DELETE("/remove/:key", app.handlerRemoveKey)
	// integer methods
	v1API.GET("/int/:key", app.getInt)
	v1API.POST("/int/:key", app.setInt)
	v1API.PUT("/int/:key", app.updateInt)
	v1API.POST("/int/incr/:key", app.incrInt)
	v1API.POST("/int/decr/:key", app.decrInt)
	// string methods
	v1API.GET("/str/:key", app.getStr)
	v1API.POST("/str/:key", app.setStr)
	v1API.PUT("/str/:key", app.updateStr)
	// list methods
	v1API.GET("/list/:key", app.getList)
	v1API.POST("/list/:key", app.setList)
	v1API.PUT("/list/:key", app.updateList)
	v1API.GET("/list/item/:key", app.getListItem)
	v1API.PUT("/list/item/:key", app.setListItem)
	// dict methods
	v1API.GET("/dict/:key", app.getDict)
	v1API.POST("/dict/:key", app.setDict)
	v1API.PUT("/dict/:key", app.updateDict)
	v1API.GET("/dict/item/:key", app.getDictItem)
	v1API.PUT("/dict/item/:key", app.setDictItem)
	v1API.DELETE("/dict/item/:key", app.removeDictItem)

	log.Printf("Server run on %s", app.cfg.address)
	echoServer.Run(standard.New(app.cfg.address))
}
