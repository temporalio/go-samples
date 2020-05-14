package main

import (
	"context"
	"time"

	"github.com/pborman/uuid"
	cw "github.com/temporalio/temporal-go-samples/child-workflow-continue-as-new"
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

	// This workflow ID can be user business logic identifier as well.
	workflowID := "parent-workflow_" + uuid.New()
	workflowOptions := client.StartWorkflowOptions{
		ID:                           workflowID,
		TaskList:                     "child-workflow-continue-as-new",
		ExecutionStartToCloseTimeout: time.Minute,
	}

	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, cw.SampleParentWorkflow)
	if err != nil {
		logger.Fatal("Unable to execute workflow", zap.Error(err))
	}
	logger.Info("Started workflow", zap.String("WorkflowID", we.GetID()), zap.String("RunID", we.GetRunID()))
}
