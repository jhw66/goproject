package internal

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	VisibilityPublic  = "public"
	VisibilityPrivate = "private"
)

var ErrForbiddenRoom = errors.New("forbidden room")
var ErrInvalidRoomPassword = errors.New("invalid room password")

// Room represents a chat room entity.
type Room struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Visibility string    `json:"visibility"`
	Password   string    `json:"-"`
	PasswordHS string    `json:"-"`
	CreatedBy  string    `json:"createdBy"`
	CreatedAt  time.Time `json:"createdAt"`
}

// Store encapsulates MySQL operations.
type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) EnsureRoomAccess(ctx context.Context, roomID, userID string) error {
	var visibility string
	err := s.db.QueryRowContext(ctx, "SELECT visibility FROM chat_rooms WHERE id = ?", roomID).Scan(&visibility)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("room %s not found", roomID)
		}
		return err
	}

	if visibility == VisibilityPublic {
		return nil
	}

	var exists int
	err = s.db.QueryRowContext(ctx,
		"SELECT 1 FROM chat_room_members WHERE room_id = ? AND user_id = ? LIMIT 1",
		roomID, userID,
	).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrForbiddenRoom
		}
		return err
	}
	return nil
}

func (s *Store) SaveMessage(ctx context.Context, msg ChatMessage) (*ChatMessage, error) {
	res, err := s.db.ExecContext(ctx,
		"INSERT INTO chat_messages(room_id, user_id, nickname, content) VALUES(?, ?, ?, ?)",
		msg.RoomID, msg.UserID, msg.Nickname, msg.Content,
	)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	msg.ID = id
	if msg.CreatedAt.IsZero() {
		msg.CreatedAt = time.Now().UTC()
	}
	return &msg, nil
}

func (s *Store) ListVisibleRooms(ctx context.Context, userID string) ([]Room, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT r.id, r.name, r.visibility, r.created_by, r.created_at
		FROM chat_rooms r
		LEFT JOIN chat_room_members m ON r.id = m.room_id AND m.user_id = ?
		WHERE r.visibility = 'public' OR m.user_id IS NOT NULL
		ORDER BY r.created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rooms := make([]Room, 0)
	for rows.Next() {
		var room Room
		if err := rows.Scan(&room.ID, &room.Name, &room.Visibility, &room.CreatedBy, &room.CreatedAt); err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}
	return rooms, rows.Err()
}

func (s *Store) CreateRoom(ctx context.Context, room Room) error {
	room.Visibility = strings.ToLower(strings.TrimSpace(room.Visibility))
	if room.Visibility != VisibilityPublic && room.Visibility != VisibilityPrivate {
		return errors.New("invalid visibility")
	}

	passwordHS := strings.TrimSpace(room.PasswordHS)
	if room.Visibility == VisibilityPrivate && passwordHS == "" && strings.TrimSpace(room.Password) != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(room.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		passwordHS = string(hash)
	}

	_, err := s.db.ExecContext(ctx,
		"INSERT INTO chat_rooms(id, name, visibility, password_hash, created_by) VALUES(?, ?, ?, ?, ?)",
		room.ID, room.Name, room.Visibility, nullable(passwordHS), room.CreatedBy,
	)
	return err
}

func (s *Store) AddMember(ctx context.Context, roomID, userID, role string) error {
	if strings.TrimSpace(role) == "" {
		role = "member"
	}
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO chat_room_members(room_id, user_id, role)
		VALUES(?, ?, ?)
		ON DUPLICATE KEY UPDATE role = VALUES(role)
	`, roomID, userID, role)
	return err
}

func (s *Store) IsRoomOwner(ctx context.Context, roomID, userID string) (bool, error) {
	var createdBy string
	err := s.db.QueryRowContext(ctx, "SELECT created_by FROM chat_rooms WHERE id = ?", roomID).Scan(&createdBy)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("room %s not found", roomID)
		}
		return false, err
	}
	if createdBy == userID {
		return true, nil
	}

	var role string
	err = s.db.QueryRowContext(ctx,
		"SELECT role FROM chat_room_members WHERE room_id = ? AND user_id = ? LIMIT 1",
		roomID, userID,
	).Scan(&role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return strings.EqualFold(role, "owner"), nil
}

func (s *Store) JoinRoomByPassword(ctx context.Context, roomID, userID, password string) error {
	password = strings.TrimSpace(password)
	if password == "" {
		return ErrInvalidRoomPassword
	}

	var visibility string
	var hash sql.NullString
	err := s.db.QueryRowContext(ctx,
		"SELECT visibility, password_hash FROM chat_rooms WHERE id = ?",
		roomID,
	).Scan(&visibility, &hash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("room %s not found", roomID)
		}
		return err
	}

	if visibility != VisibilityPrivate {
		return errors.New("room is not private")
	}
	if !hash.Valid || strings.TrimSpace(hash.String) == "" {
		return ErrInvalidRoomPassword
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hash.String), []byte(password)); err != nil {
		return ErrInvalidRoomPassword
	}
	return s.AddMember(ctx, roomID, userID, "member")
}

func (s *Store) ListMessages(ctx context.Context, roomID string, limit int) ([]ChatMessage, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, room_id, user_id, nickname, content, created_at
		FROM chat_messages
		WHERE room_id = ?
		ORDER BY id DESC
		LIMIT ?
	`, roomID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]ChatMessage, 0, limit)
	for rows.Next() {
		var msg ChatMessage
		if err := rows.Scan(&msg.ID, &msg.RoomID, &msg.UserID, &msg.Nickname, &msg.Content, &msg.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, msg)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
	return list, nil
}

func nullable(s string) interface{} {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	return s
}
