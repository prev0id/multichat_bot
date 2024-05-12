package adapter

import (
	"database/sql"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/sqlite3" // поддержка sqlite3 для goqu
	"github.com/google/uuid"
	_ "modernc.org/sqlite" // драйвер sqlite3

	"multichat_bot/internal/domain"
)

const (
	dialect = "sqlite"

	tableUser  = "user"
	columnID   = "id"
	columnUUID = "uuid"

	tablePlatform = "platform"
	columnUserID  = "user_id"
	columnName    = "name"
	columnEmail   = "email"
	columnChannel = "channel"
)

type DB struct {
	db *sql.DB
}

type user struct {
	id   int
	uuid uuid.UUID
}

type platform struct {
	name    string
	email   string
	channel string
	userID  int
}

func New(path string) (*DB, error) {
	db, err := sql.Open(dialect, path)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)

	return &DB{db: db}, nil
}

func (db *DB) ListUsers() ([]*domain.User, error) {
	users, err := db.list()
	if err != nil {
		return nil, err
	}

	platforms, err := db.listPlatform()
	if err != nil {
		return nil, err
	}

	result := make([]*domain.User, 0, len(users))
	for _, user := range users {
		domainUser := &domain.User{
			UUID:      user.uuid,
			Platforms: make(map[domain.Platform]string, len(platforms[user.id])),
		}

		for _, userPlatform := range platforms[user.id] {
			platform := domain.StringToPlatform[userPlatform.name]
			domainUser.Platforms[platform] = userPlatform.channel
		}

		result = append(result, domainUser)
	}

	return result, nil
}

func (db *DB) NewUser(userUUID string) error {
	query, _, err := goqu.Dialect(dialect).
		Insert(tableUser).
		Cols(columnUUID).
		Vals(goqu.Vals{userUUID}).
		ToSQL()

	if err != nil {
		return fmt.Errorf("db::new_user error creating query: %w", err)
	}

	_, err = db.db.Exec(query)
	if err != nil {
		return fmt.Errorf("db::new_user error executing query: %w", err)
	}

	return nil

}

func (db *DB) UpdateUserPlatform(userUUID string, platform domain.Platform, value string) error {
	query, _, err := goqu.Dialect(dialect).
		Update(tableUser).
		Set(goqu.Record{
			columnName:    platform.String(),
			columnChannel: value,
		}).
		Where(goqu.Ex{columnID: userUUID}).
		ToSQL()

	if err != nil {
		return fmt.Errorf("db::update_user_platform error creating query: %w", err)
	}

	_, err = db.db.Exec(query)
	if err != nil {
		return fmt.Errorf("db::update_user_platform error executing query: %w", err)
	}

	return nil
}

func (db *DB) list() ([]user, error) {
	query, _, err := goqu.Dialect(dialect).
		Select(columnID, columnUUID).
		From(tableUser).
		ToSQL()

	if err != nil {
		return nil, fmt.Errorf("db::list_users error creating query: %w", err)
	}

	rows, err := db.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("db::list_users error executing query: %w", err)
	}
	defer rows.Close() //nolint:errcheck

	result := make([]user, 0)
	for rows.Next() {
		var (
			userUUID string
			userID   int
		)
		if err := rows.Scan(&userID, &userUUID); err != nil {
			return nil, fmt.Errorf("db::list_users error scanning row: %w", err)
		}

		parsedUUID, err := uuid.Parse(userUUID)
		if err != nil {
			return nil, fmt.Errorf("db::list_users error parsing user uuid (%s): %w", userUUID, err)
		}

		result = append(result, user{
			id:   userID,
			uuid: parsedUUID,
		})
	}

	return result, nil
}

func (db *DB) listPlatform() (map[int][]platform, error) {
	query, _, err := goqu.Dialect(dialect).
		Select(columnUserID, columnName, columnEmail, columnChannel).
		From(tableUser).
		ToSQL()

	if err != nil {
		return nil, fmt.Errorf("db::list_platforms error creating query: %w", err)
	}

	rows, err := db.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("db::list_platforms error executing query: %w", err)
	}
	defer rows.Close() //nolint:errcheck

	result := make(map[int][]platform)
	for rows.Next() {
		var (
			userID  int
			name    string
			email   string
			channel string
		)
		if err := rows.Scan(&userID, &name, &email, &channel); err != nil {
			return nil, fmt.Errorf("db::list_platforms error scanning row: %w", err)
		}

		result[userID] = append(result[userID], platform{
			userID:  userID,
			name:    name,
			email:   email,
			channel: channel,
		})
	}

	return result, nil
}
