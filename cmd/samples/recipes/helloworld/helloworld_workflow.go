package main

import (
	"context"
	"time"

	"go.uber.org/cadence/activity"
	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"
)

/**
 * This is the hello world workflow sample.
 */

// ApplicationName is the task list for this sample
const ApplicationName = "helloWorldGroup"

// This is registration process where you register all your workflows
// and activity function handlers.
func init() {
	workflow.Register(Workflow)
	activity.Register(helloworldActivity)
}

// Workflow workflow decider
func Workflow(ctx workflow.Context, name string) error {
	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Minute,
		StartToCloseTimeout:    time.Minute,
		HeartbeatTimeout:       time.Second * 20,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	logger := workflow.GetLogger(ctx)
	logger.Info("hello world workflow started")

	lao := workflow.LocalActivityOptions{
		ScheduleToCloseTimeout: time.Minute,
	}
	lctx := workflow.WithLocalActivityOptions(ctx, lao)

	var beforeResult string
	err := workflow.ExecuteLocalActivity(lctx, localBefore, "Before Param").Get(lctx, &beforeResult)
	if err != nil {
		logger.Error("Activity failed.", zap.Error(err))
		return err
	}

	var helloworldResult string
	err = workflow.ExecuteActivity(ctx, helloworldActivity, name).Get(ctx, &helloworldResult)
	if err != nil {
		logger.Error("Activity failed.", zap.Error(err))
		return err
	}

	var afterResult string
	err = workflow.ExecuteLocalActivity(lctx, localAfter, "After Param").Get(lctx, &afterResult)
	if err != nil {
		logger.Error("Activity failed.", zap.Error(err))
		return err
	}

	logger.Info("Workflow completed.", zap.String("Result", helloworldResult))

	return nil
}

func localBefore(ctx context.Context, name string) (string, error) {
	time.Sleep(500 * time.Millisecond)
	return "Local Before", nil
}

func localAfter(ctx context.Context, name string) (string, error) {
	time.Sleep(500 * time.Millisecond)
	return "Local After", nil
}

func helloworldActivity(ctx context.Context, name string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("helloworld activity started")
	time.Sleep(1000 * time.Millisecond)
	return "Hello " + name + "!", nil
}
