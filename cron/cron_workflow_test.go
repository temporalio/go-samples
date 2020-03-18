package cron

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.temporal.io/temporal/activity"
	"go.temporal.io/temporal/encoded"
	"go.temporal.io/temporal/testsuite"
	"go.temporal.io/temporal/workflow"
)

type UnitTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}

func (s *UnitTestSuite) Test_CronWorkflow() {
	testWorkflow := func(ctx workflow.Context) error {
		ctx1 := workflow.WithChildOptions(ctx, workflow.ChildWorkflowOptions{
			ExecutionStartToCloseTimeout: time.Minute * 10,
			CronSchedule:                 "0 * * * *", // hourly
		})

		cronFuture := workflow.ExecuteChildWorkflow(ctx1, SampleCronWorkflow) // cron never stop so this future won't return

		// wait 2 hours for the cron (cron will execute 3 times)
		_ = workflow.Sleep(ctx, time.Hour*2)
		s.False(cronFuture.IsReady())
		return nil
	}

	env := s.NewTestWorkflowEnvironment()

	env.RegisterWorkflow(testWorkflow)
	env.RegisterWorkflow(SampleCronWorkflow)
	env.RegisterActivity(SampleCronActivity)

	env.OnActivity(SampleCronActivity, mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(3)

	var startTimeList, endTimeList []time.Time
	env.SetOnActivityStartedListener(func(activityInfo *activity.Info, ctx context.Context, args encoded.Values) {
		var startTime, endTime time.Time
		err := args.Get(&startTime, &endTime)
		s.NoError(err)

		startTimeList = append(startTimeList, startTime)
		endTimeList = append(endTimeList, endTime)
	})

	startTime, _ := time.Parse(time.RFC3339, "2018-12-20T16:30:00-80:00")
	env.SetStartTime(startTime)

	env.ExecuteWorkflow(testWorkflow)

	s.True(env.IsWorkflowCompleted())
	err := env.GetWorkflowError()
	s.NoError(err)
	env.AssertExpectations(s.T())

	s.Len(startTimeList, 3)
	s.Equal(time.Time{}, startTimeList[0])
	s.Equal(startTime, endTimeList[0])

	s.Equal(startTime, startTimeList[1])
	s.Equal(startTime.Add(time.Minute*30), endTimeList[1])

	s.Equal(startTime.Add(time.Minute*30), startTimeList[2])
	s.Equal(startTime.Add(time.Minute*90), endTimeList[2])
}
