package user


import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RepositoryInterface interface {

	Create(ctx context.Context, u *User) error

	GetByID(ctx context.Context, id int64) (*User, error)

	GetByEmail(ctx context.Context, email string) (*User, error)

	Update(ctx context.Context, u *User) error

	Delete(ctx context.Context, id int64) error

	List(ctx context.Context, limit, offset int) ([]*User, error)
}

var (
	ErrNotFound = errors.New("user not found")
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Create(ctx context.Context, u *User) error {
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("failed to acquire connection: %w", err)
	}
	defer conn.Release()

	query := `
		INSERT INTO users (email, password_hash, role, created_at)
		VALUES ($1, $2, $3, NOW())
		RETURNING id
	`
	err = conn.QueryRow(ctx, query, u.Email, u.PasswordHash, u.Role).Scan(&u.Id)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}
	return nil
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*User, error) {
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to acquire connection: %w", err)
	}
	defer conn.Release()

	query := `
		SELECT email, password_hash, role, created_at
		FROM users
		WHERE id = $1
	`
	var user User

	err = conn.QueryRow(ctx, query, id).Scan(&user.Id, &user.Email, &user.PasswordHash, &user.Role)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return &user, nil
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (*User, error) {
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to acquire connection: %w", err)
	}
	defer conn.Release()

	query := `
		SELECT email, password_hash, role, created_at
		FROM users
		WHERE email = $1
	`
	var user User

	err = conn.QueryRow(ctx, query, email).Scan(&user.Id, &user.Email, &user.PasswordHash, &user.Role)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}

func (r *Repository) Update(ctx context.Context, u *User) error {
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("failed to acquire connection: %w", err)
	}
	defer conn.Release()

	query := `
		UPDATE users
		SET email=$1, password_hash=$2, role=$3
		WHERE id = $4
	`
	ar, err := conn.Exec(ctx, query, u.Email, u.PasswordHash, u.Role, u.Id)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	if ar.RowsAffected() != 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("failed to acquire connection: %w", err)
	}
	defer conn.Release()

	query := `
		DELETE FROM users
		WHERE id = $1
	`
	ar, err := conn.Exec(ctx, query, id)

	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if ar.RowsAffected() != 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) List(ctx context.Context, limit, offset int) ([]*User, error) {
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to acquire connection: %w", err)
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, `
		SELECT id, email, password_hash, role, created_at
		FROM users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.Id, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan user row: %w", err)
		}
		users = append(users, &u)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("error iterating user rows: %w", rows.Err())
	}

	return users, nil
}