package dbstorage

import (
	"context"
	"strconv"

	"github.com/IfuryI/ChatAPI/internal/models"
	"github.com/jackc/pgconn"
	pgx "github.com/jackc/pgx/v4"
)

// PgxPoolIface Интерфейс для драйвера БД (использую его для возможности mock тестирования)
type PgxPoolIface interface {
	Begin(context.Context) (pgx.Tx, error)
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	Ping(context.Context) error
}

// UserRepository структура репозитория пользователя
type UserRepository struct {
	db PgxPoolIface
}

// NewUserRepository инициализация репозитория пользователя
func NewUserRepository(database PgxPoolIface) *UserRepository {
	return &UserRepository{
		db: database,
	}
}

// CreateUser создание пользователя
func (storage *UserRepository) CreateUser(user *models.User) (int, error) {
	sqlStatement := `
        INSERT INTO mdb.users (username)
        VALUES ($1)
		RETURNING id;
    `
	var newUserID int

	err := storage.db.QueryRow(context.Background(), sqlStatement, user.Username).Scan(&newUserID)

	if err != nil {
		return 0, err
	}

	return newUserID, nil
}

// GetUserByUsername получить информацию о пользователе по имени
func (storage *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User

	sqlStatement := `
        SELECT id, username, created_at
        FROM mdb.users
        WHERE username = $1
    `

	err := storage.db.
		QueryRow(context.Background(), sqlStatement, username).
		Scan(&user.ID, &user.Username, &user.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUserByID получить информацию о пользователе по ID
func (storage *UserRepository) GetUserByID(userID string) (*models.User, error) {
	var user models.User

	sqlStatement := `
        SELECT username, created_at
        FROM mdb.users
        WHERE id = $1
    `

	intUserID, err := strconv.Atoi(userID)
	if err != nil {
		return nil, err
	}

	err = storage.db.
		QueryRow(context.Background(), sqlStatement, intUserID).
		Scan(&user.Username, &user.CreatedAt)

	if err != nil {
		return nil, err
	}

	user.ID = strconv.Itoa(intUserID)

	return &user, nil
}
