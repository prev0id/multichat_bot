package adapter

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/sqlite3" // поддержка sqlite3 для goqu
	_ "modernc.org/sqlite"                             // драйвер sqlite3

	"multichat_bot/internal/domain"
)

const (
	dialect = "sqlite"

	tableUser     = "user"
	tablePlatform = "platform"

	columnID            = "id"
	columnUserID        = "user_id"
	columnName          = "name"
	columnChannel       = "channel"
	columnAccessToken   = "access_token"
	columnRefreshToken  = "refresh_token"
	columnExpiresIn     = "expires_in"
	columnDisabledUsers = "disabled_users"
	columnBannedWords   = "banned_words"
	columnIsJoined      = "is_joined"
)

type DB struct {
	db *sql.DB
}

type userRow struct {
	token string
	id    int64
}

type platformRow struct {
	Name          string `db:"name"`
	ID            string `db:"id"`
	Channel       string `db:"channel"`
	AccessToken   string `db:"access_token"`
	RefreshToken  string `db:"refresh_token"`
	ExpiresIn     string `db:"expires_in"`
	DisabledUsers string `db:"disabled_users"`
	BannedWords   string `db:"banned_words"`
	UserID        int64  `db:"user_id"`
	IsJoined      int    `db:"is_joined"`
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
	users, err := db.listUserTable()
	if err != nil {
		return nil, err
	}

	platforms, err := db.listPlatformTable()
	if err != nil {
		return nil, err
	}

	result := make([]*domain.User, 0, len(users))
	for _, user := range users {
		configs, err := convertPlatformsToDomain(platforms[user.id])
		if err != nil {
			return nil, err
		}

		result = append(result, &domain.User{
			ID:          user.id,
			AccessToken: user.token,
			Platforms:   configs,
		})
	}

	return result, nil
}

func (db *DB) NewUser(token string) (int64, error) {
	query, _, err := goqu.Dialect(dialect).
		Insert(tableUser).
		Cols(columnID, columnAccessToken).
		FromQuery(
			goqu.From(tableUser).
				Select(
					goqu.L("ifnull(max(id), 0) + 1").As(columnID),
					goqu.C(token).As(columnAccessToken),
				),
		).ToSQL()

	if err != nil {
		return 0, fmt.Errorf("db::new_user error creating query: %w", err)
	}

	res, err := db.db.Exec(query)
	if err != nil {
		return 0, fmt.Errorf("db::new_user error executing query: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("db::new_user error getting last insert id: %w", err)
	}

	return id, nil
}

func (db *DB) UpsertPlatform(id int64, platform domain.Platform, config *domain.PlatformConfig) error {
	converted := convertPlatformToDB(id, platform, config)

	query, _, err := goqu.Dialect(dialect).
		Insert(tablePlatform).
		Cols(
			columnUserID,
			columnName,
			columnID,
			columnChannel,
			columnIsJoined,
			columnAccessToken,
			columnRefreshToken,
			columnExpiresIn,
			columnDisabledUsers,
			columnBannedWords,
		).
		Rows(
			converted,
		).OnConflict(
		goqu.DoUpdate(columnID, converted),
	).ToSQL()

	if err != nil {
		return fmt.Errorf("db::upsert_platform error creating query: %w", err)
	}

	_, err = db.db.Exec(query)
	if err != nil {
		return fmt.Errorf("db::upsert_platform error exequting query: %w", err)
	}

	return nil
}

func (db *DB) ChangeJoined(id int64, platform domain.Platform, value bool) error {
	query, _, err := goqu.Dialect(dialect).
		Update(tablePlatform).
		Where(
			goqu.And(
				goqu.C(columnName).Eq(platform.String()),
				goqu.C(columnUserID).Eq(id),
			),
		).
		Set(
			goqu.Record{columnIsJoined: value},
		).
		ToSQL()

	if err != nil {
		return fmt.Errorf("db::change_joined error creating query: %w", err)
	}

	_, err = db.db.Exec(query)
	if err != nil {
		return fmt.Errorf("db::change_joined error exequting query: %w", err)
	}

	return nil
}

func (db *DB) DeleteUser(id int64) error {
	query, _, err := goqu.Dialect(dialect).
		Delete(tablePlatform).
		Where(
			goqu.C(columnUserID).Eq(id),
		).
		ToSQL()

	if err != nil {
		return fmt.Errorf("db::delete_platform error creating query: %w", err)
	}

	_, err = db.db.Exec(query)

	if err != nil {
		return fmt.Errorf("db::delete_platform error exequting query: %w", err)
	}

	return nil
}

func (db *DB) UpdateBannedUsers(id int64, platform domain.Platform, bannedUsers domain.BannedList) error {
	query, _, err := goqu.Dialect(dialect).
		Update(tablePlatform).
		Where(
			goqu.And(
				goqu.C(columnUserID).Eq(id),
				goqu.C(columnName).Eq(platform.String()),
			),
		).
		Set(
			goqu.Record{columnDisabledUsers: convertBannedListToDB(bannedUsers)},
		).
		ToSQL()
	if err != nil {
		return fmt.Errorf("db::update_banned_users error creating query: %w", err)
	}

	slog.Info(query)

	_, err = db.db.Exec(query)
	if err != nil {
		return fmt.Errorf("db::update_banned_users error exequting query: %w", err)
	}

	return nil
}

func (db *DB) UpdateBannedWords(id int64, platform domain.Platform, bannedWords domain.BannedList) error {
	query, _, err := goqu.Dialect(dialect).
		Update(tablePlatform).
		Where(
			goqu.And(
				goqu.C(columnUserID).Eq(id),
				goqu.C(columnName).Eq(platform.String()),
			),
		).
		Set(
			goqu.Record{columnBannedWords: convertBannedListToDB(bannedWords)},
		).
		ToSQL()
	if err != nil {
		return fmt.Errorf("db::update_banned_words error creating query: %w", err)
	}

	_, err = db.db.Exec(query)
	if err != nil {
		return fmt.Errorf("db::update_banned_words error exequting query: %w", err)
	}

	return nil
}

func (db *DB) listUserTable() ([]userRow, error) {
	query, _, err := goqu.Dialect(dialect).
		Select(columnID, columnAccessToken).
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

	result := make([]userRow, 0)
	for rows.Next() {
		var (
			userID int64
			token  string
		)

		if err := rows.Scan(&userID, &token); err != nil {
			return nil, fmt.Errorf("db::list_users error scanning row: %w", err)
		}

		result = append(result, userRow{
			id:    userID,
			token: token,
		})
	}

	return result, nil
}

func (db *DB) listPlatformTable() (map[int64][]platformRow, error) {
	query, _, err := goqu.Dialect(dialect).
		Select(
			columnUserID,
			columnName,
			columnID,
			columnChannel,
			columnIsJoined,
			columnAccessToken,
			columnRefreshToken,
			columnExpiresIn,
			columnDisabledUsers,
			columnBannedWords,
		).
		From(tablePlatform).
		ToSQL()

	if err != nil {
		return nil, fmt.Errorf("db::list_platforms error creating query: %w", err)
	}

	rows, err := db.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("db::list_platforms error executing query: %w", err)
	}
	defer rows.Close() //nolint:errcheck

	result := make(map[int64][]platformRow)
	for rows.Next() {
		var (
			userID        int64
			name          string
			id            string
			channel       string
			isJoined      int
			accessToken   string
			refreshToken  string
			expiresIn     string
			disabledUsers string
			bannedWords   string
		)

		err := rows.Scan(
			&userID,
			&name,
			&id,
			&channel,
			&isJoined,
			&accessToken,
			&refreshToken,
			&expiresIn,
			&disabledUsers,
			&bannedWords,
		)

		if err != nil {
			return nil, fmt.Errorf("db::list_platforms error scanning row: %w", err)
		}

		result[userID] = append(result[userID], platformRow{
			UserID:        userID,
			Name:          name,
			ID:            id,
			Channel:       channel,
			IsJoined:      isJoined,
			AccessToken:   accessToken,
			RefreshToken:  refreshToken,
			ExpiresIn:     expiresIn,
			DisabledUsers: disabledUsers,
			BannedWords:   bannedWords,
		})
	}

	return result, nil
}

func convertBannedListToDB(list domain.BannedList) string {
	result, _ := json.Marshal(list)
	return string(result)
}
