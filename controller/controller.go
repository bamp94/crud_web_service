package controller

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v4/middleware"

	"crud_web_service/application"
	"crud_web_service/config"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// this variables should be specified by '--ldflags' on the building stage
var Branch, link, author, date, summary string

// Controller is presentation tier of 3-layer architecture
type Controller struct {
	app    application.Application
	config config.Main
	router *echo.Echo
}

// New Controller constructor
func New(config config.Main, app application.Application) Controller {
	return Controller{
		app:    app,
		config: config,
		router: echo.New(),
	}
}

// ServeHTTP http server
func (c Controller) ServeHTTP(ctx context.Context, port int) {
	c.initRoutes()

	go func() {
		if err := c.router.Start(fmt.Sprint(":", port)); err != nil {
			logrus.WithError(err).Fatal("can't start serving http")
		}
	}()

	// Gracefully stopping
	<-ctx.Done()
	if err := c.router.Shutdown(ctx); err != nil {
		logrus.Error("http server shutdown error:", err)
	}
	logrus.Println("http server has stopped")
}

func (c Controller) initRoutes() {
	// cors used for success answer on OPTION request from swagger
	c.router.Use(middleware.CORS())

	c.router.GET("/"+c.config.DB.Endpoint+"/swagger/*any", echoSwagger.WrapHandler)
	c.router.GET("/"+c.config.DB.Endpoint+"v1/healthcheck", c.healthcheck)
}
