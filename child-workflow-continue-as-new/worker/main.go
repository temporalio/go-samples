package main

import (
	"os"
	"os/signal"

	cw "github.com/temporalio/temporal-go-samples/child-workflow-continue-as-new"
	"go.temporal.io/temporal/client"
	"go.temporal.io/temporal/worker"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	// The client and worker are heavyweight objects that should be created once per process.
	c, err := client.NewClient(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		logger.Fatal("Unable to create client", zap.Error(err))
	}
	defer func() { _ = c.CloseConnection() }()

	w := worker.New(c, "child-workflow-continue-as-new", worker.Options{
		Logger: logger,
	})
	defer w.Stop()

	w.RegisterWorkflow(cw.SampleParentWorkflow)
	w.RegisterWorkflow(cw.SampleChildWorkflow)

	err = w.Start()
	if err != nil {
		logger.Fatal("Unable to start worker", zap.Error(err))
	}

	// The workers are supposed to be long running process that should not exit.
	waitCtrlC()
}

func waitCtrlC() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
}
