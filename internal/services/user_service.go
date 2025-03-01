package services

import (
	"database/sql"

	"github.com/hibiken/asynq"
	"github.com/sonnyvictok/miniapp_taptoearn/internal/models"
	"github.com/sonnyvictok/miniapp_taptoearn/internal/repositories"
)

type UserService struct {
	userRepo    *repositories.UserRepository
	asynqClient *asynq.Client
}

func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}
func (s *UserService) GetUserByTelegramID(telegramID int64) (*models.User, error) {
	return s.userRepo.FindByTelegramID(telegramID)
}
func (s *UserService) GetOrCreateUser(telegramID int64, username string) (*models.User, error) {
	user, err := s.userRepo.FindByTelegramID(telegramID)
	if err != nil {
		if err == sql.ErrNoRows {
			user = &models.User{
				TelegramID: telegramID,
				Username:   username,
				Score:      0,
			}
			err = s.userRepo.Save(user)
			if err != nil {
				return nil, err
			}
			return user, nil
		}
		return nil, err
	}

	return user, nil
}

func (s *UserService) ClickToEarn(telegramID int64) (int, error) {
	user, err := s.userRepo.FindByTelegramID(telegramID)
	if err != nil {
		return 0, err
	}

	newScore := user.Score + 1
	err = s.userRepo.UpdateScore(telegramID, newScore)
	if err != nil {
		return 0, err
	}

	// Gửi thông điệp đến Redis channel
	// ctx := context.Background()
	// channel := "user_click_to_earn"
	// message := fmt.Sprintf(`{"telegram_id": %d, "new_score": %d}`, telegramID, newScore)

	// Publish message to Redis channel
	// if err := s.redisClient.Publish(ctx, channel, message).Err(); err != nil {
	// 	return 0, fmt.Errorf("failed to publish message to Redis: %w", err)
	// }

	return newScore, nil
}

// func (s *UserService) SubscribeToClickToEarn(ctx context.Context) {
// 	channel := "user_click_to_earn"
// 	pubsub := s.redisClient.Subscribe(ctx, channel)

// 	defer pubsub.Close()

// 	ch := pubsub.Channel()
// 	for msg := range ch {
// 		fmt.Printf("Received message from %s: %s\n", channel, msg.Payload)
// 	}
// }
