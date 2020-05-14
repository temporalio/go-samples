package main

import (
	"context"
	"time"

	"github.com/pborman/uuid"
	"github.com/temporalio/temporal-go-samples/choice-multi"
	"go.temporal.io/temporal/client"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	// The client is a heavyweight object that should be created once per process.
	c, err := client.NewClient(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		logger.Fatal("Unable to create client", zap.Error(err))
	}
	defer func() { _ = c.CloseConnection() }()

	workflowOptions := client.StartWorkflowOptions{
		ID:                           "multi_choice_" + uuid.New(),
		TaskList:                     "choice-multi",
		ExecutionStartToCloseTimeout: time.Minute,
	}

	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, choice_multi.MultiChoiceWorkflow)
	if err != nil {
		logger.Fatal("Unable to execute workflow", zap.Error(err))
	}
	logger.Info("Started workflow", zap.String("WorkflowID", we.GetID()), zap.String("RunID", we.GetRunID()))

}
