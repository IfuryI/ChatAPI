package usecase

import (
	"errors"
	"testing"

	"github.com/IfuryI/ChatAPI/internal/chats/mocks"
	"github.com/IfuryI/ChatAPI/internal/models"
	userMocks "github.com/IfuryI/ChatAPI/internal/users/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestReviewsUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepository(ctrl)
	usersRepo := userMocks.NewMockUserRepository(ctrl)
	uc := NewChatUseCase(repo, usersRepo)

	userIDStr := "1"
	chatIDStr := "1"

	chat := &models.Chat{
		ID:        "1",
		Name:      "Лучший чат",
		ChatUsers: []string{"1", "2", "3"},
	}

	message := &models.Message{
		ID:     "1",
		UserID: "1",
		ChatID: "2",
		Text:   "Очень хочу на стажировку в Авито",
	}

	t.Run("CreateChat", func(t *testing.T) {
		repo.EXPECT().GetChatByName(chat.Name).Return(nil, errors.New("cant get chat"))
		repo.EXPECT().CreateChat(chat).Return(1, nil)

		chatID, err := uc.CreateChat(chat)
		assert.NoError(t, err)
		assert.Equal(t, 1, chatID)
	})

	t.Run("AddMessageToChat", func(t *testing.T) {
		repo.EXPECT().IsUserInChat(message.UserID, message.ChatID).Return(nil)
		repo.EXPECT().AddMessageToChat(message).Return(1, nil)

		msgID, err := uc.AddMessageToChat(message)
		assert.NoError(t, err)
		assert.Equal(t, 1, msgID)
	})

	t.Run("GetAllUserChats", func(t *testing.T) {
		usersRepo.EXPECT().GetUserByID(userIDStr).Return(nil, nil)
		repo.EXPECT().GetAllUserChats(userIDStr).Return([]*models.Chat{}, nil)

		_, err := uc.GetAllUserChats(userIDStr)
		assert.NoError(t, err)
	})

	t.Run("GetAllChatMessages", func(t *testing.T) {
		repo.EXPECT().GetChatByID(chatIDStr).Return(nil, nil)
		repo.EXPECT().GetAllChatMessages(chatIDStr).Return([]*models.Message{}, nil)

		_, err := uc.GetAllChatMessages(chatIDStr)
		assert.NoError(t, err)
	})
}
