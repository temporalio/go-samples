package helloworld

import (
	"context"
	"time"

	"go.temporal.io/temporal/activity"
	"go.temporal.io/temporal/workflow"
	"go.uber.org/zap"
)

// Workflow is a Hello World workflow definition.
func Workflow(ctx workflow.Context, name string) (string, error) {
	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Minute,
		StartToCloseTimeout:    time.Minute,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	logger := workflow.GetLogger(ctx)
	logger.Info("HelloWorld workflow started", zap.String("name", name))

	var result string
	err := workflow.ExecuteActivity(ctx, Activity, name).Get(ctx, &result)
	if err != nil {
		logger.Error("Activity failed.", zap.Error(err))
		return "", err
	}

	logger.Info("HelloWorld workflow completed.", zap.String("result", result))

	return result, nil
}

func Activity(ctx context.Context, name string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Activity", zap.String("name", name))
	return "Hello " + name + "!", nil
}
