package users

import (
	"github.com/IfuryI/ChatAPI/internal/models"
)

// UseCase интерфейс usecase для работы с логикой пользовательских запросов
//go:generate mockgen -destination=mocks/usecase.go -package=mocks . UseCase
type UseCase interface {
	CreateUser(user *models.User) error
}
