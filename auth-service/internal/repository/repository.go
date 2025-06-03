package repository

import (
	"context"
	"time"

	"github.com/Thanhbinh1905/realtime-chat-v2-go/auth-service/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Register(ctx context.Context, req *model.RegisterRequest) error
	GetAccountByEmail(ctx context.Context, email string) (*model.Account, error)
	CheckEmailExists(ctx context.Context, email string) (bool, error)
	SaveRefreshToken(ctx context.Context, token, userID string, expiresAt time.Time) error
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Register(ctx context.Context, req *model.RegisterRequest) error {
	_, err := r.db.Exec(ctx, "INSERT INTO accounts (id, email, password_hash) VALUES ($1, $2, $3)",
		req.ID, req.Email, req.PasswordHash)
	return err
}

func (r *repository) Login(ctx context.Context, input *model.LoginInput) (*model.Account, error) {
	var account model.Account

	err := r.db.QueryRow(ctx, "SELECT id, email, password_hash, created_at, updated_at FROM accounts WHERE email = $1", input.Email).
		Scan(&account.ID, &account.Email, &account.PasswordHash, &account.CreatedAt, &account.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (r *repository) GetAccountByEmail(ctx context.Context, email string) (*model.Account, error) {
	var account model.Account
	err := r.db.QueryRow(ctx, "SELECT id, email, password_hash, created_at, updated_at FROM accounts WHERE email = $1", email).Scan(&account.ID, &account.Email, &account.PasswordHash, &account.CreatedAt, &account.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (r *repository) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM accounts WHERE email = $1)", email).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *repository) SaveRefreshToken(ctx context.Context, token, userID string, expiresAt time.Time) error {
	_, err := r.db.Exec(ctx,
		"INSERT INTO refresh_tokens (token, user_id, expires_at) VALUES ($1, $2, $3)",
		token, userID, expiresAt)
	return err
}

// func (r *repository) GetAccountById(ctx context.Context, id uuid.UUID) (*model.Account, error) {
// 	var account model.Account
// 	err := r.db.QueryRow(ctx, "SELECT id, email, created_at, updated_at FROM accounts WHERE id = $1", email).Scan(&account.ID, &account.Email, &account.CreatedAt, &account.UpdatedAt)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &account, nil
// }
