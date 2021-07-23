package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/IfuryI/ChatAPI/internal/users"
	usersHttp "github.com/IfuryI/ChatAPI/internal/users/delivery/http"
	usersDBStorage "github.com/IfuryI/ChatAPI/internal/users/repository/dbstorage"
	usersUseCase "github.com/IfuryI/ChatAPI/internal/users/usecase"

	"github.com/IfuryI/ChatAPI/internal/chats"
	chatsHttp "github.com/IfuryI/ChatAPI/internal/chats/delivery/http"
	chatsDBStorage "github.com/IfuryI/ChatAPI/internal/chats/repository/dbstorage"
	chatsUseCase "github.com/IfuryI/ChatAPI/internal/chats/usecase"

	"github.com/IfuryI/ChatAPI/internal/logger"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
)

// App структура главного приложения
type App struct {
	server  *http.Server
	usersUC users.UseCase
	chatsUC chats.UseCase
	logger  *logger.Logger
}

// NewApp инициализация приложения
func NewApp() *App {
	accessLogger := logger.NewAccessLogger()

	connStr, connected := os.LookupEnv("DB_CONNECT")
	if !connected {
		fmt.Println(os.Getwd())
		log.Fatal("Failed to read DB connection data")
	}

	dbpool, err := pgxpool.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	usersRepo := usersDBStorage.NewUserRepository(dbpool)
	chatsRepo := chatsDBStorage.NewChatRepository(dbpool)

	usersUC := usersUseCase.NewUsersUseCase(usersRepo)
	chatsUC := chatsUseCase.NewChatUseCase(chatsRepo, usersRepo)

	return &App{
		usersUC: usersUC,
		chatsUC: chatsUC,
		logger:  accessLogger,
	}
}

// Run запуск приложения
func (app *App) Run(port string) error {
	router := gin.Default()

	router.Use(gin.Recovery())

	// Для версионирования апи можно добавить префикс (однако этого не было в тз, поэтому закомментировано):
	// api := router.Group("/api")
	// v1 := api.Group("/v1")

	usersHttp.RegisterHTTPEndpoints(router, app.usersUC, app.logger)
	chatsHttp.RegisterHTTPEndpoints(router, app.chatsUC, app.usersUC, app.logger)

	app.server = &http.Server{
		Addr:           ":" + port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		err := app.server.ListenAndServe()
		if err != nil {
			log.Fatal("Failed to listen and serve: ", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return app.server.Shutdown(ctx)
}

// init функция инициализации .env файла
func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}
