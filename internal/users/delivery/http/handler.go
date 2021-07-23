package http

import (
	"fmt"
	"net/http"

	"github.com/IfuryI/ChatAPI/internal/logger"
	"github.com/IfuryI/ChatAPI/internal/models"
	"github.com/IfuryI/ChatAPI/internal/users"
	"github.com/gin-gonic/gin"
)

// Handler структура хендлера пользователя
type Handler struct {
	useCase users.UseCase
	Log     *logger.Logger
}

// NewHandler инициализация структуры хендлера пользователя
func NewHandler(useCase users.UseCase, Log *logger.Logger) *Handler {
	return &Handler{
		useCase: useCase,
		Log:     Log,
	}
}

// CreateUser создание пользователя
func (h *Handler) CreateUser(ctx *gin.Context) {
	userData := new(models.User)
	err := ctx.BindJSON(userData)
	if err != nil {
		msg := "Failed to bind request data " + err.Error()
		h.Log.LogWarning(ctx, "users", "CreateUser", msg)
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	if userData.Username == "" {
		err := fmt.Errorf("%s", "invalid value in user data")
		h.Log.LogWarning(ctx, "users", "CreateUser", err.Error())
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	err = h.useCase.CreateUser(userData)
	if err != nil {
		h.Log.LogError(ctx, "users", "CreateUser", err)
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	ctx.Status(http.StatusCreated) // 201
}
