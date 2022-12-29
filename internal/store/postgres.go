package store

import (
	"context"
	"fmt"

	"github.com/georgysavva/scany/sqlscan"
	_ "github.com/jackc/pgx/v4/stdlib" // use pgx
	"github.com/jmoiron/sqlx"

	pb "github.com/VladBag2022/gokeeper/internal/proto"
)

// PostgresStore is the Postgres implementation of the Store.
type PostgresStore struct {
	database *sqlx.DB
}

// NewPostgresStore connects to Postgres and returns PostgresStore.
func NewPostgresStore(
	ctx context.Context,
	databaseDSN string,
) (*PostgresStore, error) {
	db, err := sqlx.Open("pgx", databaseDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %s", err)
	}

	p := &PostgresStore{
		database: db,
	}
	return p, p.createSchema(ctx)
}

// Ping checks postgres connectivity.
func (p *PostgresStore) Ping(ctx context.Context) error {
	if err := p.database.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %s", err)
	}

	return nil
}

// Close closes postgres connection.
func (p *PostgresStore) Close() error {
	if err := p.database.Close(); err != nil {
		return fmt.Errorf("failed to close database: %s", err)
	}

	return nil
}

func (p *PostgresStore) createSchema(ctx context.Context) error {
	tables := []string{
		"CREATE EXTENSION IF NOT EXISTS pgcrypto",
		"CREATE TABLE IF NOT EXISTS users (" +
			"id SERIAL PRIMARY KEY, " +
			"username TEXT NOT NULL UNIQUE, " +
			"password TEXT NOT NULL)",
		"CREATE TABLE IF NOT EXISTS secrets (" +
			"id SERIAL PRIMARY KEY, " +
			"user_id INTEGER NOT NULL, " +
			"data BYTEA NOT NULL, " +
			"kind INTEGER NOT NULL, " +
			"FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE)",
		"CREATE TABLE IF NOT EXISTS meta (" +
			"id SERIAL PRIMARY KEY, " +
			"secret_id INTEGER NOT NULL, " +
			"text TEXT NOT NULL, " +
			"FOREIGN KEY (secret_id) REFERENCES secrets (id) ON DELETE CASCADE)",
	}
	for _, table := range tables {
		_, err := p.database.ExecContext(ctx, table)
		if err != nil {
			return fmt.Errorf("failed to execute statement: %s", err)
		}
	}

	return nil
}

// IsUsernameAvailable checks whether provided username is available.
func (p *PostgresStore) IsUsernameAvailable(
	ctx context.Context,
	username string,
) (available bool, err error) {
	var count int64

	row := p.database.QueryRowContext(ctx, "SELECT COUNT(*) FROM users WHERE username = $1", username)
	err = row.Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to scan username count: %s", err)
	}

	return count == 0, nil
}

// SignIn returns user ID from provided credentials.
func (p *PostgresStore) SignIn(ctx context.Context, credentials Credentials) (id int64, err error) {
	err = sqlscan.Get(ctx, p.database, &id,
		"SELECT id FROM users WHERE username = $1 AND password = crypt($2, password)",
		credentials.Username, credentials.Password)
	if err != nil {
		return 0, fmt.Errorf("failed to select provied user: %s", err)
	}

	return id, nil
}

// SignUp registers new user and returns his/her new ID.
func (p *PostgresStore) SignUp(ctx context.Context, credentials Credentials) (id int64, err error) {
	err = p.database.GetContext(ctx, &id,
		"INSERT INTO users (username, password) VALUES ($1, crypt($2, gen_salt('bf'))) RETURNING id",
		credentials.Username, credentials.Password)
	if err != nil {
		return 0, fmt.Errorf("failed to insert new user: %s", err)
	}

	return id, nil
}

// StoreSecret stores user secret and returns its new ID.
func (p *PostgresStore) StoreSecret(ctx context.Context, userID int64, secret Secret) (id int64, err error) {
	err = p.database.GetContext(ctx, &id,
		"INSERT INTO secrets (user_id, data, kind) VALUES ($1, $2, $3) RETURNING id",
		userID, secret.Data, secret.Kind)
	if err != nil {
		return id, fmt.Errorf("failed to insert secret: %s", err)
	}

	if secret.Kind == SecretKind(pb.SecretKind_SECRET_ENCRYPTED_KEY) {
		_, err = p.database.ExecContext(ctx,
			"DELETE FROM secrets WHERE user_id = $1 AND kind = $2 AND id != $3",
			userID, pb.SecretKind_SECRET_ENCRYPTED_KEY, id)
		if err != nil {
			return 0, fmt.Errorf("failed to delete user' other encrypted keys: %s", err)
		}
	}

	return id, nil
}

// UpdateSecret updates secret by its ID.
func (p *PostgresStore) UpdateSecret(ctx context.Context, id int64, secret Secret) error {
	if _, err := p.database.ExecContext(ctx, "UPDATE secrets SET data = $1 WHERE id = $2", secret.Data, id); err != nil {
		return fmt.Errorf("failed to update secret: %s", err)
	}

	return nil
}

// DeleteSecret deletes secret by its ID.
func (p *PostgresStore) DeleteSecret(ctx context.Context, id int64) error {
	if _, err := p.database.ExecContext(ctx, "DELETE FROM secrets WHERE id = $1", id); err != nil {
		return fmt.Errorf("failed to delete secret: %s", err)
	}

	return nil
}

// StoreMeta stores secret meta and returns ins new ID.
func (p *PostgresStore) StoreMeta(ctx context.Context, secretID int64, meta Meta) (id int64, err error) {
	err = p.database.GetContext(ctx, &id,
		"INSERT INTO meta (secret_id, text) VALUES ($1, $2) RETURNING id",
		secretID, meta)
	if err != nil {
		return 0, fmt.Errorf("failed to insert meta: %s", err)
	}

	return id, nil
}

// UpdateMeta updates meta by its ID.
func (p *PostgresStore) UpdateMeta(ctx context.Context, id int64, meta Meta) error {
	if _, err := p.database.ExecContext(ctx, "UPDATE meta SET text = $1 WHERE id = $2", meta, id); err != nil {
		return fmt.Errorf("failed to update meta: %s", err)
	}

	return nil
}

// DeleteMeta deletes meta by its ID.
func (p *PostgresStore) DeleteMeta(ctx context.Context, id int64) error {
	if _, err := p.database.ExecContext(ctx, "DELETE FROM meta WHERE id = $1", id); err != nil {
		return fmt.Errorf("failed to delete meta: %s", err)
	}

	return nil
}

// GetSecrets returns user' secrets.
func (p *PostgresStore) GetSecrets(ctx context.Context, userID int64) (secrets []StoredSecret, err error) {
	rows, err := p.database.QueryContext(ctx,
		"SELECT id, data, kind FROM secrets WHERE user_id = $1 AND kind != $2",
		userID, pb.SecretKind_SECRET_ENCRYPTED_KEY)
	if err != nil {
		return nil, fmt.Errorf("failed to query user' secrets: %s", err)
	}
	defer rows.Close()

	for rows.Next() {
		var secret StoredSecret

		err = rows.Scan(&secret.ID, &secret.Secret.Data, &secret.Secret.Kind)
		if err != nil {
			return nil, fmt.Errorf("failed to scan secret: %s", err)
		}

		secret.Meta, err = p.getSecretMeta(ctx, secret.ID)
		if err != nil {
			return nil, err
		}

		secrets = append(secrets, secret)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user' secrets: %s", err)
	}

	return secrets, nil
}

// GetEncryptedKey returns user's encrypted key.
func (p *PostgresStore) GetEncryptedKey(ctx context.Context, userID int64) (secret StoredSecret, err error) {
	row := p.database.QueryRowContext(ctx,
		"SELECT id, data, kind FROM secrets WHERE user_id = $1 AND kind = $2 LIMIT 1",
		userID, pb.SecretKind_SECRET_ENCRYPTED_KEY)

	if err = row.Scan(&secret.ID, &secret.Secret.Data, &secret.Secret.Kind); err != nil {
		return secret, fmt.Errorf("failed to scan secret: %s", err)
	}

	return secret, nil
}

// IsUserSecret checks whether secret belongs to user.
func (p *PostgresStore) IsUserSecret(ctx context.Context, userID, secretID int64) (bool, error) {
	var count int64

	row := p.database.QueryRowContext(ctx, "SELECT COUNT(*) FROM secrets WHERE id = $1 AND user_id = $2",
		secretID, userID)

	if err := row.Scan(&count); err != nil {
		return false, fmt.Errorf("failed to scan count: %s", err)
	}

	return count > 0, nil
}

// IsUserMeta checks whether meta belongs to meta.
func (p *PostgresStore) IsUserMeta(ctx context.Context, userID, metaID int64) (bool, error) {
	var count int64

	row := p.database.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM meta "+
			"JOIN secrets ON secrets.id = meta.secret_id "+
			"AND meta.id = $1 AND user_id = $2",
		metaID, userID)

	if err := row.Scan(&count); err != nil {
		return false, fmt.Errorf("failed to scan user' meta count: %s", err)
	}

	return count > 0, nil
}

func (p *PostgresStore) getSecretMeta(ctx context.Context, secretID int64) (meta []StoredMeta, err error) {
	rows, err := p.database.QueryContext(ctx, "SELECT id, text FROM meta WHERE secret_id = $1", secretID)
	if err != nil {
		return meta, fmt.Errorf("failed to query secret' meta: %s", err)
	}
	defer rows.Close()

	for rows.Next() {
		var m StoredMeta

		err = rows.Scan(&m.ID, &m.Meta)
		if err != nil {
			return meta, fmt.Errorf("failed to scan secret's meta: %s", err)
		}

		meta = append(meta, m)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve secret' meta: %s", err)
	}

	return meta, nil
}
