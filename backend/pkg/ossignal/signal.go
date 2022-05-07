package ossignal

import (
	"os"
	"os/signal"
	"syscall"
)

func Quit() chan os.Signal {
	quit := make(chan os.Signal, 2)
	signal.Notify(quit,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	)
	return quit
}

func GetExitCode(s os.Signal) int {
	var code int
	switch s {
	case syscall.SIGHUP:
		code = 1
	case syscall.SIGINT:
		code = 2
	case syscall.SIGQUIT:
		code = 3
	case syscall.SIGTERM:
		code = 15
	default:
		code = 9
	}
	return code
}
