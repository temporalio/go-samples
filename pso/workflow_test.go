package pso

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.temporal.io/temporal/activity"
	"go.temporal.io/temporal/encoded"
	"go.temporal.io/temporal/testsuite"
	"go.temporal.io/temporal/workflow"
)

func Test_Workflow(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()

	env.RegisterWorkflow(PSOChildWorkflow)

	env.RegisterActivityWithOptions(
		InitParticleActivity,
		activity.RegisterOptions{Name: InitParticleActivityName},
	)
	env.RegisterActivityWithOptions(
		UpdateParticleActivity,
		activity.RegisterOptions{Name: UpdateParticleActivityName},
	)

	var activityCalled []string

	var dataConverter = NewJSONDataConverter()
	env.SetDataConverter(dataConverter)

	// env.SetWorkflowTimeout(time.Minute * 5)
	// env.SetTestTimeout(time.Minute * 5)

	env.SetOnActivityStartedListener(func(activityInfo *activity.Info, ctx context.Context, args encoded.Values) {
		activityType := activityInfo.ActivityType.Name
		activityCalled = append(activityCalled, activityType)
		switch activityType {
		case "initParticleActivityName":
		case "updateParticleActivityName":
		default:
			panic("unexpected activity call")
		}
	})

	var childWorkflowID string
	env.SetOnChildWorkflowStartedListener(func(workflowInfo *workflow.Info, ctx workflow.Context, args encoded.Values) {
		childWorkflowID = workflowInfo.WorkflowExecution.ID
	})

	env.ExecuteWorkflow(PSOWorkflow, "sphere")

	require.True(t, env.IsWorkflowCompleted())
	queryAndVerify(t, env, "child", childWorkflowID)
	//queryAndVerify(t, env, "iteration", "???")
	// consider recreating a new test env on every iteration and calling execute workflow
	// with the arguments from the previous iteration (contained in ContinueAsNewError)
	require.Equal(t, "ContinueAsNew", env.GetWorkflowError().Error())
}

func queryAndVerify(t *testing.T, env *testsuite.TestWorkflowEnvironment, query string, expectedState string) {
	result, err := env.QueryWorkflow(query)
	require.NoError(t, err)
	var state string
	err = result.Get(&state)
	require.NoError(t, err)
	require.Equal(t, expectedState, state)
}
