package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Gergenus/commerce/user-service/internal/models"
	dbpkg "github.com/Gergenus/commerce/user-service/pkg/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
)

type PostgresRepository struct {
	db dbpkg.PostgresDB
}

func NewPostgresRepository(db dbpkg.PostgresDB) PostgresRepository {
	return PostgresRepository{
		db: db,
	}
}

func (p *PostgresRepository) AddUser(ctx context.Context, user models.User) (*uuid.UUID, error) {
	const op = ""
	tx, err := p.db.DB.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback(ctx)

	uid := uuid.New()
	_, err = p.db.DB.Exec(ctx, "INSERT INTO users (id, username, email, role, hashpassword) VALUES($1, $2, $3, $4, $5)", uid.String(),
		user.Username, user.Email, user.Role, user.Password)
	if err != nil {
		var pgxErr *pgconn.PgError
		if errors.As(err, &pgxErr) {
			if pgxErr.Code == "23505" {
				return nil, fmt.Errorf("%s: %w", op, ErrUserAlreadyExists)
			}
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &uid, nil
}

func (p *PostgresRepository) GetUser(ctx context.Context, email string) (*models.User, error) {
	const op = "repository.GetUser"
	var user models.User
	err := p.db.DB.QueryRow(ctx, "SELECT id, username, email, verified, role, hashpassword FROM users WHERE email = $1", email).Scan(&user.ID,
		&user.Username, &user.Email, &user.Verified, &user.Role, &user.Password)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &user, nil

}

func (p *PostgresRepository) CreateJWTSession(ctx context.Context, user models.User, refreshTokenHash, fingerprint, ip string, expiresIn int64) error {
	const op = "repository.CreateJWTSession"
	_, err := p.db.DB.Exec(ctx, "INSERT INTO refreshsessions (userId, refreshToken, fingerprint, ip, expiresIn) VALUES($1, $2, $3, $4, $5)",
		user.ID.String(), refreshTokenHash, fingerprint, ip, expiresIn)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
