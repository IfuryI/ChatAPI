package chats

import "github.com/IfuryI/ChatAPI/internal/models"

// UseCase интерфейс usecase для работы с логикой запросов чата
//go:generate mockgen -destination=mocks/usecase.go -package=mocks . UseCase
type UseCase interface {
	CreateChat(chat *models.Chat) (int, error)

	AddMessageToChat(message *models.Message) (int, error)

	GetAllUserChats(userID string) ([]*models.Chat, error)

	GetAllChatMessages(chatID string) ([]*models.Message, error)
}
