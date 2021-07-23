package usecase

import (
	"errors"

	"github.com/IfuryI/ChatAPI/internal/models"
	"github.com/IfuryI/ChatAPI/internal/users"
)

// UsersUseCase структура usecase пользователя
type UsersUseCase struct {
	userRepository users.UserRepository
}

// NewUsersUseCase инициализация usecase пользователя
func NewUsersUseCase(repo users.UserRepository) *UsersUseCase {
	return &UsersUseCase{
		userRepository: repo,
	}
}

// CreateUser создание пользователя
func (usersUC *UsersUseCase) CreateUser(user *models.User) error {
	_, err := usersUC.userRepository.GetUserByUsername(user.Username) // Проверка существует ли пользователь
	if err == nil {
		return errors.New("user already exists")
	}

	return usersUC.userRepository.CreateUser(user)
}
