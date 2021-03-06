package usecase

import (
	"errors"

	"github.com/IfuryI/ChatAPI/internal/chats"
	"github.com/IfuryI/ChatAPI/internal/models"
	"github.com/IfuryI/ChatAPI/internal/users"
)

// ChatUseCase структура usecase чата
type ChatUseCase struct {
	chatRepository chats.Repository
	userRepository users.UserRepository
}

// NewChatUseCase инициализация структуры usecase чата
func NewChatUseCase(chatRepo chats.Repository, userRepo users.UserRepository) *ChatUseCase {
	return &ChatUseCase{
		chatRepository: chatRepo,
		userRepository: userRepo,
	}
}

// CreateChat создание чата
func (chatUC *ChatUseCase) CreateChat(chat *models.Chat) (int, error) {
	_, err := chatUC.chatRepository.GetChatByName(chat.Name)
	if err == nil {
		return 0, errors.New("chat already exists")
	}

	return chatUC.chatRepository.CreateChat(chat)
}

// AddMessageToChat создание сообщения в чате
func (chatUC *ChatUseCase) AddMessageToChat(message *models.Message) (int, error) {
	err := chatUC.chatRepository.IsUserInChat(message.UserID, message.ChatID)
	if err != nil {
		return 0, errors.New("user is not in chat")
	}

	return chatUC.chatRepository.AddMessageToChat(message)
}

// GetAllUserChats получение всех чатов пользователя
func (chatUC *ChatUseCase) GetAllUserChats(userID string) ([]*models.Chat, error) {
	_, err := chatUC.userRepository.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("user does not exists")
	}

	return chatUC.chatRepository.GetAllUserChats(userID)
}

// GetAllChatMessages получение всех сообщений чата
func (chatUC *ChatUseCase) GetAllChatMessages(chatID string) ([]*models.Message, error) {
	_, err := chatUC.chatRepository.GetChatByID(chatID)
	if err != nil {
		return nil, errors.New("chat does not exists")
	}

	return chatUC.chatRepository.GetAllChatMessages(chatID)
}
