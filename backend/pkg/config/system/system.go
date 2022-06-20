package system

import (
	"blockchain/pkg/utils/stringutil"
	"os"
	"time"
)

type Config struct {
	Env      string
	TimeZone string
	Http     Http
}

type Http struct {
	Addr                   string
	ContextTimeoutSec      time.Duration
	SessionSecure          bool
	SessionExpireMinute    int64
	OAuthTokenExpireMinute int64
	NodeGateway            string
}

func NewConfig() Config {
	http := Http{
		Addr: os.Getenv("HTTP_ADDR"),
		// nolint:durationcheck // ConvertDurationは最小単位のDurationを返すため、Secondをかけるのが正しい
		ContextTimeoutSec:      stringutil.ConvertDuration(os.Getenv("HTTP_CONTEXT_TIMEOUT_SEC")) * time.Second,
		SessionSecure:          stringutil.ConvertBool(os.Getenv("HTTP_SESSION_SECURE")),
		SessionExpireMinute:    stringutil.ConvertStringToInt64(os.Getenv("HTTP_SESSION_EXPIRE_MINUTE")),
		OAuthTokenExpireMinute: stringutil.ConvertStringToInt64(os.Getenv("OAUTH_TOKEN_EXPIRE_MINUTE")),
		NodeGateway:            os.Getenv("NODE_GATEWAY"),
	}
	return Config{
		Env:      os.Getenv("ENV"),
		TimeZone: os.Getenv("TZ"),
		Http:     http,
	}
}
