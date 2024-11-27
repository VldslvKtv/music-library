package pg

import (
	"context"
	"errors"
	"fmt"
	"music_library/config"
	"music_library/internal/http_server/lib/utils"
	"music_library/internal/http_server/models"
	"music_library/internal/http_server/storage"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const errCode = "23505"

type Storage struct {
	DB *pgxpool.Pool
}

func New(cfg *config.Config) (*Storage, error) {
	const op = "storage.pg.New"

	databaseUrl := cfg.StoragePath
	dbPool, err := pgxpool.New(context.Background(), databaseUrl)
	if err != nil {
		return nil, fmt.Errorf("%s :%w", op, err)
	}
	if err := runMigrations(cfg); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}
	return &Storage{DB: dbPool}, nil

}

func runMigrations(cfg *config.Config) error {
	const op = "storage.pg.runMigrations"
	m, err := migrate.New(cfg.MigrationsPath, cfg.StoragePath)
	if err != nil {
		return fmt.Errorf("%s :%w", op, err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrate up: %w", err)
	}

	return nil

}
func (s *Storage) Close() {
	defer s.DB.Close()
}

func getWhereClauses(filter map[string]interface{}) (string, []interface{}, int) {
	var whereClauses []string
	var args []interface{}
	argID := 1

	// Генерация условий фильтрации
	for key, value := range filter {
		whereClauses = append(whereClauses, fmt.Sprintf("%s = $%d", key, argID))
		args = append(args, value)
		argID++
	}

	whereSQL := ""
	if len(whereClauses) > 0 {
		whereSQL = "WHERE " + strings.Join(whereClauses, " AND ")
	}
	return whereSQL, args, argID
}

func (s *Storage) GetCountSongs(ctx context.Context, filter map[string]interface{}) (int, error) {
	const op = "storage.pg.GetCountSongs"

	whereSQL, _, _ := getWhereClauses(filter)

	countSongs := fmt.Sprintf(`
	 SELECT COUNT(*) 
	 FROM groups 
	 JOIN songs ON groups.id = songs.group_id
	 JOIN song_details ON songs.id = song_details.song_id
     %s
	`, whereSQL)

	var totalSongs int
	if err := s.DB.QueryRow(ctx, countSongs).Scan(&totalSongs); err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return totalSongs, nil

}

func (s *Storage) GetData(ctx context.Context, filter map[string]interface{}, page int, pageSize int) ([]models.Data, error) {
	const op = "storage.pg.GetData"
	offset := (page - 1) * pageSize

	whereSQL, args, argID := getWhereClauses(filter)

	query := fmt.Sprintf(`
        SELECT groups.name, songs.name, release_date, text, link
        FROM groups
		JOIN songs ON groups.id = songs.group_id
		JOIN song_details ON songs.id = song_details.song_id
        %s
        ORDER BY songs.id
        LIMIT $%d OFFSET $%d
    `, whereSQL, argID, argID+1)

	args = append(args, pageSize, offset)

	rows, err := s.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var songs []models.Data
	for rows.Next() {
		var song models.Data
		err := rows.Scan(&song.Group, &song.Song, &song.ReleaseDate.Time, &song.Text, &song.Link)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		songs = append(songs, song)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("%s: %w", op, rows.Err())
	}

	return songs, nil
}

func (s *Storage) GetSong(ctx context.Context, group string, song string) (string, error) {
	const op = "storage.pg.GetSong"

	query := `
        SELECT song_details.text
        FROM songs
        JOIN groups ON groups.id = songs.group_id
        JOIN song_details ON songs.id = song_details.song_id
        WHERE groups.name = $1 AND songs.name = $2
    `

	row := s.DB.QueryRow(ctx, query, group, song)

	var textSong string
	err := row.Scan(&textSong)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("%s; %w", op, storage.ErrSongNotFound)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return textSong, nil
}

func (s *Storage) CreateSong(ctx context.Context, data models.Data) error {
	const op = "storage.pg.CreateSong"

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s; failed to begin transaction: %w", op, err)
	}
	defer tx.Rollback(ctx)

	var groupID int
	err = tx.QueryRow(ctx, `
        INSERT INTO groups (name)
        VALUES ($1)
		RETURNING id
    `, data.Group).Scan(&groupID)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == errCode {
			return fmt.Errorf("%s; %w", op, storage.ErrGroupExists)
		}
		return fmt.Errorf("%s; failed to insert into groups: %w", op, err)
	}

	var songID int
	err = tx.QueryRow(ctx, `
        INSERT INTO songs (group_id, name)
        VALUES ($1, $2)
		RETURNING id
    `, groupID, data.Song).Scan(&songID)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == errCode {
			return fmt.Errorf("%s; %w", op, storage.ErrSongExists)
		}
		return fmt.Errorf("%s; failed to insert into songs: %w", op, err)
	}

	_, err = tx.Exec(ctx, `
        INSERT INTO song_details (song_id, release_date, text, link)
        VALUES ($1, $2, $3, $4)
    `, songID, data.ReleaseDate.Time.Format(models.CustomTimeFormat), data.Text, data.Link)
	if err != nil {
		return fmt.Errorf("%s; failed to insert into song_details: %w", op, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("%s; failed to commit transaction: %w", op, err)
	}

	return nil
}

func (s *Storage) DeleteSong(ctx context.Context, idSong int) error {
	const op = "storage.pg.DeleteSong"

	query := `
        DELETE FROM songs
        WHERE id = $1
    `
	result, err := s.DB.Exec(ctx, query, idSong)
	if err != nil {
		return fmt.Errorf("%s: failed to delete from songs: %w", op, err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("%s: %w", op, storage.ErrSongNotFound)
	}

	return nil
}

func (s *Storage) PatchSong(ctx context.Context, idSong int, data models.Data) error {
	const op = "storage.pg.PatchSong"

	var exists bool
	err := s.DB.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM songs WHERE id = $1)", idSong).Scan(&exists)
	if err != nil {
		return fmt.Errorf("%s: failed to check song existence: %w", op, err)
	}
	if !exists {
		return fmt.Errorf("%s: song id does not exist %w", op, storage.ErrSongNotFound)
	}

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: failed to begin transaction: %w", op, err)
	}
	defer tx.Rollback(ctx)

	// Конвертируем структуру в map для удобства обработки
	mapData := utils.ConvertStruct(data)
	if len(mapData) == 0 {
		return fmt.Errorf("%s: no changes", op)
	}

	if group, ok := mapData["groups.name"]; ok {
		query := `
            UPDATE groups
            SET name = $1
            WHERE id = (SELECT group_id FROM songs WHERE id = $2)
        `
		_, err := tx.Exec(ctx, query, group, idSong)
		if err != nil {
			if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == errCode {
				return fmt.Errorf("%s; %w", op, storage.ErrGroupExists)
			}
			return fmt.Errorf("%s: failed to update group: %w", op, err)
		}
		delete(mapData, "groups.name")
	}

	if song, ok := mapData["songs.name"]; ok {
		query := `
            UPDATE songs
            SET name = $1
            WHERE id = $2
        `
		_, err := tx.Exec(ctx, query, song, idSong)
		if err != nil {
			if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == errCode {
				return fmt.Errorf("%s; %w", op, storage.ErrSongExists)
			}
			return fmt.Errorf("%s: failed to update song name: %w", op, err)
		}
		delete(mapData, "songs.name")
	}

	var setClauses []string
	var args []interface{}
	argID := 1

	for key, value := range mapData {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", strings.Replace(key, ".", "_", -1), argID))
		args = append(args, value)
		argID++
	}

	if len(setClauses) > 0 {
		setSQL := "SET " + strings.Join(setClauses, ", ")
		query := fmt.Sprintf(`
            UPDATE song_details
            %s
            WHERE song_id = $%d
        `, setSQL, argID)

		args = append(args, idSong)

		_, err := tx.Exec(ctx, query, args...)
		if err != nil {
			return fmt.Errorf("%s: failed to update song details: %w", op, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("%s: failed to commit transaction: %w", op, err)
	}

	return nil
}
