package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const TaskTapToEarn = "task:tap_to_earn"

type PayloadHandleTapToEarn struct {
	TelegramID int64 `json:"telegram_id"`
	Score      int   `json:"score"`
}

func (distributor *RedisTaskDistributor) DistributeTaskHandleTapToEarn(ctx context.Context, payload *PayloadHandleTapToEarn, opts ...asynq.Option) error {
	json, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Create Task
	task := asynq.NewTask(TaskTapToEarn, json, opts...)

	// Put task on  Queues
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Str("queue", info.Queue).Int("max_retry", info.MaxRetry).Msg("enqueued task")
	return nil
}

func (processor *RedisTaskProcessor) ProcessTaskTapToEarn(ctx context.Context, task *asynq.Task) error {
	var payload PayloadHandleTapToEarn
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}
	log.Info().Interface("payload", payload).Msg("Processing TapToEarn task")

	score, err := processor.userService.ClickToEarn(payload.TelegramID)
	if err != nil {
		return fmt.Errorf("failed to process TapToEarn: %w", err)
	}

	log.Info().Msgf("✅ Updated score for user %d: %d points", payload.TelegramID, score)

	return nil
}
