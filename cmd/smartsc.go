package cmd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/emicklei/go-restful"
	"github.com/urfave/cli"
	"net/http"
	"os"
	"smartsc/internal/controller"
)

const (
	appName = "smartsc"
	version = "dev"
	address = "0.0.0.0:8000"
)

func main() {
	app := cli.NewApp()
	app.Name = appName
	app.Version = version

	app.Action = startSmartsc
	app.Run(os.Args)
}

func startSmartsc(ctx *cli.Context) error {
	log.Info("Starting Smartsc...")

	// main container
	container := restful.NewContainer()

	c := controller.NewController()
	c.Register(container)

	log.Infof("smart scheduler listening at: %v", address)
	return http.ListenAndServe(address, container)
}
