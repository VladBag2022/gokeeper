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
	database, err := sqlx.Open("pgx", databaseDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	p := &PostgresStore{
		database: database,
	}

	return p, p.createSchema(ctx)
}

// Ping checks postgres connectivity.
func (p *PostgresStore) Ping(ctx context.Context) error {
	if err := p.database.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return nil
}

// Close closes postgres connection.
func (p *PostgresStore) Close() error {
	if err := p.database.Close(); err != nil {
		return fmt.Errorf("failed to close database: %w", err)
	}

	return nil
}

// IsUsernameAvailable checks whether provided username is available.
func (p *PostgresStore) IsUsernameAvailable(
	ctx context.Context,
	username string,
) (bool, error) {
	var count int64

	row := p.database.QueryRowContext(ctx, "SELECT COUNT(*) FROM users WHERE username = $1", username)
	if err := row.Scan(&count); err != nil {
		return false, fmt.Errorf("failed to scan username count: %w", err)
	}

	return count == 0, nil
}

// SignIn returns user ID from provided credentials.
func (p *PostgresStore) SignIn(ctx context.Context, credentials Credentials) (int64, error) {
	var id int64

	err := sqlscan.Get(ctx, p.database, &id,
		"SELECT id FROM users WHERE username = $1 AND password = crypt($2, password)",
		credentials.Username, credentials.Password)
	if err != nil {
		return 0, fmt.Errorf("failed to select provied user: %w", err)
	}

	return id, nil
}

// SignUp registers new user and returns his/her new ID.
func (p *PostgresStore) SignUp(ctx context.Context, credentials Credentials) (int64, error) {
	var id int64

	err := p.database.GetContext(ctx, &id,
		"INSERT INTO users (username, password) VALUES ($1, crypt($2, gen_salt('bf'))) RETURNING id",
		credentials.Username, credentials.Password)
	if err != nil {
		return 0, fmt.Errorf("failed to insert new user: %w", err)
	}

	return id, nil
}

// StoreSecret stores user secret and returns its new ID.
func (p *PostgresStore) StoreSecret(ctx context.Context, userID int64, secret Secret) (int64, error) {
	var id int64

	err := p.database.GetContext(ctx, &id,
		"INSERT INTO secrets (user_id, data, kind) VALUES ($1, $2, $3) RETURNING id",
		userID, secret.Data, secret.Kind)
	if err != nil {
		return id, fmt.Errorf("failed to insert secret: %w", err)
	}

	if secret.Kind == SecretKind(pb.SecretKind_SECRET_ENCRYPTED_KEY) {
		_, err = p.database.ExecContext(ctx,
			"DELETE FROM secrets WHERE user_id = $1 AND kind = $2 AND id != $3",
			userID, pb.SecretKind_SECRET_ENCRYPTED_KEY, id)
		if err != nil {
			return 0, fmt.Errorf("failed to delete user' other encrypted keys: %w", err)
		}
	}

	return id, nil
}

// UpdateSecret updates secret by its ID.
func (p *PostgresStore) UpdateSecret(ctx context.Context, id int64, secret Secret) error {
	if _, err := p.database.ExecContext(ctx, "UPDATE secrets SET data = $1 WHERE id = $2", secret.Data, id); err != nil {
		return fmt.Errorf("failed to update secret: %w", err)
	}

	return nil
}

// DeleteSecret deletes secret by its ID.
func (p *PostgresStore) DeleteSecret(ctx context.Context, id int64) error {
	if _, err := p.database.ExecContext(ctx, "DELETE FROM secrets WHERE id = $1", id); err != nil {
		return fmt.Errorf("failed to delete secret: %w", err)
	}

	return nil
}

// StoreMeta stores secret meta and returns ins new ID.
func (p *PostgresStore) StoreMeta(ctx context.Context, secretID int64, meta Meta) (int64, error) {
	var id int64

	err := p.database.GetContext(ctx, &id,
		"INSERT INTO meta (secret_id, text) VALUES ($1, $2) RETURNING id",
		secretID, meta)
	if err != nil {
		return 0, fmt.Errorf("failed to insert meta: %w", err)
	}

	return id, nil
}

// UpdateMeta updates meta by its ID.
func (p *PostgresStore) UpdateMeta(ctx context.Context, id int64, meta Meta) error {
	if _, err := p.database.ExecContext(ctx, "UPDATE meta SET text = $1 WHERE id = $2", meta, id); err != nil {
		return fmt.Errorf("failed to update meta: %w", err)
	}

	return nil
}

// DeleteMeta deletes meta by its ID.
func (p *PostgresStore) DeleteMeta(ctx context.Context, id int64) error {
	if _, err := p.database.ExecContext(ctx, "DELETE FROM meta WHERE id = $1", id); err != nil {
		return fmt.Errorf("failed to delete meta: %w", err)
	}

	return nil
}

// GetSecrets returns user' secrets.
func (p *PostgresStore) GetSecrets(ctx context.Context, userID int64) ([]StoredSecret, error) {
	var secrets []StoredSecret

	rows, err := p.database.QueryContext(ctx,
		"SELECT id, data, kind FROM secrets WHERE user_id = $1 AND kind != $2",
		userID, pb.SecretKind_SECRET_ENCRYPTED_KEY)
	if err != nil {
		return nil, fmt.Errorf("failed to query user' secrets: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var secret StoredSecret

		err = rows.Scan(&secret.ID, &secret.Secret.Data, &secret.Secret.Kind)
		if err != nil {
			return nil, fmt.Errorf("failed to scan secret: %w", err)
		}

		secret.Meta, err = p.getSecretMeta(ctx, secret.ID)
		if err != nil {
			return nil, err
		}

		secrets = append(secrets, secret)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user' secrets: %w", err)
	}

	return secrets, nil
}

// GetEncryptedKey returns user's encrypted key.
func (p *PostgresStore) GetEncryptedKey(ctx context.Context, userID int64) (StoredSecret, error) {
	row := p.database.QueryRowContext(ctx,
		"SELECT id, data, kind FROM secrets WHERE user_id = $1 AND kind = $2 LIMIT 1",
		userID, pb.SecretKind_SECRET_ENCRYPTED_KEY)

	var secret StoredSecret

	if err := row.Scan(&secret.ID, &secret.Secret.Data, &secret.Secret.Kind); err != nil {
		return secret, fmt.Errorf("failed to scan secret: %w", err)
	}

	return secret, nil
}

// IsUserSecret checks whether secret belongs to user.
func (p *PostgresStore) IsUserSecret(ctx context.Context, userID, secretID int64) (bool, error) {
	var count int64

	row := p.database.QueryRowContext(ctx, "SELECT COUNT(*) FROM secrets WHERE id = $1 AND user_id = $2",
		secretID, userID)

	if err := row.Scan(&count); err != nil {
		return false, fmt.Errorf("failed to scan count: %w", err)
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
		return false, fmt.Errorf("failed to scan user' meta count: %w", err)
	}

	return count > 0, nil
}

func (p *PostgresStore) getSecretMeta(ctx context.Context, secretID int64) ([]StoredMeta, error) {
	var metas []StoredMeta

	rows, err := p.database.QueryContext(ctx, "SELECT id, text FROM meta WHERE secret_id = $1", secretID)
	if err != nil {
		return metas, fmt.Errorf("failed to query secret' meta: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var storedMeta StoredMeta

		err = rows.Scan(&storedMeta.ID, &storedMeta.Meta)
		if err != nil {
			return metas, fmt.Errorf("failed to scan secret's meta: %w", err)
		}

		metas = append(metas, storedMeta)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve secret' meta: %w", err)
	}

	return metas, nil
}
