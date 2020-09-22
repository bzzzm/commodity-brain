package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/bzzzm/commodity-brain/pkg/commodity"
	"github.com/bzzzm/commodity-brain/pkg/utils"
	"github.com/octago/sflags/gen/gflag"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

func main() {

	// create config and parse options
	cfg := utils.NewConfig()
	err := gflag.ParseToDef(&cfg)
	if err != nil {
		panic(fmt.Errorf("cannot initiate config: %v", err))
	}
	flag.Parse()

	// main context
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// create quit and error channels
	qchan := make(chan os.Signal, 10)
	echan := make(chan error, 10)
	defer close(qchan)
	defer close(echan)

	// os signals
	signal.Notify(qchan, syscall.SIGINT, syscall.SIGTERM)

	// logger
	lc := zap.NewDevelopmentConfig()
	lc.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, _ := lc.Build()
	_ = zap.ReplaceGlobals(logger)

	// create commodity and watch buttons
	c := commodity.NewCommodity(ctx, &cfg, echan, qchan)
	defer c.Close()

	// start
	if runtime.GOARCH == "arm" {
		c.Start()
	}
	zap.S().Infof("started commodity v1 on arch=%v", runtime.GOARCH)

	// todo: initiate default program

	for {
		select {
		case <-qchan:
			cancel()
			c.Close()
			zap.S().Info("system signal received, quiting..")
			return
		case <-echan:
			//zap.S().Errorf("received error: %v", err)
			continue
		}
	}
}
