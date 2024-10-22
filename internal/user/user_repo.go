package user

import (
	"context"
	"fmt"
	"github.com/bertoxic/graphqlChat/internal/database"
	"github.com/bertoxic/graphqlChat/internal/database/postgres"
	errorx "github.com/bertoxic/graphqlChat/internal/error"
	"github.com/bertoxic/graphqlChat/internal/models"
	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

type UserRepo interface {
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	GetUsersFollowers(ctx context.Context, userID string, limit, offset int) ([]models.User, error)
	UpdateUserDetails(ctx context.Context, userDetails models.UpdateUserInput, userID string) (*models.UserDetails, error)
	GetUserStats(ctx context.Context, userID string) (*models.UserStats, error)
	GetUsersFollowing(ctx context.Context, userID string, limit, offset int) ([]models.User, error)
	FollowUser(ctx context.Context, followerID, followedID string) (*models.UserResponse, error)
	UnfollowUser(ctx context.Context, followerID, followedID string) (*models.UserResponse, error)
	GetUserNotifications(ctx context.Context, userID string, limit, offset int) ([]models.Notification, error)
	MarkNotificationAsRead(ctx context.Context, notificationID string) error
	SearchUsers(ctx context.Context, query string, limit int) ([]models.User, error)
	GetSuggestedUsers(ctx context.Context, userID string, limit int) ([]models.User, error)
	CheckUsernameAvailability(ctx context.Context, username string) (bool, error)
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) (*models.UserResponse, error)
	DeleteUser(ctx context.Context, userID string) (*models.UserResponse, error)
	ChangePassword(ctx context.Context, userID, currentPassword, newPassword string) (*models.UserResponse, error)
	GetUserLikedPosts(ctx context.Context, userID string, limit, offset int) ([]models.Post, error)
	GetUserBookmarkedPosts(ctx context.Context, userID string, limit, offset int) ([]models.Post, error)
	BookmarkPost(ctx context.Context, userID, postID string) (*models.UserResponse, error)
	RemoveBookmark(ctx context.Context, userID, postID string) (*models.UserResponse, error)
	GetUserFeed(ctx context.Context, userID string, limit, offset int) ([]models.Post, error)
	GetTrendingPosts(ctx context.Context, limit int) ([]models.Post, error)
}

type Repository struct {
	DB database.DatabaseRepo
}

func NewUserRepo(db database.DatabaseRepo) *Repository {
	return &Repository{
		DB: db,
	}
}

func (us *Repository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	db, ok := us.DB.(*postgres.PostgresDBRepo)
	if !ok {
		return nil, errorx.New(errorx.ErrCodeDatabase, "", errorx.ErrDatabase)
	}

	query := `
        SELECT id, username, email, full_name, bio, date_of_birth, profile_picture_url, cover_picture_url, 
               location, website, is_private, created_at, updated_at
        FROM users 
        WHERE email = $1
    `

	var user models.User
	err := db.DB.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.UserName, &user.Email, &user.FullName, &user.Bio, &user.DateOfBirth,
		&user.ProfilePictureURL, &user.CoverPictureURL, &user.Location, &user.Website,
		&user.IsPrivate, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}
func (us *Repository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	db, ok := us.DB.(*postgres.PostgresDBRepo)
	if !ok {
		return nil, errorx.New(errorx.ErrCodeDatabase, "", errorx.ErrDatabase)
	}

	query := `
        SELECT id, username, email, full_name, bio, date_of_birth, profile_picture_url, cover_picture_url, 
               location, website, is_private, created_at, updated_at
        FROM users 
        WHERE username = $1
    `

	var user models.User
	err := db.DB.QueryRow(ctx, query, username).Scan(
		&user.ID, &user.UserName, &user.Email, &user.FullName, &user.Bio, &user.DateOfBirth,
		&user.ProfilePictureURL, &user.CoverPictureURL, &user.Location, &user.Website,
		&user.IsPrivate, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}

	return &user, nil
}
func (us *Repository) GetUsersFollowers(ctx context.Context, userID string, limit, offset int) ([]*models.User, error) {
	db, ok := us.DB.(*postgres.PostgresDBRepo)
	if !ok {
		return nil, errorx.New(errorx.ErrCodeDatabase, "", errorx.ErrDatabase)
	}

	query := `
        SELECT u.id, u.username, u.email, u.full_name, u.bio, u.profile_picture_url
        FROM users u
        INNER JOIN follows f ON u.id = f.follower_id
        WHERE f.followed_id = $1
        ORDER BY f.created_at DESC
        LIMIT $2 OFFSET $3
    `

	rows, err := db.DB.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get user's followers: %w", err)
	}
	defer rows.Close()

	var followers []*models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.UserName, &user.Email, &user.FullName, &user.Bio, &user.ProfilePictureURL)
		if err != nil {
			return nil, fmt.Errorf("failed to scan follower: %w", err)
		}
		followers = append(followers, &user)
	}

	return followers, nil
}
func (us *Repository) UpdateUserDetails(ctx context.Context, userDetails models.UpdateUserInput, userID string) (*models.UserDetails, error) {
	db, ok := us.DB.(*postgres.PostgresDBRepo)
	if !ok {
		return &models.UserDetails{}, errorx.New(errorx.ErrCodeDatabase, "", errorx.ErrDatabase)
	}

	tx, err := db.DB.Begin(ctx)
	if err != nil {
		return &models.UserDetails{}, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	newUserDetails, err := updateUserDetails(ctx, tx, userDetails, userID)
	if err != nil {
		return &models.UserDetails{}, err
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return newUserDetails, nil
}
func updateUserDetails(ctx context.Context, tx pgx.Tx, details models.UpdateUserInput, userID string) (*models.UserDetails, error) {
	if details.UpdatedAt == nil {
		now := time.Now()
		details.UpdatedAt = &now
	}
	query := `
       UPDATE users 
		SET username = COALESCE($1, username),
		    email = COALESCE($2, email),
		    full_name = COALESCE($3, full_name),
		    bio = COALESCE($4, bio),
		    date_of_birth = COALESCE($5, date_of_birth),
		    profile_picture_url = COALESCE($6, profile_picture_url),
		    cover_picture_url = COALESCE($7, cover_picture_url),
		    location = COALESCE($8, location),
		    website = COALESCE($9, website),
		    is_private = COALESCE($10, is_private),
		    updated_at = $11
		WHERE id = $12
		RETURNING id, username, email, full_name, bio, date_of_birth, profile_picture_url, cover_picture_url, location, website, is_private, updated_at
	`

	var updatedDetails models.UserDetails
	err := tx.QueryRow(ctx, query,
		details.UserName, details.Email, details.FullName, details.Bio,
		details.DateOfBirth, details.ProfilePictureURL, details.CoverPictureURL,
		details.Location, details.Website, details.IsPrivate, details.UpdatedAt,
		userID,
	).Scan(
		&updatedDetails.ID, &updatedDetails.UserName, &updatedDetails.Email,
		&updatedDetails.FullName, &updatedDetails.Bio, &updatedDetails.DateOfBirth,
		&updatedDetails.ProfilePictureURL, &updatedDetails.CoverPictureURL,
		&updatedDetails.Location, &updatedDetails.Website, &updatedDetails.IsPrivate,
		&updatedDetails.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update user details: %w", err)
	}

	return &updatedDetails, nil
}
func (us *Repository) GetUserStats(ctx context.Context, userID string) (*models.UserStats, error) {
	db, ok := us.DB.(*postgres.PostgresDBRepo)
	if !ok {
		return nil, errorx.New(errorx.ErrCodeDatabase, "", errorx.ErrDatabase)
	}

	query := `
        SELECT 
            (SELECT COUNT(*) FROM posts WHERE user_id = $1) as total_posts,
            (SELECT COUNT(*) FROM follows WHERE followed_id = $1) as total_followers,
            (SELECT COUNT(*) FROM follows WHERE follower_id = $1) as total_following
    `

	var stats models.UserStats
	err := db.DB.QueryRow(ctx, query, userID).Scan(&stats.TotalPosts, &stats.TotalFollowers, &stats.TotalFollowing)
	if err != nil {
		return nil, fmt.Errorf("failed to get user stats: %w", err)
	}

	return &stats, nil
}
func xupdateUserDetails(ctx context.Context, tx pgx.Tx, details models.UpdateUserInput, userID string) (*models.UserDetails, error) {
	// Map to store the column-value pairs that need to be updated
	updates := map[string]interface{}{}

	// Add only non-empty fields to the updates map
	if *details.UserName != "" {
		updates["username"] = details.UserName
	}
	if *details.Email != "" {
		updates["email"] = details.Email
	}
	if *details.FullName != "" {
		updates["full_name"] = details.FullName
	}
	if *details.Bio != "" {
		updates["bio"] = details.Bio
	}
	if *details.DateOfBirth != "" {
		updates["date_of_birth"] = details.DateOfBirth
	}
	if *details.ProfilePictureURL != "" {
		updates["profile_picture_url"] = details.ProfilePictureURL
	}
	if *details.CoverPictureURL != "" {
		updates["cover_picture_url"] = details.CoverPictureURL
	}
	if *details.Location != "" {
		updates["location"] = details.Location
	}
	if *details.Website != "" {
		updates["website"] = details.Website
	}
	if details.IsPrivate != nil {
		updates["is_private"] = *details.IsPrivate
	}

	// Always update updated_at field
	updates["updated_at"] = details.UpdatedAt

	// Construct the SQL query dynamically
	var query string
	query = "UPDATE users SET "
	params := []interface{}{}
	i := 1

	for col, val := range updates {
		if i > 1 {
			query += ", "
		}
		query += fmt.Sprintf("%s = $%d", col, i)
		params = append(params, val)
		i++
	}

	query += fmt.Sprintf(" WHERE id = $%d RETURNING id, username, email, full_name, bio, date_of_birth, profile_picture_url, cover_picture_url, location, website, is_private, updated_at", i)
	params = append(params, userID)

	// Execute the query
	var updatedDetails models.UserDetails
	err := tx.QueryRow(ctx, query, params...).Scan(
		&updatedDetails.ID, &updatedDetails.UserName, &updatedDetails.Email,
		&updatedDetails.FullName, &updatedDetails.Bio, &updatedDetails.DateOfBirth,
		&updatedDetails.ProfilePictureURL, &updatedDetails.CoverPictureURL,
		&updatedDetails.Location, &updatedDetails.Website, &updatedDetails.IsPrivate,
		&updatedDetails.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update user details: %w", err)
	}

	return &updatedDetails, nil
}
func (us *Repository) GetUsersFollowing(ctx context.Context, userID string, limit, offset int) ([]*models.User, error) {
	db, ok := us.DB.(*postgres.PostgresDBRepo)
	if !ok {
		return nil, errorx.New(errorx.ErrCodeDatabase, "", errorx.ErrDatabase)
	}

	query := `
        SELECT u.id, u.username, u.email, u.full_name, u.bio, u.profile_picture_url
        FROM users u
        INNER JOIN follows f ON u.id = f.followed_id
        WHERE f.follower_id = $1
        ORDER BY f.created_at DESC
        LIMIT $2 OFFSET $3
    `

	rows, err := db.DB.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get user's following: %w", err)
	}
	defer rows.Close()

	var following []*models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.UserName, &user.Email, &user.FullName, &user.Bio, &user.ProfilePictureURL)
		if err != nil {
			return nil, fmt.Errorf("failed to scan following user: %w", err)
		}
		following = append(following, &user)
	}

	return following, nil
}
func (us *Repository) FollowUser(ctx context.Context, followerID, followedID string) (*models.UserResponse, error) {
	db, ok := us.DB.(*postgres.PostgresDBRepo)
	if !ok {
		return nil, errorx.New(errorx.ErrCodeDatabase, "", errorx.ErrDatabase)
	}

	query := `
        INSERT INTO follows (follower_id, followed_id)
        VALUES ($1, $2)
        ON CONFLICT (follower_id, followed_id) DO NOTHING
    `

	_, err := db.DB.Exec(ctx, query, followerID, followedID)
	if err != nil {
		return nil, fmt.Errorf("failed to follow user: %w", err)
	}

	return &models.UserResponse{
		Success: true,
		Message: "Successfully followed user",
	}, nil
}
func (us *Repository) UnfollowUser(ctx context.Context, followerID, followedID string) (*models.UserResponse, error) {
	db, ok := us.DB.(*postgres.PostgresDBRepo)
	if !ok {
		return nil, errorx.New(errorx.ErrCodeDatabase, "", errorx.ErrDatabase)
	}

	query := `
        DELETE FROM follows
        WHERE follower_id = $1 AND followed_id = $2
    `

	_, err := db.DB.Exec(ctx, query, followerID, followedID)
	if err != nil {
		return nil, fmt.Errorf("failed to unfollow user: %w", err)
	}

	return &models.UserResponse{
		Success: true,
		Message: "Successfully unfollowed user",
	}, nil
}
func (us *Repository) GetUserNotifications(ctx context.Context, userID string, limit, offset int) ([]*models.Notification, error) {
	db, ok := us.DB.(*postgres.PostgresDBRepo)
	if !ok {
		return nil, errorx.New(errorx.ErrCodeDatabase, "", errorx.ErrDatabase)
	}

	query := `
        SELECT id, user_id, type, content, is_read, created_at
        FROM notifications
        WHERE user_id = $1
        ORDER BY created_at DESC
        LIMIT $2 OFFSET $3
    `

	rows, err := db.DB.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get user notifications: %w", err)
	}
	defer rows.Close()

	var notifications []*models.Notification
	for rows.Next() {
		var notification models.Notification
		err := rows.Scan(&notification.ID, &notification.UserID, &notification.Type, &notification.Content, &notification.IsRead, &notification.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan notification: %w", err)
		}
		notifications = append(notifications, &notification)
	}

	return notifications, nil
}
func (us *Repository) MarkNotificationAsRead(ctx context.Context, notificationID string) error {
	db, ok := us.DB.(*postgres.PostgresDBRepo)
	if !ok {
		return errorx.New(errorx.ErrCodeDatabase, "", errorx.ErrDatabase)
	}

	query := `
        UPDATE notifications
        SET is_read = true
        WHERE id = $1
    `

	_, err := db.DB.Exec(ctx, query, notificationID)
	if err != nil {
		return fmt.Errorf("failed to mark notification as read: %w", err)
	}

	return nil
}
func (us *Repository) SearchUsers(ctx context.Context, query string, limit int) ([]*models.User, error) {
	db, ok := us.DB.(*postgres.PostgresDBRepo)
	if !ok {
		return nil, errorx.New(errorx.ErrCodeDatabase, "", errorx.ErrDatabase)
	}

	sqlQuery := `
        SELECT id, username, email, full_name, bio, profile_picture_url
        FROM users
        WHERE username ILIKE $1 OR full_name ILIKE $1
        LIMIT $2
    `

	rows, err := db.DB.Query(ctx, sqlQuery, "%"+query+"%", limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
	}
	defer rows.Close()

	var users = []*models.User{}
	for rows.Next() {
		var user = &models.User{}
		err := rows.Scan(&user.ID, &user.UserName, &user.Email, &user.FullName, &user.Bio, &user.ProfilePictureURL)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user in search results: %w", err)
		}

		if user.ID == "" || user.UserName == "" || user.Email == "" {
			log.Printf("User object has missing fields: %+v", user)
		}
		users = append(users, user)

		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("error iterating over rows: %w", err)
		}

	}

	return users, nil
}
func (us *Repository) GetSuggestedUsers(ctx context.Context, userID string, limit int) ([]*models.User, error) {
	db, ok := us.DB.(*postgres.PostgresDBRepo)
	if !ok {
		return nil, errorx.New(errorx.ErrCodeDatabase, "", errorx.ErrDatabase)
	}

	query := `
        SELECT u.id, u.username, u.email, u.full_name, u.bio, u.profile_picture_url
        FROM users u
        WHERE u.id NOT IN (
            SELECT followed_id FROM follows WHERE follower_id = $1
        ) AND u.id != $1
        ORDER BY RANDOM()
        LIMIT $2
    `

	rows, err := db.DB.Query(ctx, query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get suggested users: %w", err)
	}
	defer rows.Close()

	var suggestedUsers []*models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.UserName, &user.Email, &user.FullName, &user.Bio, &user.ProfilePictureURL)
		if err != nil {
			return nil, fmt.Errorf("failed to scan suggested user: %w", err)
		}
		suggestedUsers = append(suggestedUsers, &user)
	}

	return suggestedUsers, nil
}
func (us *Repository) CheckUsernameAvailability(ctx context.Context, username string) (bool, error) {
	db, ok := us.DB.(*postgres.PostgresDBRepo)
	if !ok {
		return false, errorx.New(errorx.ErrCodeDatabase, "", errorx.ErrDatabase)
	}

	query := `
        SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)
    `

	var exists bool
	err := db.DB.QueryRow(ctx, query, username).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check username availability: %w", err)
	}

	return !exists, nil
}
func (us *Repository) GetUserByID(ctx context.Context, id string) (models.User, error) {
	db, ok := us.DB.(*postgres.PostgresDBRepo)
	if !ok {
		return models.User{}, errorx.New(errorx.ErrCodeDatabase, "", errorx.ErrDatabase)
	}

	query := `
        SELECT id, username, email, full_name, bio, profile_picture_url, 
               cover_picture_url, location, website, is_private, created_at, updated_at
        FROM users
        WHERE id = $1
    `

	var user models.User
	err := db.DB.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.UserName, &user.Email, &user.FullName, &user.Bio, &user.ProfilePictureURL,
		&user.CoverPictureURL, &user.Location, &user.Website, &user.IsPrivate, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return models.User{}, fmt.Errorf("user not found: %w", err)
		}
		return models.User{}, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return user, nil
}

//func (us *Repository) BlockUser(ctx context.Context, blockerID, blockedID string) (*models.UserResponse, error) {
//	db, ok := us.DB.(*postgres.PostgresDBRepo)
//	if !ok {
//		return nil, errorx.New(errorx.ErrCodeDatabase, "", errorx.ErrDatabase)
//	}
//
//	tx, err := db.DB.Begin(ctx)
//	if err != nil {
//		return nil, fmt.Errorf("failed to begin transaction: %w", err)
//	}
//	defer tx.Rollback(ctx)
//
//	// First, unfollow each other if they were following
//	unfollowQuery := `
//        DELETE FROM followers
//        WHERE (follower_id = $1 AND followed_id = $2)
//           OR (follower_id = $2 AND followed_id = $1)
//    `
//	_, err = tx.Exec(ctx, unfollowQuery, blockerID, blockedID)
//	if err != nil {
//		return nil, fmt.Errorf("failed to unfollow users: %w", err)
//	}
//
//	// Then, add to blocked users
//	blockQuery := `
//        INSERT INTO blocked_users (blocker_id, blocked_id)
//        VALUES ($1, $2)
//        ON CONFLICT (blocker_id, blocked_id) DO NOTHING
//    `
//	_, err = tx.Exec(ctx, blockQuery, blockerID, blockedID)
//	if err != nil {
//		return nil, fmt.Errorf("failed to block user: %w", err)
//	}
//
//	if err = tx.Commit(ctx); err != nil {
//		return nil, fmt.Errorf("failed to commit transaction: %w", err)
//	}
//
//	return &models.UserResponse{
//		Success: true,
//		Message: "Successfully blocked user",
//	}, nil
//}

//func (us *Repository) UnblockUser(ctx context.Context, unblockerID, unblockedID string) (*models.UserResponse, error) {
//	db, ok := us.DB.(*postgres.PostgresDBRepo)
//	if !ok {
//		return nil, errorx.New(errorx.ErrCodeDatabase, "", errorx.ErrDatabase)
//	}
//
//	query := `
//        DELETE FROM blocked_users
//        WHERE blocker_id = $1 AND blocked_id = $2
//    `
//
//	_, err := db.DB.Exec(ctx, query, unblockerID, unblockedID)
//	if err != nil {
//		return nil, fmt.Errorf("failed to unblock user: %w", err)
//	}
//
//	return &models.UserResponse{
//		Success: true,
//		Message: "Successfully unblocked user",
//	}, nil
//}

// ... [Previous code remains the same] ....................................................>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

func (us *Repository) CreateUser(ctx context.Context, user *models.User) (*models.UserResponse, error) {
	db, ok := us.DB.(*postgres.PostgresDBRepo)
	if !ok {
		return nil, errorx.New(errorx.ErrCodeDatabase, "", errorx.ErrDatabase)
	}

	query := `
        INSERT INTO users (username, email, full_name, bio, date_of_birth, profile_picture_url, cover_picture_url, 
                           location, website, is_private, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $11)
        RETURNING id
    `

	var userID string
	err := db.DB.QueryRow(ctx, query,
		user.UserName, user.Email, user.FullName, user.Bio, user.DateOfBirth,
		user.ProfilePictureURL, user.CoverPictureURL, user.Location, user.Website,
		user.IsPrivate, time.Now(),
	).Scan(&userID)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	user.ID = userID

	return &models.UserResponse{
		Success: true,
		Message: "User created successfully",
		Data:    map[string]interface{}{"user": user},
	}, nil
}
func (us *Repository) DeleteUser(ctx context.Context, userID string) (*models.UserResponse, error) {
	db, ok := us.DB.(*postgres.PostgresDBRepo)
	if !ok {
		return nil, errorx.New(errorx.ErrCodeDatabase, "", errorx.ErrDatabase)
	}

	tx, err := db.DB.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Delete user's posts
	_, err = tx.Exec(ctx, "DELETE FROM posts WHERE user_id = $1", userID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete user's posts: %w", err)
	}

	// Delete user's followers and following relationships
	_, err = tx.Exec(ctx, "DELETE FROM follows WHERE follower_id = $1 OR followed_id = $1", userID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete user's follower relationships: %w", err)
	}

	// Delete user's notifications
	_, err = tx.Exec(ctx, "DELETE FROM notifications WHERE user_id = $1", userID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete user's notifications: %w", err)
	}

	// Finally, delete the user
	_, err = tx.Exec(ctx, "DELETE FROM users WHERE id = $1", userID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete user: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &models.UserResponse{
		Success: true,
		Message: "User deleted successfully",
	}, nil
}
func (us *Repository) ResetPassword(ctx context.Context, userID, currentPassword, newPassword string) (*models.UserResponse, error) {
	db, ok := us.DB.(*postgres.PostgresDBRepo)
	if !ok {
		return nil, errorx.New(errorx.ErrCodeDatabase, "", errorx.ErrDatabase)
	}

	// First, verify the current password
	var storedHash string
	err := db.DB.QueryRow(ctx, "SELECT password FROM users WHERE id = $1", userID).Scan(&storedHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get user's password hash: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(currentPassword)); err != nil {
		return &models.UserResponse{
			Success: false,
			Message: "Current password is incorrect",
		}, nil
	}

	// Hash the new password
	newHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash new password: %w", err)
	}

	// Update the password
	_, err = db.DB.Exec(ctx, "UPDATE users SET password = $1 WHERE id = $2", newHash, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to update password: %w", err)
	}

	return &models.UserResponse{
		Success: true,
		Message: "Password changed successfully",
	}, nil
}
func (us *Repository) GetUserLikedPosts(ctx context.Context, userID string, limit, offset int) ([]*models.Post, error) {
	db, ok := us.DB.(*postgres.PostgresDBRepo)
	if !ok {
		return nil, errorx.New(errorx.ErrCodeDatabase, "", errorx.ErrDatabase)
	}

	query := `
        SELECT p.id, p.user_id, p.title, p.content, p.image_url, p.video_url, p.audio_url, 
               p.is_edited, p.is_draft, p.parent_id, p.created_at, p.updated_at, p.likes, p.reposts
        FROM posts p
        INNER JOIN post_likes l ON p.id = l.post_id
        WHERE l.user_id = $1
        ORDER BY l.created_at DESC
        LIMIT $2 OFFSET $3
    `

	rows, err := db.DB.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get user's liked posts: %w", err)
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(
			&post.ID, &post.UserID, &post.Title, &post.Content, &post.ImageURL, &post.VideoURL, &post.AudioURL,
			&post.IsEdited, &post.IsDraft, &post.ParentID, &post.CreatedAt, &post.UpdatedAt, &post.Likes, &post.Reposts,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan liked post: %w", err)
		}
		posts = append(posts, &post)
	}

	return posts, nil
}

//	func (us *Repository) ReportUser(ctx context.Context, reporterID, reportedID, reason string) (*models.UserResponse, error) {
//		db, ok := us.DB.(*postgres.PostgresDBRepo)
//		if !ok {
//			return nil, errorx.New(errorx.ErrCodeDatabase, "", errorx.ErrDatabase)
//		}
//
//		query := `
//	       INSERT INTO user_reports (reporter_id, reported_id, reason, created_at)
//	       VALUES ($1, $2, $3, $4)
//	   `
//
//		_, err := db.DB.Exec(ctx, query, reporterID, reportedID, reason, time.Now())
//		if err != nil {
//			return nil, fmt.Errorf("failed to report user: %w", err)
//		}
//
//		return &models.UserResponse{
//			Success: true,
//			Message: "User reported successfully",
//		}, nil
//	}
func (us *Repository) GetUserBookmarkedPosts(ctx context.Context, userID string, limit, offset int) ([]*models.Post, error) {
	db, ok := us.DB.(*postgres.PostgresDBRepo)
	if !ok {
		return nil, errorx.New(errorx.ErrCodeDatabase, "", errorx.ErrDatabase)
	}

	query := `
        SELECT p.id, p.user_id, p.title, p.content, p.image_url, p.video_url, p.audio_url, 
               p.is_edited, p.is_draft, p.parent_id, p.created_at, p.updated_at, p.likes, p.reposts
        FROM posts p
        INNER JOIN bookmarks b ON p.id = b.post_id
        WHERE b.user_id = $1
        ORDER BY b.created_at DESC
        LIMIT $2 OFFSET $3
    `

	rows, err := db.DB.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get user's bookmarked posts: %w", err)
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(
			&post.ID, &post.UserID, &post.Title, &post.Content, &post.ImageURL, &post.VideoURL, &post.AudioURL,
			&post.IsEdited, &post.IsDraft, &post.ParentID, &post.CreatedAt, &post.UpdatedAt, &post.Likes, &post.Reposts,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan bookmarked post: %w", err)

		}
		posts = append(posts, &post)
	}

	return posts, nil
}
func (us *Repository) BookmarkPost(ctx context.Context, userID, postID string) (*models.UserResponse, error) {
	db, ok := us.DB.(*postgres.PostgresDBRepo)
	if !ok {
		return nil, errorx.New(errorx.ErrCodeDatabase, "", errorx.ErrDatabase)
	}

	query := `
        INSERT INTO bookmarks (user_id, post_id, created_at)
        VALUES ($1, $2, $3)
        ON CONFLICT (user_id, post_id) DO NOTHING
    `

	_, err := db.DB.Exec(ctx, query, userID, postID, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to bookmark post: %w", err)
	}

	return &models.UserResponse{
		Success: true,
		Message: "Post bookmarked successfully",
	}, nil
}

func (us *Repository) RemoveBookmark(ctx context.Context, userID, postID string) (*models.UserResponse, error) {
	db, ok := us.DB.(*postgres.PostgresDBRepo)
	if !ok {
		return nil, errorx.New(errorx.ErrCodeDatabase, "", errorx.ErrDatabase)
	}

	query := `
        DELETE FROM bookmarks
        WHERE user_id = $1 AND post_id = $2
    `

	_, err := db.DB.Exec(ctx, query, userID, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to remove bookmark: %w", err)
	}

	return &models.UserResponse{
		Success: true,
		Message: "Bookmark removed successfully",
	}, nil
}

func (us *Repository) GetUserFeed(ctx context.Context, userID string, limit, offset int) ([]*models.Post, error) {
	db, ok := us.DB.(*postgres.PostgresDBRepo)
	if !ok {
		return nil, errorx.New(errorx.ErrCodeDatabase, "", errorx.ErrDatabase)
	}

	query := `
        SELECT p.id, p.user_id, p.title, p.content, p.image_url, p.video_url, p.audio_url, 
               p.is_edited, p.is_draft, p.parent_id, p.created_at, p.updated_at, p.likes, p.reposts
        FROM posts p
        INNER JOIN follows f ON p.user_id = f.followed_id
        WHERE f.follower_id = $1
        ORDER BY p.created_at DESC
        LIMIT $2 OFFSET $3
    `

	rows, err := db.DB.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get user's feed: %w", err)
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(
			&post.ID, &post.UserID, &post.Title, &post.Content, &post.ImageURL, &post.VideoURL, &post.AudioURL,
			&post.IsEdited, &post.IsDraft, &post.ParentID, &post.CreatedAt, &post.UpdatedAt, &post.Likes, &post.Reposts,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan feed post: %w", err)
		}
		posts = append(posts, &post)
	}

	return posts, nil
}

func (us *Repository) GetTrendingPosts(ctx context.Context, limit int) ([]*models.Post, error) {
	db, ok := us.DB.(*postgres.PostgresDBRepo)
	if !ok {
		return nil, errorx.New(errorx.ErrCodeDatabase, "", errorx.ErrDatabase)
	}

	query := `
        SELECT id, user_id, title, content, image_url, video_url, audio_url, 
               is_edited, is_draft, parent_id, created_at, updated_at, likes, reposts
        FROM posts
        WHERE created_at > $1
        ORDER BY (likes + reposts) DESC
        LIMIT $2
    `

	// Get posts from the last 24 hours
	timeThreshold := time.Now().Add(-24 * time.Hour)

	rows, err := db.DB.Query(ctx, query, timeThreshold, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get trending posts: %w", err)
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(
			&post.ID, &post.UserID, &post.Title, &post.Content, &post.ImageURL, &post.VideoURL, &post.AudioURL,
			&post.IsEdited, &post.IsDraft, &post.ParentID, &post.CreatedAt, &post.UpdatedAt, &post.Likes, &post.Reposts,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan trending post: %w", err)
		}
		posts = append(posts, &post)
	}

	return posts, nil
}
