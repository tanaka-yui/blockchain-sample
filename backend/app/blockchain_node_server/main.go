package main

import (
	"blockchain/app/blockchain_node_server/api"
	"blockchain/pkg/command"
	"blockchain/pkg/config"
	"blockchain/pkg/logger"
	"blockchain/pkg/server/cron"
	"blockchain/pkg/server/http"
	"blockchain/pkg/utils/stringutil"
	"flag"
	"fmt"
	"strings"
)

func main() {
	config.LoadConfig()
	cfg := config.GetConfig()
	logger.LoadLogger(cfg)
	defer logger.Logging.Sync()

	c := command.NewCommand()
	httpFlag := http.SetFlag(c)
	var addr string
	flag.StringVar(&addr, "addr", "0.0.0.0:8080", "TCP Port Number for Wallet Server")
	flag.Parse()
	logger.Logging.Info(fmt.Sprintf("localAddr:%s", addr))

	port := stringutil.ConvertStringToUint16(strings.Split(addr, ":")[1])
	cronServer := cron.New(logger.Logging.GetZapLogger())
	a := api.NewApi(port, cfg, &cronServer)
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
