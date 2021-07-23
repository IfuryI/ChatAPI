package dbstorage

import (
	"context"
	"strconv"
	"testing"

	"github.com/IfuryI/ChatAPI/internal/models"

	"github.com/pashagolub/pgxmock"
)

func TestCreateRating(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	chat := &models.Chat{
		ID:        "1",
		Name:      "Лучший чат",
		ChatUsers: []string{"1", "2"},
	}

	rows := pgxmock.NewRows([]string{"id"}).AddRow(1)

	chatRepo := NewChatRepository(mock)

	mock.ExpectQuery("INSERT INTO mdb.chats").WithArgs(chat.Name).WillReturnRows(rows)

	intChatID, _ := strconv.Atoi(chat.ID)
	intArg0, _ := strconv.Atoi(chat.ChatUsers[0])
	intArg1, _ := strconv.Atoi(chat.ChatUsers[1])
	mock.ExpectExec("INSERT INTO mdb.users_in_chats").WithArgs(intArg0, intChatID).WillReturnResult(pgxmock.NewResult("INSERT", 1))
	mock.ExpectExec("INSERT INTO mdb.users_in_chats").WithArgs(intArg1, intChatID).WillReturnResult(pgxmock.NewResult("INSERT", 1))

	if _, err = chatRepo.CreateChat(chat); err != nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCreateRatingFail(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	chat := &models.Chat{
		ID:        "1",
		Name:      "Лучший чат",
		ChatUsers: []string{"1", "qwe"},
	}

	rows := pgxmock.NewRows([]string{"id"}).AddRow(1)

	chatRepo := NewChatRepository(mock)

	mock.ExpectQuery("INSERT INTO mdb.chats").WithArgs(chat.Name).WillReturnRows(rows)

	intChatID, _ := strconv.Atoi(chat.ID)
	intArg0, _ := strconv.Atoi(chat.ChatUsers[0])

	mock.ExpectExec("INSERT INTO mdb.users_in_chats").WithArgs(intArg0, intChatID).WillReturnResult(pgxmock.NewResult("INSERT", 1))
	mock.ExpectExec("INSERT INTO mdb.users_in_chats").WithArgs(chat.ChatUsers[1], intChatID).WillReturnResult(pgxmock.NewResult("INSERT", 1))

	if _, err = chatRepo.CreateChat(chat); err == nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err == nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
