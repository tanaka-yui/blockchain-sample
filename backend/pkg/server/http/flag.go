package http

import (
	"blockchain/pkg/command"
	"blockchain/pkg/config"
	"time"
)

const (
	flagAddr    = "addr"
	flagTimeout = "timeout"
)

type flag struct {
	addr    string
	timeout time.Duration
}

func SetFlag(c *command.Command) *flag {
	cfg := config.GetConfig().System.Http
	flag := &flag{}
	c.Flags().StringVar(&flag.addr, flagAddr, cfg.Addr, "http Addr")
	c.Flags().DurationVar(&flag.timeout, flagTimeout, cfg.ContextTimeoutSec, "http context Timeout")
	return flag
}
