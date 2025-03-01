package repositories

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/sonnyvictok/miniapp_taptoearn/internal/models"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByTelegramID(telegramID int64) (*models.User, error) {
	var user models.User
	err := r.db.Get(&user, "SELECT * FROM users WHERE telegram_id = $1", telegramID)
	return &user, err
}

func (r *UserRepository) Save(user *models.User) error {
	_, err := r.db.Exec("INSERT INTO users (telegram_id, username, score) VALUES ($1, $2, $3)",
		user.TelegramID, user.Username, user.Score)
	return err
}
func (r *UserRepository) GetScore(telegramID int64) (int, error) {
	var score int
	err := r.db.QueryRow("SELECT score FROM users WHERE telegram_id = $1", telegramID).Scan(&score)
	return score, err
}

func (r *UserRepository) UpdateScore(telegramID int64, score int) error {
	result, err := r.db.Exec("UPDATE users SET score = $1 WHERE telegram_id = $2", score, telegramID)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no user found with telegram_id: %d", telegramID)
	}

	return nil
}
