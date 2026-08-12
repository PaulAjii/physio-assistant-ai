package repo

import (
	"context"
	"errors"

	"github.com/PaulAjii/physio-assistant-ai/internal/db"
	"github.com/PaulAjii/physio-assistant-ai/internal/models"
	"github.com/jackc/pgx/v5"
)

// CreateUserParams is the input for creating a user. The clinic is passed
// separately (it is both the tenant context and the row's clinic_id), and the
// password is already hashed by the caller — the repo never sees plaintext.
type CreateUserParams struct {
	Email        string
	PasswordHash string
	FullName     string
	Role         string
	LicenseNo    *string
}

// Users is the persistence contract for user rows, so callers can depend on the
// interface and swap in a mock. Login and refresh do NOT go through here: they
// run before a tenant context exists and use the SECURITY DEFINER functions.
type Users interface {
	Create(ctx context.Context, clinicID string, p CreateUserParams) (*models.User, error)
	GetByID(ctx context.Context, clinicID, id string) (*models.User, error)
}

// UserRepo is the Postgres-backed implementation of Users.
type UserRepo struct {
	db *db.DB
}

// NewUserRepo returns a UserRepo backed by the given database.
func NewUserRepo(database *db.DB) *UserRepo {
	return &UserRepo{db: database}
}

// Compile-time check that UserRepo satisfies Users.
var _ Users = (*UserRepo)(nil)

const createUserSQL = `
INSERT INTO users (clinic_id, email, password_hash, full_name, role, license_no)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, clinic_id, email, full_name, role, license_no, password_hash, created_at, updated_at`

// Create inserts a user into clinicID. The clinic_id written equals the tenant
// context set by WithTenant, so the row satisfies the RLS WITH CHECK. A
// duplicate live email surfaces as ErrConflict rather than a raw pgx error.
func (r *UserRepo) Create(ctx context.Context, clinicID string, p CreateUserParams) (*models.User, error) {
	var u models.User
	err := r.db.WithTenant(ctx, clinicID, func(tx pgx.Tx) error {
		return tx.QueryRow(ctx, createUserSQL,
			clinicID, p.Email, p.PasswordHash, p.FullName, p.Role, p.LicenseNo,
		).Scan(
			&u.ID, &u.ClinicID, &u.Email, &u.FullName, &u.Role,
			&u.LicenseNo, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt,
		)
	})
	if err != nil {
		if isUniqueViolation(err) {
			return nil, ErrConflict
		}
		return nil, err
	}
	return &u, nil
}

const getUserByIDSQL = `
SELECT id, clinic_id, email, full_name, role, license_no, password_hash, created_at, updated_at
FROM users
WHERE id = $1 AND deleted_at IS NULL`

// GetByID fetches a live user by id within clinicID. RLS also confines the query
// to that clinic, so a valid id from another tenant simply returns ErrNotFound.
func (r *UserRepo) GetByID(ctx context.Context, clinicID, id string) (*models.User, error) {
	var u models.User
	err := r.db.WithTenant(ctx, clinicID, func(tx pgx.Tx) error {
		return tx.QueryRow(ctx, getUserByIDSQL, id).Scan(
			&u.ID, &u.ClinicID, &u.Email, &u.FullName, &u.Role,
			&u.LicenseNo, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt,
		)
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}
