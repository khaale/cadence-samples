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

type (
	// ChunkResult contains the result for this sample
	ChunkResult struct {
		NumberOfItemsInChunk int
		SumInChunk           int
	}
)

// This is registration process where you register all your workflows
// and activity function handlers.
func init() {
	workflow.Register(Workflow)
	activity.Register(chunkProcessingActivity)
}

// Workflow workflow decider
func Workflow(ctx workflow.Context, workerCount int) (ChunkResult, error) {
	chunkResultChannel := workflow.NewChannel(ctx)
	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Minute,
		StartToCloseTimeout:    time.Minute,
		HeartbeatTimeout:       time.Second * 20,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	for i := 1; i <= workerCount; i++ {
		chunkID := i
		workflow.Go(ctx, func(ctx workflow.Context) {
			var result ChunkResult
			err := workflow.ExecuteActivity(ctx, chunkProcessingActivity, chunkID).Get(ctx, &result)
			if err == nil {
				chunkResultChannel.Send(ctx, result)
			} else {
				chunkResultChannel.Send(ctx, err)
			}
		})
	}

	var totalItemCount, totalSum int
	for i := 1; i <= workerCount; i++ {
		var v interface{}
		chunkResultChannel.Receive(ctx, &v)
		switch r := v.(type) {
		case error:
		// failed to process this chunk
		// some proper error handling code here
		case ChunkResult:
			totalItemCount += r.NumberOfItemsInChunk
			totalSum += r.SumInChunk
		}
	}

	return ChunkResult{totalItemCount, totalSum}, nil
}

func chunkProcessingActivity(ctx context.Context, chunkID int) (result ChunkResult, err error) {
	// some fake processing logic here
	numberOfItemsInChunk := chunkID
	sumInChunk := chunkID * chunkID

	activity.GetLogger(ctx).Debug("Chunck processed", zap.Int("chunkID", chunkID))
	return ChunkResult{numberOfItemsInChunk, sumInChunk}, nil
}
