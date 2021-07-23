package dbstorage

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	"github.com/IfuryI/ChatAPI/internal/models"
	"github.com/jackc/pgconn"
	pgx "github.com/jackc/pgx/v4"
)

// PgxPoolIface Интерфейс для драйвера БД (использую его для mock тестирования)
type PgxPoolIface interface {
	Begin(context.Context) (pgx.Tx, error)
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	Ping(context.Context) error
}

// ChatRepository структура репозитория чата
type ChatRepository struct {
	db PgxPoolIface
}

// NewChatRepository инициализация структуры репозитория чата
func NewChatRepository(database PgxPoolIface) *ChatRepository {
	return &ChatRepository{
		db: database,
	}
}

// CreateChat создание чата
func (storage *ChatRepository) CreateChat(chat *models.Chat) (int, error) {
	sqlStatement := `
        INSERT INTO mdb.chats (name)
        VALUES ($1)
        RETURNING "id";
    `
	var newChatID int

	err := storage.db.
		QueryRow(context.Background(), sqlStatement, chat.Name).
		Scan(&newChatID)

	if err != nil {
		return 0, err
	}

	chat.ID = strconv.Itoa(newChatID)

	for i := 0; i < len(chat.ChatUsers); i++ {
		intID, err := strconv.Atoi(chat.ChatUsers[i])
		if err != nil {
			storage.deleteChat(newChatID)
			return 0, err
		}

		err = storage.addUserToChat(intID, newChatID)
		if err != nil {
			storage.deleteChat(newChatID)
			return 0, err
		}
	}

	return newChatID, err
}

// AddMessageToChat добавление сообщения в чат
func (storage *ChatRepository) AddMessageToChat(message *models.Message) (int, error) {
	sqlStatement := `
        INSERT INTO mdb.messages (chat_id, user_id, content)
        VALUES ($1, $2, $3)
        RETURNING "id";
    `
	var newMessageID int

	err := storage.db.
		QueryRow(context.Background(), sqlStatement, message.ChatID, message.UserID, message.Text).
		Scan(&newMessageID)

	if err != nil {
		return 0, err
	}

	message.ID = strconv.Itoa(newMessageID)

	return newMessageID, nil
}

// IsUserInChat проверка находится ли пользователь в чате или нет
func (storage *ChatRepository) IsUserInChat(userID string, chatID string) error {
	sqlStatement := `
        SELECT user_id, chat_id
        FROM mdb.users_in_chats
        WHERE user_id = $1 AND chat_id = $2
    `

	intUserID, err := strconv.Atoi(userID)
	if err != nil {
		return err
	}

	intChatID, err := strconv.Atoi(chatID)
	if err != nil {
		return err
	}

	err = storage.db.
		QueryRow(context.Background(), sqlStatement, intUserID, intChatID).
		Scan(&intUserID, &intChatID)
	if err != nil {
		return err
	}

	return nil
}

// GetChatByName получение чата по имени (использую для проверки существования чата)
func (storage *ChatRepository) GetChatByName(chatName string) (*models.Chat, error) {
	sqlStatement := `
        SELECT id, created_at
        FROM mdb.chats
        WHERE name = $1
    `
	chatData := &models.Chat{}

	var intChatID int

	err := storage.db.
		QueryRow(context.Background(), sqlStatement, chatName).
		Scan(&intChatID, &chatData.CreatedAt)

	if err != nil {
		return nil, err
	}

	chatData.ID = strconv.Itoa(intChatID)
	chatData.Name = chatName

	return chatData, nil
}

// GetChatByID получение чата по ID (использую для проверки существования чата)
func (storage *ChatRepository) GetChatByID(chatID string) (*models.Chat, error) {
	sqlStatement := `
        SELECT name, created_at
        FROM mdb.chats
        WHERE id = $1
    `
	chatData := &models.Chat{}

	intChatID, err := strconv.Atoi(chatID)
	if err != nil {
		return nil, err
	}
	err = storage.db.
		QueryRow(context.Background(), sqlStatement, intChatID).
		Scan(&chatData.Name, &chatData.CreatedAt)

	if err != nil {
		return nil, err
	}

	chatData.ID = chatID

	return chatData, nil
}

// addUserToChat добавление нового пользователя в чат
func (storage *ChatRepository) addUserToChat(userID int, chatID int) error {
	sqlStatement := `
        INSERT INTO mdb.users_in_chats (user_id, chat_id)
        VALUES ($1, $2);
    `
	_, err := storage.db.Exec(context.Background(), sqlStatement, userID, chatID)
	if err != nil {
		return err
	}

	return nil
}

// deleteChat удаление чата
func (storage *ChatRepository) deleteChat(chatID int) error {
	sqlStatement := `
        DELETE FROM mdb.chats
        WHERE id = $1;
    `

	_, err := storage.db.Exec(context.Background(), sqlStatement, chatID)
	if err != nil {
		return err
	}

	return nil
}

// GetAllUserChats получение всех чатов, в которых учавствует пользователь
// учитывая, что чаты отсортированы по времени последнего сообщения в нем
func (storage *ChatRepository) GetAllUserChats(userID string) ([]*models.Chat, error) {
	sqlStatement := `
        SELECT uic.chat_id, ch.name, array_agg(DISTINCT uic2.user_id) as user_in_chat, ch.created_at
        FROM mdb.users_in_chats uic JOIN mdb.chats ch ON uic.chat_id = ch.id AND uic.user_id = $1
		LEFT JOIN mdb.messages msg ON uic.chat_id = msg.chat_id
		JOIN mdb.users_in_chats uic2 ON uic.chat_id = uic2.chat_id
		GROUP BY uic.chat_id, ch.name, ch.created_at
		ORDER BY MAX(msg.created_at) DESC
    `

	var userChats []*models.Chat

	intUserID, err := strconv.Atoi(userID)
	if err != nil {
		return nil, err
	}

	rows, err := storage.db.Query(context.Background(), sqlStatement, intUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		chat := &models.Chat{}
		var newID int
		var newUserIDArr []int

		err = rows.Scan(&newID, &chat.Name, &newUserIDArr, &chat.CreatedAt)
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}

		chat.ID = strconv.Itoa(newID)
		for i := 0; i < len(newUserIDArr); i++ {
			chat.ChatUsers = append(chat.ChatUsers, strconv.Itoa(newUserIDArr[i]))
		}

		userChats = append(userChats, chat)
	}

	return userChats, nil
}

// GetAllChatMessages получение всех сообщений чата
func (storage *ChatRepository) GetAllChatMessages(chatID string) ([]*models.Message, error) {
	sqlStatement := `
        SELECT id, user_id, content, created_at
        FROM mdb.messages msg
        WHERE msg.chat_id = $1
		ORDER BY msg.created_at DESC
    `

	var chatMessages []*models.Message

	intChatID, err := strconv.Atoi(chatID)
	if err != nil {
		return nil, err
	}

	rows, err := storage.db.Query(context.Background(), sqlStatement, intChatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		message := &models.Message{}
		var newID int
		var newUserID int

		err = rows.Scan(&newID, &newUserID, &message.Text, &message.CreatedAt)
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}

		message.ID = strconv.Itoa(newID)
		message.UserID = strconv.Itoa(newUserID)
		message.ChatID = chatID

		chatMessages = append(chatMessages, message)
	}

	if len(chatMessages) == 0 {
		return nil, errors.New("no messages in chat")
	}

	return chatMessages, nil
}
