package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/bertoxic/graphqlChat/internal/drivers"
	errorx "github.com/bertoxic/graphqlChat/internal/error"
	"github.com/bertoxic/graphqlChat/internal/models"
	"github.com/bertoxic/graphqlChat/pkg/config"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PostgresDBRepo struct {
	App *config.AppConfig
	DB  *pgxpool.Pool
}

func NewPostgresDBRepo(a *config.AppConfig, db *drivers.PostgresDB) *PostgresDBRepo {
	return &PostgresDBRepo{
		App: a,
		DB:  db.Pool,
	}
}
func (pr *PostgresDBRepo) CreateUser(ctx context.Context, user models.RegistrationInput) (*models.UserDetails, error) {
	// Start the transaction
	tx, err := pr.DB.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		// Rollback in case of failure
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	// Private function to hide query details
	userDetails, err := pr.createUserTx(ctx, tx, user)
	if err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	// Commit the transaction if no errors
	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return userDetails, nil
}

// GetUserByEmail fetches a user by email with transaction management
func (pr *PostgresDBRepo) GetUserByEmail(ctx context.Context, email string) (*models.UserDetails, error) {
	// Start the transaction
	tx, err := pr.DB.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		// Rollback in case of failure
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	// Private function to hide query details
	userDetails, err := pr.getUserByEmailTx(ctx, tx, email)
	if err != nil {
		return nil, fmt.Errorf("error fetching user by email: %w", err)
	}

	// Commit the transaction if no errors
	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return userDetails, nil
}

// Private function that performs the SQL query for creating a user within a transaction
func (pr *PostgresDBRepo) createUserTx(ctx context.Context, tx pgx.Tx, user models.RegistrationInput) (*models.UserDetails, error) {
	query := `
		INSERT INTO users (username, email, password) 
		VALUES ($1, $2, $3) 
		RETURNING   id, username, email, password;
	`
	userNew := user

	// Struct to hold the created user's details
	var userDetails models.UserDetails

	// Execute the query within the transaction
	err := tx.QueryRow(ctx, query, userNew.Username, userNew.Email, userNew.Password).Scan(
		&userDetails.ID,
		&userDetails.UserName,
		&userDetails.Email,
		&userDetails.Password,
	)
	if err != nil {
		// Handle unique violation (email already exists)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // PostgreSQL unique violation code
			return nil, fmt.Errorf("user with email %s already exists", userNew.Email)
		}
		return nil, fmt.Errorf("error executing query: %w", err)
	}

	return &userDetails, nil
}

// Private function that performs the SQL query for fetching a user by email within a transaction
func (pr *PostgresDBRepo) getUserByEmailTx(ctx context.Context, tx pgx.Tx, email string) (*models.UserDetails, error) {
	query := `
		SELECT id, username, email, password 
		FROM users 
		WHERE email = $1;
	`

	// Struct to hold the fetched user's details
	var userDetails models.UserDetails

	// Execute the query within the transaction
	err := tx.QueryRow(ctx, query, email).Scan(
		&userDetails.ID,
		&userDetails.UserName,
		&userDetails.Email,
		&userDetails.Password,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("no user found with email: %s, %w", email, errorx.ErrNotFound)
		}
		return nil, fmt.Errorf("error executing query: %w", err)
	}

	return &userDetails, nil
}
