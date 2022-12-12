package storage

import (
	"context"
	"database/sql"

	"github.com/georgysavva/scany/sqlscan"
	_ "github.com/jackc/pgx/v4/stdlib"

	pb "github.com/VladBag2022/gokeeper/internal/proto"
)

type PostgresStore struct {
	database *sql.DB
}

func NewPostgresStore(
	ctx context.Context,
	databaseDSN string,
) (*PostgresStore, error) {
	db, err := sql.Open("pgx", databaseDSN)
	if err != nil {
		return nil, err
	}
	p := &PostgresStore{
		database: db,
	}
	err = p.createSchema(ctx)
	return p, err
}

func (p *PostgresStore) Ping(ctx context.Context) error {
	return p.database.PingContext(ctx)
}

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

func (p *PostgresStore) SignIn(ctx context.Context, credentials *pb.Credentials) (id int64, err error) {
	r, err := p.database.ExecContext(ctx,
		"INSERT INTO users (username, password) VALUES ($1, crypt($2, gen_salt('bf')))",
		credentials.GetUsername(), credentials.GetPassword())
	if err != nil {
		return 0, err
	}
	return r.LastInsertId()
}

func (p *PostgresStore) SignUp(ctx context.Context, credentials *pb.Credentials) (id int64, err error) {
	err = sqlscan.Get(ctx, p.database, &id,
		"SELECT id FROM users WHERE username = $1 AND password = crypt($2, password)",
		credentials.GetUsername(), credentials.GetPassword())
	return
}

func (p *PostgresStore) StoreSecret(ctx context.Context, userID int64, secret *pb.Secret) (id int64, err error) {
	r, err := p.database.ExecContext(ctx,
		"INSERT INTO secrets (user_id, data, type_id) VALUES ($1, $2, $3)",
		userID, secret.GetData(), secret.GetKind())
	if err != nil {
		return 0, err
	}
	return r.LastInsertId()
}

func (p *PostgresStore) UpdateSecret(ctx context.Context, id int64, secret *pb.Secret) error {
	_, err := p.database.ExecContext(ctx, "UPDATE secrets SET data = $1 WHERE id = $2", secret.GetData(), id)
	return err
}

func (p *PostgresStore) DeleteSecret(ctx context.Context, id int64) error {
	_, err := p.database.ExecContext(ctx, "DELETE FROM secrets WHERE id = $1", id)
	return err
}

func (p *PostgresStore) StoreMeta(ctx context.Context, secretID int64, meta *pb.Meta) (id int64, err error) {
	r, err := p.database.ExecContext(ctx,
		"INSERT INTO meta (secret_id, text) VALUES ($1, $2)",
		secretID, meta.GetText())
	if err != nil {
		return 0, err
	}
	return r.LastInsertId()
}

func (p *PostgresStore) UpdateMeta(ctx context.Context, id int64, meta *pb.Meta) error {
	_, err := p.database.ExecContext(ctx, "UPDATE meta SET text = $1 WHERE id = $2", meta.GetText(), id)
	return err
}

func (p *PostgresStore) DeleteMeta(ctx context.Context, id int64) error {
	_, err := p.database.ExecContext(ctx, "DELETE FROM meta WHERE id = $1", id)
	return err
}

func (p *PostgresStore) GetSecrets(ctx context.Context, userID int64) (secrets *pb.ClientSecrets, err error) {
	secrets = &pb.ClientSecrets{}

	rows, err := p.database.QueryContext(ctx, "SELECT id, data, kind FROM secrets WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		s := &pb.ClientSecret{Secret: &pb.Secret{}}
		err = rows.Scan(&s.Id, &s.Secret.Data, &s.Secret.Kind)
		if err != nil {
			return nil, err
		}

		s.Meta, err = p.getClientMeta(ctx, s.GetId())
		if err != nil {
			return nil, err
		}

		secrets.Secrets = append(secrets.Secrets, s)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return secrets, nil
}

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

func (p *PostgresStore) getClientMeta(ctx context.Context, secretID int64) (meta []*pb.ClientMeta, err error) {
	rows, err := p.database.QueryContext(ctx, "SELECT id, text FROM meta WHERE secret_id = $1", secretID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		m := &pb.ClientMeta{Meta: &pb.Meta{}}
		err = rows.Scan(&m.Id, &m.Meta.Text)
		if err != nil {
			return nil, err
		}
		meta = append(meta, m)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return meta, nil
}
