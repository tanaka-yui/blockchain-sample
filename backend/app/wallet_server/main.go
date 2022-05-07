package main

import (
	"blockchain/app/wallet_server/api"
	"blockchain/pkg/command"
	"blockchain/pkg/config"
	"blockchain/pkg/logger"
	"blockchain/pkg/server/http"
)

func main() {
	config.LoadConfig()
	cfg := config.GetConfig()
	logger.LoadLogger(cfg)
	defer logger.Logging.Sync()

	c := command.NewCommand()
	httpFlag := http.SetFlag(c)

	a := api.NewApi(cfg)
	httpServer := http.NewServer(a.Router)

	signal := c.WaitSignal(func() {
		httpServer.Start(
			http.WithFlag(httpFlag),
		)
	})

	httpServer.Stop(signal)
}
