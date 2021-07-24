package http

import (
	"net/http"
	"strconv"

	"github.com/IfuryI/ChatAPI/internal/chats"
	"github.com/IfuryI/ChatAPI/internal/logger"
	"github.com/IfuryI/ChatAPI/internal/models"
	"github.com/IfuryI/ChatAPI/internal/users"
	"github.com/gin-gonic/gin"
)

// Handler структура хендлера
type Handler struct {
	chatUC  chats.UseCase
	usersUC users.UseCase
	Log     *logger.Logger
}

// NewHandler инициализация структуры хендлера чата
func NewHandler(chatUC chats.UseCase, usersUC users.UseCase, Log *logger.Logger) *Handler {
	return &Handler{
		chatUC:  chatUC,
		usersUC: usersUC,
		Log:     Log,
	}
}

// userRequest структура для анмаршалинга данных из POST запроса
type userRequest struct {
	UserID string `json:"user"`
}

// chatRequest структура для анмаршалинга данных из POST запроса
type chatRequest struct {
	ChatID string `json:"chat"`
}

// idResponse структура для маршалинга данных (для отправки ID)
type idResponse struct {
	ID string `json:"id"`
}

// CreateChat создание чата
func (h *Handler) CreateChat(ctx *gin.Context) {
	chatData := new(models.Chat)
	err := ctx.BindJSON(chatData)
	if err != nil {
		msg := "Failed to bind request data" + err.Error()
		h.Log.LogWarning(ctx, "chats", "CreateChat", msg)
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	chatID, err := h.chatUC.CreateChat(chatData)
	if err != nil {
		h.Log.LogWarning(ctx, "chats", "CreateChat", err.Error())
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	ctx.JSON(http.StatusCreated, idResponse{ID: strconv.Itoa(chatID)}) // 201
}

// AddMessageToChat добавление сообщения в чат
func (h *Handler) AddMessageToChat(ctx *gin.Context) {
	messageData := new(models.Message)
	err := ctx.BindJSON(messageData)
	if err != nil {
		msg := "Failed to bind request data" + err.Error()
		h.Log.LogWarning(ctx, "chats", "AddMessageToChat", msg)
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	messageID, err := h.chatUC.AddMessageToChat(messageData)
	if err != nil {
		h.Log.LogWarning(ctx, "chats", "AddMessageToChat", err.Error())
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	ctx.JSON(http.StatusCreated, idResponse{ID: strconv.Itoa(messageID)}) // 201
}

// GetAllUserChats получить все чаты, в которых учавствует пользователь
func (h *Handler) GetAllUserChats(ctx *gin.Context) {
	userData := new(userRequest)
	err := ctx.BindJSON(userData)
	if err != nil {
		msg := "Failed to bind request data" + err.Error()
		h.Log.LogWarning(ctx, "chats", "GetAllUserChats", msg)
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	userChats, err := h.chatUC.GetAllUserChats(userData.UserID)
	if err != nil {
		h.Log.LogWarning(ctx, "chats", "GetAllUserChats", err.Error())
		ctx.AbortWithStatus(http.StatusNotFound) // 404
		return
	}

	ctx.JSON(http.StatusOK, userChats)
}

// GetAllChatMessages получить все сообщения из чата
func (h *Handler) GetAllChatMessages(ctx *gin.Context) {
	chatData := new(chatRequest)
	err := ctx.BindJSON(chatData)
	if err != nil {
		msg := "Failed to bind request data" + err.Error()
		h.Log.LogWarning(ctx, "chats", "GetAllChatMessages", msg)
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	userChats, err := h.chatUC.GetAllChatMessages(chatData.ChatID)
	if err != nil {
		h.Log.LogWarning(ctx, "chats", "GetAllChatMessages", err.Error())
		ctx.AbortWithStatus(http.StatusNotFound) // 404
		return
	}

	ctx.JSON(http.StatusOK, userChats)
}
