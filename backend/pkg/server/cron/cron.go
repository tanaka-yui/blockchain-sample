package cron

import (
	"blockchain/pkg/logger"
	"fmt"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"os"
)

type (
	Job interface {
		Start()
		Stop(sig os.Signal)
		AddFunc(spec string, cmd func())
	}
	Function interface {
		Interval() string
		RegisterJob()
	}
	jobImpl struct {
		cron *cron.Cron
	}
)

func New(zapLogger *zap.Logger) Job {
	job := jobImpl{
		cron: cron.New(cron.WithChain(
			cron.SkipIfStillRunning(&jobLogger{logger: zapLogger})),
		),
	}

	return &job
}

func (j *jobImpl) Start() {
	go j.cron.Start()
	logger.Logging.Info(fmt.Sprintf("Start cron server."))
}

func (j *jobImpl) Stop(sig os.Signal) {
	<-j.cron.Stop().Done()
	logger.Logging.Info(fmt.Sprintf("stopping cron server... Signal: %s", sig.String()))
}

func (j *jobImpl) AddFunc(spec string, cmd func()) {
	_, err := j.cron.AddFunc(spec, cmd)
	if err != nil {
		logger.Logging.Fatal("add function error")
	}
}
