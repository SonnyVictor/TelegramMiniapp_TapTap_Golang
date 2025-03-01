package domain

import "time"

type TapEventPayload struct {
	TelegramID int64     `json:"telegram_id"`
	Timestamp  time.Time `json:"timestamp"`
	Score      int       `json:"score"`
}
