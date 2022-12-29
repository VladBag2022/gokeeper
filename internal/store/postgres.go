package store

import (
	"context"

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
		return nil, err
	}
	p := &PostgresStore{
		database: db,
	}
	err = p.createSchema(ctx)
	return p, err
}

// Ping checks postgres connectivity.
func (p *PostgresStore) Ping(ctx context.Context) error {
	return p.database.PingContext(ctx)
}

// Close closes postgres connection.
func (p *PostgresStore) Close() error {
	return p.database.Close()
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
			return err
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
		return false, err
	}
	return count == 0, err
}

// SignIn returns user ID from provided credentials.
func (p *PostgresStore) SignIn(ctx context.Context, credentials Credentials) (id int64, err error) {
	return id, sqlscan.Get(ctx, p.database, &id,
		"SELECT id FROM users WHERE username = $1 AND password = crypt($2, password)",
		credentials.Username, credentials.Password)
}

// SignUp registers new user and returns his/her new ID.
func (p *PostgresStore) SignUp(ctx context.Context, credentials Credentials) (id int64, err error) {
	return id, p.database.GetContext(ctx, &id,
		"INSERT INTO users (username, password) VALUES ($1, crypt($2, gen_salt('bf'))) RETURNING id",
		credentials.Username, credentials.Password)
}

// StoreSecret stores user secret and returns its new ID.
func (p *PostgresStore) StoreSecret(ctx context.Context, userID int64, secret Secret) (id int64, err error) {
	err = p.database.GetContext(ctx, &id,
		"INSERT INTO secrets (user_id, data, kind) VALUES ($1, $2, $3) RETURNING id",
		userID, secret.Data, secret.Kind)
	if err != nil {
		return
	}

	if secret.Kind == SecretKind(pb.SecretKind_SECRET_ENCRYPTED_KEY) {
		_, err = p.database.ExecContext(ctx,
			"DELETE FROM secrets WHERE user_id = $1 AND kind = $2 AND id != $3",
			userID, pb.SecretKind_SECRET_ENCRYPTED_KEY, id)
	}
	return
}

// UpdateSecret updates secret by its ID.
func (p *PostgresStore) UpdateSecret(ctx context.Context, id int64, secret Secret) error {
	_, err := p.database.ExecContext(ctx, "UPDATE secrets SET data = $1 WHERE id = $2", secret.Data, id)
	return err
}

// DeleteSecret deletes secret by its ID.
func (p *PostgresStore) DeleteSecret(ctx context.Context, id int64) error {
	_, err := p.database.ExecContext(ctx, "DELETE FROM secrets WHERE id = $1", id)
	return err
}

// StoreMeta stores secret meta and returns ins new ID.
func (p *PostgresStore) StoreMeta(ctx context.Context, secretID int64, meta Meta) (id int64, err error) {
	return id, p.database.GetContext(ctx, &id,
		"INSERT INTO meta (secret_id, text) VALUES ($1, $2) RETURNING id",
		secretID, meta)
}

// UpdateMeta updates meta by its ID.
func (p *PostgresStore) UpdateMeta(ctx context.Context, id int64, meta Meta) error {
	_, err := p.database.ExecContext(ctx, "UPDATE meta SET text = $1 WHERE id = $2", meta, id)
	return err
}

// DeleteMeta deletes meta by its ID.
func (p *PostgresStore) DeleteMeta(ctx context.Context, id int64) error {
	_, err := p.database.ExecContext(ctx, "DELETE FROM meta WHERE id = $1", id)
	return err
}

// GetSecrets returns client' secrets.
func (p *PostgresStore) GetSecrets(ctx context.Context, userID int64) (secrets []ClientSecret, err error) {
	rows, err := p.database.QueryContext(ctx,
		"SELECT id, data, kind FROM secrets WHERE user_id = $1 AND kind != $2",
		userID, pb.SecretKind_SECRET_ENCRYPTED_KEY)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var s ClientSecret
		err = rows.Scan(&s.ID, &s.Secret.Data, &s.Secret.Kind)
		if err != nil {
			return nil, err
		}

		s.Meta, err = p.getClientMeta(ctx, s.ID)
		if err != nil {
			return nil, err
		}

		secrets = append(secrets, s)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return secrets, nil
}

// GetEncryptedKey returns client's encrypted key.
func (p *PostgresStore) GetEncryptedKey(ctx context.Context, userID int64) (secret ClientSecret, err error) {
	row := p.database.QueryRowContext(ctx,
		"SELECT id, data, kind FROM secrets WHERE user_id = $1 AND kind = $2 LIMIT 1",
		userID, pb.SecretKind_SECRET_ENCRYPTED_KEY)
	return secret, row.Scan(&secret.ID, &secret.Secret.Data, &secret.Secret.Kind)
}

// IsUserSecret checks whether secret belongs to user.
func (p *PostgresStore) IsUserSecret(ctx context.Context, userID, secretID int64) (bool, error) {
	var count int64
	row := p.database.QueryRowContext(ctx, "SELECT COUNT(*) FROM secrets WHERE id = $1 AND user_id = $2",
		secretID, userID)
	err := row.Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, err
}

// IsUserMeta checks whether meta belongs to meta.
func (p *PostgresStore) IsUserMeta(ctx context.Context, userID, metaID int64) (bool, error) {
	var count int64
	row := p.database.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM meta "+
			"JOIN secrets ON secrets.id = meta.secret_id "+
			"AND meta.id = $1 AND user_id = $2",
		metaID, userID)
	err := row.Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, err
}

func (p *PostgresStore) getClientMeta(ctx context.Context, secretID int64) (meta []ClientMeta, err error) {
	rows, err := p.database.QueryContext(ctx, "SELECT id, text FROM meta WHERE secret_id = $1", secretID)
	if err != nil {
		return meta, err
	}
	defer rows.Close()

	for rows.Next() {
		var m ClientMeta
		err = rows.Scan(&m.ID, &m.Meta)
		if err != nil {
			return meta, err
		}
		meta = append(meta, m)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return meta, nil
}
