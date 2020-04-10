package main

import (
	"context"
	"os"
	"os/signal"
	"sync"

	"crud_web_service/application"
	"crud_web_service/config"
	"crud_web_service/controller"
	"crud_web_service/model"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

// @title CRUD Web Service
// @description Документация для http сервера приложения
func main() {
	appCli := cli.NewApp()
	if appCli.Version == "" {
		appCli.Version = "version not specified"
	}

	appCli.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config",
			Value: "/etc/crud_web_service/config.json",
			Usage: "optional config path",
		},
	}

	appCli.Action = func(cliContext *cli.Context) {
		cli.ShowVersion(cliContext)
		config := config.New(cliContext.String("config"))

		m := model.NewFromConfig(config.DB)
		if err := m.CheckMigrations(); err != nil {
			logrus.WithError(err).Fatal("invalid database condition")
		}

		shutdown := make(chan int, 16)
		wg := sync.WaitGroup{}
		appNew := application.New(m, config)
		c := controller.New(config, appNew)
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			for {
				select {
				case <-shutdown:
					cancel()
					return
				}
			}
		}()

		wg.Add(1)
		go func() {
			c.ServeHTTP(ctx, config.Port)
			wg.Done()
		}()

		gracefulClosing(&wg, shutdown)
	}

	appCli.Commands = []cli.Command{
		{
			Name:  "migrate",
			Usage: "update migrations to the latest stage",
			Action: func(cliContext *cli.Context) {
				config := config.New(cliContext.GlobalString("config"))
				m := model.NewFromConfig(config.DB)
				m.Migrate()
			},
		},
	}

	if err := appCli.Run(os.Args); err != nil {
		panic(err)
	}
}

func gracefulClosing(servicesWg *sync.WaitGroup, shutdown chan int) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig
	logrus.Info("stopping services... (Double enter Ctrl + C to force close)")

	close(shutdown)

	quit := make(chan struct{})
	go func() {
		<-sig
		<-sig
		logrus.Info("services unsafe stopped")
		<-quit
	}()

	go func() {
		servicesWg.Wait()
		logrus.Info("services gracefully stopped")
		<-quit
	}()

	quit <- struct{}{}
	close(quit)
}
