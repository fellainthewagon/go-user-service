package db

import (
	"context"
	"errors"
	"rest-api/internal/apperror"
	"rest-api/internal/user"
	"rest-api/pkg/client/postgresql"
	"rest-api/pkg/logging"

	"github.com/jackc/pgconn"
)

type postgresqlStorage struct {
	client postgresql.Client
	logger *logging.Logger
}

func NewPostgresqlStorage(client postgresql.Client, logger *logging.Logger) user.Storage {
	return &postgresqlStorage{
		client: client,
		logger: logger,
	}
}

func (ps *postgresqlStorage) Create(ctx context.Context, user user.User) (string, error) {
	ps.logger.Debug("Create user...")

	q := `INSERT INTO users (email, username, password) VALUES ($1, $2, $3) RETURNING id`
	ps.logger.Tracef("SQL Query: %s", q)

	err := ps.client.QueryRow(ctx, q, user.Email, user.Username, user.PasswordHash).Scan(&user.ID)
	if err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) {
			pgError = err.(*pgconn.PgError)
			ps.logger.Errorf("SQL error: %s, details: %s, where: %s", pgError.Message, pgError.Detail, pgError.Where)
			return "", pgError
		}
		return "", err
	}

	return user.ID, nil
}

func (ps *postgresqlStorage) FindOne(ctx context.Context, id string) (user.User, error) {
	ps.logger.Debug("FindOne user...")

	q := `SELECT id, email, username, password FROM users WHERE id = $1`
	ps.logger.Tracef("SQL Query: %s", q)

	var usr user.User
	err := ps.client.QueryRow(ctx, q, id).Scan(&usr.ID, &usr.Email, &usr.Username, &usr.PasswordHash)
	if err != nil {
		ps.logger.Errorf("SQL error: %s", err)
		return user.User{}, apperror.NotFoundError
	}

	return usr, nil
}

func (ps *postgresqlStorage) FindAll(ctx context.Context) ([]user.User, error) {
	ps.logger.Debug("FindAll users...")

	q := `SELECT id, email, username, password FROM users`
	ps.logger.Tracef("SQL Query: %s", q)

	rows, err := ps.client.Query(ctx, q)
	if err != nil {
		ps.logger.Errorf("SQL error: %s", err)
		return nil, err
	}

	users := make([]user.User, 0)
	for rows.Next() {
		var usr user.User

		err = rows.Scan(&usr.ID, &usr.Email, &usr.Username, &usr.PasswordHash)
		if err != nil {
			ps.logger.Errorf("Scan error: %s", err)
			return nil, err
		}

		users = append(users, usr)
	}

	if err = rows.Err(); err != nil {
		ps.logger.Errorf("SQL error: %s", err)
		return nil, err
	}

	return users, nil
}

func (ps *postgresqlStorage) Update(ctx context.Context, user user.User) error {
	ps.logger.Debug("Update user...")

	q := `UPDATE users SET email = $1, username = $2, password = $3 WHERE id = $4`
	ps.logger.Tracef("SQL Query: %s", q)

	tx, err := ps.client.Begin(ctx)
	if err != nil {
		ps.logger.Errorf("Transaction error: ", err)
		return err
	}
	defer tx.Rollback(ctx)

	commandTag, err := ps.client.Exec(ctx, q, user.Email, user.Username, user.PasswordHash, user.ID)
	if err != nil {
		ps.logger.Errorf("SQL error: %s", err)
		return err
	}
	if commandTag.RowsAffected() != 1 {
		ps.logger.Errorf("RowsAffected: %s", err)
		return apperror.NotFoundError
	}

	tx.Commit(ctx)
	ps.logger.Infof("Successfully updated, count: %d", commandTag.RowsAffected())
	return nil
}

func (ps *postgresqlStorage) Delete(ctx context.Context, id string) error {
	ps.logger.Debug("Delete user...")

	q := `DELETE FROM users WHERE id = $1`
	ps.logger.Tracef("SQL Query: %s", q)

	commandTag, err := ps.client.Exec(ctx, q, id)
	if err != nil {
		ps.logger.Errorf("SQL error: %s", err)
		return err
	}
	if commandTag.RowsAffected() != 1 {
		ps.logger.Errorf("RowsAffected: %s", err)
		return apperror.NotFoundError
	}

	ps.logger.Infof("Successfully deleted, count: %d", commandTag.RowsAffected())
	return nil
}
