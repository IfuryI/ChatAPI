package http

import (
	"github.com/IfuryI/ChatAPI/internal/chats"
	"github.com/IfuryI/ChatAPI/internal/logger"
	"github.com/IfuryI/ChatAPI/internal/users"
	"github.com/gin-gonic/gin"
)

// RegisterHTTPEndpoints регистрация хендлеров
func RegisterHTTPEndpoints(router *gin.Engine, chatUC chats.UseCase, usersUC users.UseCase, Log *logger.Logger) {
	handler := NewHandler(chatUC, usersUC, Log)

	router.POST("/chats/add", handler.CreateChat)
	router.POST("/chats/get", handler.GetAllUserChats)

	router.POST("/messages/add", handler.AddMessageToChat)
	router.POST("/messages/get", handler.GetAllChatMessages)
}
