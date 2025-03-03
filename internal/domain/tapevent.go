package domain

import "time"

type TapEventPayload struct {
	TelegramID int64     `json:"telegram_id"`
	Timestamp  time.Time `json:"timestamp"`
	Score      int       `json:"score"`
}

type UserTelegramPayload struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
