package users

import (
	"github.com/IfuryI/ChatAPI/internal/models"
)

// UserRepository интерфейс репозитория для работы с пользователями
//go:generate mockgen -destination=mocks/repository.go -package=mocks . UserRepository
type UserRepository interface {
	CreateUser(user *models.User) error

	GetUserByUsername(username string) (*models.User, error)

	GetUserByID(userID string) (*models.User, error)
}
