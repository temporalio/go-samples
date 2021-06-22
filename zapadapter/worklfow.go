package zapadapter

import (
	"context"
	"errors"
	"time"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/workflow"
)

func Workflow(ctx workflow.Context, name string) error {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	logger := workflow.GetLogger(ctx)
	logger.Info("Logging from workflow", "name", name)

	var result interface{}
	err := workflow.ExecuteActivity(ctx, Activity, name).Get(ctx, &result)
	if err != nil {
		logger.Error("Activity failed.", "Error", err)
		return err
	}

	err = workflow.ExecuteActivity(ctx, ActivityError).Get(ctx, &result)
	if err != nil {
		logger.Error("Activity failed.", "Error", err)
		return err
	}

	logger.Info("Workflow completed.")
	return nil
}

func Activity(ctx context.Context, name string) error {
	logger := activity.GetLogger(ctx)
	logger.Info("Executing Activity.", "name", name)
	logger.Debug("Debugging Activity.", "value", "important debug data")
	return nil
}

func ActivityError(ctx context.Context) error {
	logger := activity.GetLogger(ctx)
	logger.Warn("Ignore next error message. It is just for demo purpose.")
	logger.Error("Unable to execute ActivityError.", "error", errors.New("random error"))
	return nil
}
