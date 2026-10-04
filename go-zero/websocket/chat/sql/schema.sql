CREATE TABLE IF NOT EXISTS chat_rooms (
    id VARCHAR(64) PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    visibility ENUM('public', 'private') NOT NULL DEFAULT 'public',
    password_hash VARCHAR(255) NULL,
    created_by VARCHAR(64) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS chat_room_members (
    room_id VARCHAR(64) NOT NULL,
    user_id VARCHAR(64) NOT NULL,
    role VARCHAR(32) NOT NULL DEFAULT 'member',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (room_id, user_id),
    CONSTRAINT fk_room_members_room
        FOREIGN KEY (room_id) REFERENCES chat_rooms(id)
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS chat_messages (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    room_id VARCHAR(64) NOT NULL,
    user_id VARCHAR(64) NOT NULL,
    nickname VARCHAR(64) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_messages_room_time (room_id, created_at),
    CONSTRAINT fk_messages_room
        FOREIGN KEY (room_id) REFERENCES chat_rooms(id)
        ON DELETE CASCADE
);
