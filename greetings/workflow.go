package greetings

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

// GreetingSample workflow definition.
// This greetings sample workflow executes 3 activities in sequential.
// It gets greeting and name from 2 different activities,
// and then pass greeting and name as input to a 3rd activity to generate final greetings.
func GreetingSample(ctx workflow.Context) (string, error) {
	logger := workflow.GetLogger(ctx)

	ao := workflow.ActivityOptions{
		StartToCloseTimeout:    time.Minute,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	var a *Activities

	var greetResult string
	err := workflow.ExecuteActivity(ctx, a.GetGreeting).Get(ctx, &greetResult)
	if err != nil {
		logger.Error("Get greeting failed.", "Error", err)
		return "", err
	}

	// Get Name.
	var nameResult string
	err = workflow.ExecuteActivity(ctx, a.GetName).Get(ctx, &nameResult)
	if err != nil {
		logger.Error("Get name failed.", "Error", err)
		return "", err
	}

	// Say Greeting.
	var sayResult string
	err = workflow.ExecuteActivity(ctx, a.SayGreeting, greetResult, nameResult).Get(ctx, &sayResult)
	if err != nil {
		logger.Error("Marshalling failed with error.", "Error", err)
		return "", err
	}

	logger.Info("GreetingSample completed.", "Result", sayResult)
	return sayResult, nil
}
