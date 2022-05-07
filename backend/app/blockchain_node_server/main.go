package main

import (
	"blockchain/app/blockchain_node_server/api"
	"blockchain/pkg/command"
	"blockchain/pkg/config"
	"blockchain/pkg/logger"
	"blockchain/pkg/server/cron"
	"blockchain/pkg/server/http"
)

func main() {
	config.LoadConfig()
	cfg := config.GetConfig()
	logger.LoadLogger(cfg)
	defer logger.Logging.Sync()

	c := command.NewCommand()
	httpFlag := http.SetFlag(c)
	cronServer := cron.New(logger.Logging.GetZapLogger())

	a := api.NewApi(cfg, &cronServer)
	httpServer := http.NewServer(a.Router)

	signal := c.WaitSignal(func() {
		cronServer.Start()
		httpServer.Start(
			http.WithFlag(httpFlag),
		)
	})

	cronServer.Stop(signal)
	httpServer.Stop(signal)
}
