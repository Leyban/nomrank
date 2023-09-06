package shutdown

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

var (
	onlyOneSignalHandler = make(chan struct{})
	shutdownHandler      chan os.Signal
	shutdownSignals      = []os.Signal{os.Interrupt, syscall.SIGTERM}
)

func SetupSignalHandler() <-chan struct{} {
	return SetupSignalContext().Done()
}

func SetupSignalContext() context.Context {
	close(onlyOneSignalHandler)

	shutdownHandler = make(chan os.Signal, 2)

	ctx, cancel := context.WithCancel(context.Background())
	signal.Notify(shutdownHandler, shutdownSignals...)
	go func() {
		<-shutdownHandler
		cancel()
		<-shutdownHandler
		os.Exit(1)
	}()

	return ctx
}
