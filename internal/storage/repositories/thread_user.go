package repositories

import (
	"database/sql"
	"muttr_chat/internal/storage/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ThreadUserRepository struct {
	DB *sqlx.DB
}

func (r *ThreadUserRepository) GetManyByThreadId(threadId uuid.UUID) ([]*models.ThreadUser, error) {
	var threadUsers []*models.ThreadUser
	err := r.DB.Select(&threadUsers,
		`
		SELECT user_id, thread_id, user_role
		FROM
			thread_users
		WHERE
			threadId = $1
	`, threadId)
	if err != nil {
		return nil, err
	}
	return threadUsers, nil
}

func (r *ThreadUserRepository) Upsert(tx *sql.Tx, threadUser *models.ThreadUser) error {
	_, err := tx.Exec(
		`
		INSERT INTO thread_users (user_id, thread_id, user_role)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, thread_id)
		DO
			UPDATE SET
				user_role = EXCLUDED.user_role
			WHERE
				(thread_users.user_role) IS DISTINCT FROM (EXCLUDED.user_role)	
		`, threadUser.UserId, threadUser.ThreadId, threadUser.UserRole,
	)
	return err
}

func (r *ThreadUserRepository) Delete(tx *sql.Tx, userId uuid.UUID, threadId uuid.UUID) error {
	_, err := r.DB.Exec(
		`
		DELETE FROM thread_users
		WHERE
			user_id = $1 AND thread_id = $2	
		`, userId, threadId,
	)
	if err != nil {
		return err
	}
	return nil
}
