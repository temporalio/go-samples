package main

import (
	"context"
	"github.com/temporalio/temporal-go-samples/branch"
	"go.temporal.io/temporal/client"
	"go.uber.org/zap"
	"time"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	// The client is a heavyweight object that should be created once per process.
	c, err := client.NewClient(client.Options{})
	if err != nil {
		panic(err)
	}
	defer func() { _ = c.CloseConnection() }()

	workflowOptions := client.StartWorkflowOptions{
		TaskList:                     "branch",
		ExecutionStartToCloseTimeout: time.Minute,
	}
	ctx := context.Background()
	we, err := c.ExecuteWorkflow(ctx, workflowOptions, branch.SampleBranchWorkflow, 10)
	if err != nil {
		logger.Fatal("Failure starting workflow", zap.Error(err))
	}
	logger.Info("Started workflow", zap.String("WorkflowID", we.GetID()), zap.String("RunID", we.GetRunID()))

	// Wait for workflow completion. This is rarely needed in real use cases
	// when workflows are potentially long running
	var result []string
	err = we.Get(ctx, &result)
	if err != nil {
		panic(err)
	}
	logger.Info("Started workflow", zap.String("WorkflowID", we.GetID()), zap.String("RunID", we.GetRunID()))
}
