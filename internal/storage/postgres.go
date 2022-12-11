package storage

import (
	"context"
	"database/sql"

	"github.com/georgysavva/scany/sqlscan"
	_ "github.com/jackc/pgx/v4/stdlib"

	pb "github.com/VladBag2022/gokeeper/internal/proto"
)

type PostgresRepository struct {
	database *sql.DB
}

func NewPostgresRepository(
	ctx context.Context,
	databaseDSN string,
) (*PostgresRepository, error) {
	db, err := sql.Open("pgx", databaseDSN)
	if err != nil {
		return nil, err
	}
	p := &PostgresRepository{
		database: db,
	}
	err = p.createSchema(ctx)
	return p, err
}

func (p *PostgresRepository) Ping(ctx context.Context) error {
	return p.database.PingContext(ctx)
}

func (p *PostgresRepository) Close() error {
	return p.database.Close()
}

func (p *PostgresRepository) createSchema(ctx context.Context) error {
	tables := []string{
		"CREATE EXTENSION IF NOT EXISTS pgcrypto",
		"CREATE TABLE IF NOT EXISTS users (" +
			"id SERIAL PRIMARY KEY, " +
			"username TEXT NOT NULL UNIQUE, " +
			"password TEXT NOT NULL)",
		"CREATE TABLE IF NOT EXISTS secrets (" +
			"id BIGINT PRIMARY KEY, " +
			"user_id INTEGER NOT NULL, " +
			"data BYTEA NOT NULL, " +
			"kind INTEGER NOT NULL, " +
			"FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE)",
		"CREATE TABLE IF NOT EXISTS meta (" +
			"id BIGINT PRIMARY KEY, " +
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

func (p *PostgresRepository) IsUsernameAvailable(
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

func (p *PostgresRepository) SignIn(ctx context.Context, credentials *pb.Credentials) (id int64, err error) {
	r, err := p.database.ExecContext(ctx,
		"INSERT INTO users (username, password) VALUES ($1, crypt($2, gen_salt('bf')))",
		credentials.GetUsername(), credentials.GetPassword())
	if err != nil {
		return 0, err
	}
	return r.LastInsertId()
}

func (p *PostgresRepository) SignUp(ctx context.Context, credentials *pb.Credentials) (id int64, err error) {
	err = sqlscan.Get(ctx, p.database, &id,
		"SELECT id FROM users WHERE username = $1 AND password = crypt($2, password)",
		credentials.GetUsername(), credentials.GetPassword())
	return
}

func (p *PostgresRepository) StoreSecret(ctx context.Context, userID int64, secret *pb.Secret) (id int64, err error) {
	r, err := p.database.ExecContext(ctx,
		"INSERT INTO secrets (user_id, data, type_id) VALUES ($1, $2, $3)",
		userID, secret.GetData(), secret.GetKind())
	if err != nil {
		return 0, err
	}
	return r.LastInsertId()
}

func (p *PostgresRepository) UpdateSecret(ctx context.Context, id int64, secret *pb.Secret) error {
	_, err := p.database.ExecContext(ctx,"UPDATE secrets SET data = $1 WHERE id = $2", secret.GetData(), id)
	return err
}

func (p *PostgresRepository) DeleteSecret(ctx context.Context, id int64) error {
	_, err := p.database.ExecContext(ctx,"DELETE FROM secrets WHERE id = $1", id)
	return err
}

func (p *PostgresRepository) StoreMeta(ctx context.Context, secretID int64, meta *pb.Meta) (id int64, err error) {
	r, err := p.database.ExecContext(ctx,
		"INSERT INTO meta (secret_id, text) VALUES ($1, $2)",
		secretID, meta.GetText())
	if err != nil {
		return 0, err
	}
	return r.LastInsertId()
}

func (p *PostgresRepository) UpdateMeta(ctx context.Context, id int64, meta *pb.Meta) error {
	_, err := p.database.ExecContext(ctx,"UPDATE meta SET text = $1 WHERE id = $2", meta.GetText(), id)
	return err
}

func (p *PostgresRepository) DeleteMeta(ctx context.Context, id int64) error {
	_, err := p.database.ExecContext(ctx,"DELETE FROM meta WHERE id = $1", id)
	return err
}