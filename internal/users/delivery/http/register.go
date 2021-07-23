package http

import (
	"github.com/IfuryI/ChatAPI/internal/logger"
	"github.com/IfuryI/ChatAPI/internal/users"
	"github.com/gin-gonic/gin"
)

// RegisterHTTPEndpoints регистрация хендлеров
func RegisterHTTPEndpoints(router *gin.Engine, usersUC users.UseCase, Log *logger.Logger) {
	handler := NewHandler(usersUC, Log)

	router.POST("/users/add", handler.CreateUser)
}
