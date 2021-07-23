\set ON_ERROR_STOP 1

\connect postgres

DROP DATABASE IF EXISTS mdb;
-- DROP user IF EXISTS mdb;
CREATE DATABASE mdb;
-- CREATE user mdb WITH PASSWORD 'mdb';

\connect mdb

CREATE SCHEMA mdb;
GRANT usage ON SCHEMA mdb TO mdb;


CREATE TABLE mdb.users
(
    id                  serial PRIMARY KEY,
    username            VARCHAR(100) UNIQUE,
    created_at          timestamp NOT NULL DEFAULT NOW()
);
GRANT SELECT, INSERT, UPDATE, DELETE ON mdb.users TO mdb;
COMMENT ON TABLE mdb.users IS 'Пользователи';


CREATE TABLE mdb.chats
(
    id                  serial PRIMARY KEY,
    name                VARCHAR(100) UNIQUE,
    created_at          timestamp NOT NULL DEFAULT NOW()
);
GRANT SELECT, INSERT, UPDATE, DELETE ON mdb.chats TO mdb;
COMMENT ON TABLE mdb.chats IS 'Чаты';


CREATE TABLE mdb.users_in_chats
(
    user_id         INTEGER REFERENCES mdb.users (id) ON DELETE CASCADE,
    chat_id         INTEGER REFERENCES mdb.chats (id) ON DELETE CASCADE
);
GRANT SELECT, INSERT, UPDATE, DELETE ON mdb.users_in_chats TO mdb;
COMMENT ON TABLE mdb.users_in_chats IS 'Связочная таблица users-chats';


CREATE TABLE mdb.messages
(
    id              serial PRIMARY KEY,
    chat_id         INTEGER REFERENCES mdb.chats (id) ON DELETE CASCADE,
    user_id         INTEGER REFERENCES mdb.users (id),
    content         VARCHAR(1000),
    created_at      timestamp NOT NULL DEFAULT NOW()
);
GRANT SELECT, INSERT, UPDATE, DELETE ON mdb.messages TO mdb;
COMMENT ON TABLE mdb.messages IS 'Сообщения';


CREATE TABLE mdb.messages_in_chat
(
    chat_id         INTEGER REFERENCES mdb.chats (id) ON DELETE CASCADE,
    message_id      INTEGER REFERENCES mdb.messages (id) ON DELETE CASCADE
);
GRANT SELECT, INSERT, UPDATE, DELETE ON mdb.messages_in_chat TO mdb;
COMMENT ON TABLE mdb.messages_in_chat IS 'Связочная таблица chat-messages';

GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA mdb TO mdb;
