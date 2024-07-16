package repositories

import (
	"database/sql"
	"muttr_chat/internal/storage/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ThreadRepository struct {
	DB *sqlx.DB
}

func (r *ThreadRepository) GetAll() ([]*models.Thread, error) {
	var threads []*models.Thread
	err := r.DB.Select(&threads,
		`
			SELECT id, thread_type, server_id, name, voice_enabled, created_at, updated_at, deleted_at
			FROM threads
		`,
	)
	if err != nil {
		return nil, err
	}
	return threads, nil
}

func (r *ThreadRepository) GetById(id uuid.UUID) (*models.Thread, error) {
	var thread models.Thread
	err := r.DB.Get(&thread,
		`
		SELECT id, thread_type, server_id, name, voice_enabled, created_at, updated_at, deleted_at
		FROM threads
		WHERE
			id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	return &thread, nil
}

func (r *ThreadRepository) Upsert(tx *sql.Tx, thread *models.Thread) error {
	_, err := tx.Exec(
		`
		INSERT INTO threads (id, thread_type, server_id, name, voice_enabled, created_at, updated_at, deleted_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (id)
		DO
			UPDATE SET
				thread_type = EXCLUDED.thread_type,
				name = EXCLUDED.name,
				voice_enabled = EXCLUDED.voice_enabled,
				updated_at = now(),
				deleted_at = EXCLUDED.deleted_at
			WHERE
				(threads.thread_type, threads.name, threads.voice_enabled, threads.deleted_at) IS DISTINCT FROM
				(EXCLUDED.thread_type, EXCLUDED.name, EXCLUDED.voice_enabled, EXCLUDED.deleted_at)
	`, thread.Id, thread.ThreadType, thread.ServerId, thread.Name, thread.VoiceEnabled, thread.CreatedAt, thread.UpdatedAt, thread.DeletedAt)
	return err
}
