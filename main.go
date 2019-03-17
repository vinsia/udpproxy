package main

import (
	"github.com/vinsia/udpproxy/udpproxy"
	"gopkg.in/urfave/cli.v1"
	"log"
	"os"
	"time"
)

func initFlags() (app *cli.App) {
	app = cli.NewApp()
	app.Flags = []cli.Flag{
		cli.IntSliceFlag{
			Name: "server_port, s",
			Usage: "listening port",
		},
		cli.IntSliceFlag{
			Name: "client_port, c",
			Usage: "connecting port",
		},
		cli.StringFlag{
			Name: "mode, m",
			Usage: "specify proxy mode",
		},
	}
	return
}

func main() {
	app := initFlags()
	app.Action = func(context *cli.Context) (err error) {
		proxy := udpproxy.NewProxy(context.IntSlice("server_port"), context.IntSlice("client_port"))
		proxy.Init()
		proxy.Start()
		for {
			time.Sleep(time.Duration(10 * time.Hour))
		}
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalf("can not parse arguments: %e", err)
	}
}
