package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/bertoxic/graphqlChat/internal/models"
)

type Service struct {
	UserRepo *Repository
}

func NewService(userRepo *Repository) *Service {
	return &Service{UserRepo: userRepo}
}

type UserService interface {
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	GetUsersFollowers(ctx context.Context, userID string, limit, offset int) ([]*models.User, error)
	GetUsersFollowing(ctx context.Context, userID string, limit, offset int) ([]*models.User, error)
	UpdateUserDetails(ctx context.Context, userDetails models.UpdateUserInput, userID string) (*models.UserDetails, error)
	FollowUser(ctx context.Context, followerID, followedID string) (*models.UserResponse, error)
	UnfollowUser(ctx context.Context, followerID, followedID string) (*models.UserResponse, error)
	GetUserStats(ctx context.Context, userID string) (*models.UserStats, error)
	GetUserNotifications(ctx context.Context, userID string, limit, offset int) ([]*models.Notification, error)
	MarkNotificationAsRead(ctx context.Context, notificationID string) error
	GetSuggestedUsers(ctx context.Context, userID string, limit int) ([]*models.User, error)
	GetUserByID(ctx context.Context, id string) (models.User, error)
	CheckUsernameAvailability(ctx context.Context, username string) (bool, error)
	SearchUsers(ctx context.Context, query string, limit int) ([]*models.User, error)
	CreateUser(ctx context.Context, user *models.User) (*models.UserResponse, error)
	DeleteUser(ctx context.Context, userID string) (*models.UserResponse, error)
	ResetPassword(ctx context.Context, userID, currentPassword, newPassword string) (*models.UserResponse, error)
	GetUserLikedPosts(ctx context.Context, userID string, limit, offset int) ([]*models.Post, error)
	GetUserBookmarkedPosts(ctx context.Context, userID string, limit, offset int) ([]*models.Post, error)
	BookmarkPost(ctx context.Context, userID, postID string) (*models.UserResponse, error)
	RemoveBookmark(ctx context.Context, userID, postID string) (*models.UserResponse, error)
	GetUserFeed(ctx context.Context, userID string, limit, offset int) ([]*models.Post, error)
	GetTrendingPosts(ctx context.Context, limit int) ([]*models.Post, error)
}

func (s Service) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := s.UserRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s Service) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	user, err := s.UserRepo.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s Service) UnfollowUser(ctx context.Context, followerID, followedID string) (*models.UserResponse, error) {
	resp, err := s.UserRepo.UnfollowUser(ctx, followerID, followedID)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s Service) GetUsersFollowers(ctx context.Context, userID string, limit, offset int) ([]*models.User, error) {
	followers, err := s.UserRepo.GetUsersFollowers(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	return followers, nil
}

func (s Service) GetUsersFollowing(ctx context.Context, userID string, limit, offset int) ([]*models.User, error) {
	following, err := s.UserRepo.GetUsersFollowing(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	return following, nil
}

func (s Service) FollowUser(ctx context.Context, followerID, followedID string) (*models.UserResponse, error) {
	if followerID == followedID {
		return nil, errors.New("users cannot follow themselves")
	}

	resp, err := s.UserRepo.FollowUser(ctx, followerID, followedID)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s Service) UpdateUserDetails(ctx context.Context, userDetails models.UpdateUserInput, userID string) (*models.UserDetails, error) {
	newUserDetails, err := s.UserRepo.UpdateUserDetails(ctx, userDetails, userID)
	if err != nil {
		return nil, err
	}
	return newUserDetails, nil
}

func (s Service) GetUserStats(ctx context.Context, userID string) (*models.UserStats, error) {
	stats, err := s.UserRepo.GetUserStats(ctx, userID)
	if err != nil {
		return nil, err
	}

	return stats, err
}

func (s Service) GetUserNotifications(ctx context.Context, userID string, limit, offset int) ([]*models.Notification, error) {
	notifications, err := s.UserRepo.GetUserNotifications(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get notifications for user %s: %w", userID, err)
	}

	if len(notifications) == 0 {
		return []*models.Notification{}, nil
	}

	return notifications, nil
}

func (s Service) MarkNotificationAsRead(ctx context.Context, notificationID string) error {
	err := s.UserRepo.MarkNotificationAsRead(ctx, notificationID)
	if err != nil {
		return fmt.Errorf("failed to mark notification %s as read: %w", notificationID, err)
	}

	return nil
}
func (s Service) GetSuggestedUsers(ctx context.Context, userID string, limit int) ([]*models.User, error) {
	suggestedUsers, err := s.UserRepo.GetSuggestedUsers(ctx, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get suggested users for user %s: %w", userID, err)
	}

	if len(suggestedUsers) == 0 {
		return []*models.User{}, nil
	}

	return suggestedUsers, nil
}
func (s Service) GetUserByID(ctx context.Context, id string) (models.User, error) {
	user, err := s.UserRepo.GetUserByID(ctx, id)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to get user with ID %s: %w", id, err)
	}

	return user, nil
}

func (s Service) CheckUsernameAvailability(ctx context.Context, username string) (bool, error) {
	exists, err := s.UserRepo.CheckUsernameAvailability(ctx, username)
	if err != nil {
		return false, fmt.Errorf("failed to check username availability for %s: %w", username, err)
	}

	return !exists, nil
}
func (s Service) SearchUsers(ctx context.Context, query string, limit int) ([]*models.User, error) {
	search, err := s.UserRepo.SearchUsers(ctx, query, limit)
	if err != nil {
		return nil, err
	}

	return search, nil
}

func (s Service) CreateUser(ctx context.Context, user *models.User) (*models.UserResponse, error) {
	response, err := s.UserRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return response, nil
}

func (s Service) DeleteUser(ctx context.Context, userID string) (*models.UserResponse, error) {
	response, err := s.UserRepo.DeleteUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete user with ID %s: %w", userID, err)
	}
	return response, nil
}
func (s Service) ResetPassword(ctx context.Context, userID, currentPassword, newPassword string) (*models.UserResponse, error) {
	response, err := s.UserRepo.ResetPassword(ctx, userID, currentPassword, newPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to change password for user with ID %s: %w", userID, err)
	}
	return response, nil
}
func (s Service) GetUserLikedPosts(ctx context.Context, userID string, limit, offset int) ([]*models.Post, error) {
	posts, err := s.UserRepo.GetUserLikedPosts(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get liked posts for user with ID %s: %w", userID, err)
	}
	return posts, nil
}
func (s Service) GetUserBookmarkedPosts(ctx context.Context, userID string, limit, offset int) ([]*models.Post, error) {
	posts, err := s.UserRepo.GetUserBookmarkedPosts(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get bookmarked posts for user with ID %s: %w", userID, err)
	}
	return posts, nil
}
func (s Service) BookmarkPost(ctx context.Context, userID, postID string) (*models.UserResponse, error) {
	response, err := s.UserRepo.BookmarkPost(ctx, userID, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to bookmark post %s for user %s: %w", postID, userID, err)
	}
	return response, nil
}
func (s Service) RemoveBookmark(ctx context.Context, userID, postID string) (*models.UserResponse, error) {
	response, err := s.UserRepo.RemoveBookmark(ctx, userID, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to remove bookmark for post %s from user %s: %w", postID, userID, err)
	}
	return response, nil
}
func (s Service) GetUserFeed(ctx context.Context, userID string, limit, offset int) ([]*models.Post, error) {
	posts, err := s.UserRepo.GetUserFeed(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get feed for user with ID %s: %w", userID, err)
	}
	return posts, nil
}
func (s Service) GetTrendingPosts(ctx context.Context, limit int) ([]*models.Post, error) {
	posts, err := s.UserRepo.GetTrendingPosts(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get trending posts: %w", err)
	}
	return posts, nil
}
