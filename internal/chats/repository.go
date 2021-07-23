package chats

import "github.com/IfuryI/ChatAPI/internal/models"

// Repository интерфейс репозитория для работы с чатом
//go:generate mockgen -destination=mocks/repository.go -package=mocks . Repository
type Repository interface {
	CreateChat(chat *models.Chat) (int, error)

	GetChatByName(chatName string) (*models.Chat, error)

	GetChatByID(chatID string) (*models.Chat, error)

	AddMessageToChat(message *models.Message) (int, error)

	IsUserInChat(userID string, chatID string) error

	GetAllUserChats(userID string) ([]*models.Chat, error)

	GetAllChatMessages(chatID string) ([]*models.Message, error)
}
