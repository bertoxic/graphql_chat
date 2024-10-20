package models

import (
	"time"
)

type UserDetails struct {
	ID                string     `json:"id"`                            // UUID
	UserName          string     `json:"username"`                      // Unique username
	Email             string     `json:"email"`                         // Unique email
	Password          string     `json:"password,omitempty"`            // Hashed password
	FullName          *string    `json:"full_name,omitempty"`           // Optional full name
	Bio               *string    `json:"bio,omitempty"`                 // Optional bio
	DateOfBirth       *time.Time `json:"date_of_birth,omitempty"`       // Optional date of birth
	ProfilePictureURL *string    `json:"profile_picture_url,omitempty"` // Optional profile picture URL
	CoverPictureURL   *string    `json:"cover_picture_url,omitempty"`   // Optional cover picture URL
	Location          *string    `json:"location,omitempty"`            // Optional location
	Website           *string    `json:"website,omitempty"`             // Optional website
	IsVerified        bool       `json:"is_verified"`                   // Verified account status
	IsPrivate         bool       `json:"is_private"`                    // Private account status
	FollowerCount     int        `json:"follower_count"`                // Number of followers
	FollowingCount    int        `json:"following_count"`               // Number of following accounts
	PostCount         int        `json:"post_count"`                    // Number of posts
	LastLogin         *time.Time `json:"last_login,omitempty"`          // Optional last login time
	CreatedAt         time.Time  `json:"created_at"`                    // Timestamp of account creation
	UpdatedAt         time.Time  `json:"updated_at"`                    // Timestamp of the last update
}

type User struct {
	ID                string     `json:"id"`                            // UUID
	UserName          string     `json:"username"`                      // Unique username
	Email             string     `json:"email"`                         // Unique email
	Password          string     `json:"password,omitempty"`            // Hashed password
	FullName          *string    `json:"full_name,omitempty"`           // Optional full name
	Bio               *string    `json:"bio,omitempty"`                 // Optional bio
	DateOfBirth       *time.Time `json:"date_of_birth,omitempty"`       // Optional date of birth
	ProfilePictureURL *string    `json:"profile_picture_url,omitempty"` // Optional profile picture URL
	CoverPictureURL   *string    `json:"cover_picture_url,omitempty"`   // Optional cover picture URL
	Location          *string    `json:"location,omitempty"`            // Optional location
	Website           *string    `json:"website,omitempty"`             // Optional website
	IsVerified        bool       `json:"is_verified"`                   // Verified account status
	IsPrivate         bool       `json:"is_private"`                    // Private account status
	Followers         []User     `json:"followers"`                     // Number of followers
	Followings        []User     `json:"followings"`                    // Number of following accounts
	Posts             []Post     `json:"post_count"`                    // Number of posts
	LastLogin         *time.Time `json:"last_login,omitempty"`          // Optional last login time
	CreatedAt         time.Time  `json:"created_at"`                    // Timestamp of account creation
	UpdatedAt         time.Time  `json:"updated_at"`                    // Timestamp of the last update
}

type UpdateUserInput struct {
	FullName          *string    `json:"fullName,omitempty"`
	UserName          *string    `json:"userName,omitempty"`
	Bio               *string    `json:"bio,omitempty"`
	Email             *string    `json:"email,omitempty"`
	DateOfBirth       *string    `json:"dateOfBirth,omitempty"`
	ProfilePictureURL *string    `json:"profilePictureUrl,omitempty"`
	CoverPictureURL   *string    `json:"coverPictureUrl,omitempty"`
	Location          *string    `json:"location,omitempty"`
	Website           *string    `json:"website,omitempty"`
	IsPrivate         *bool      `json:"isPrivate,omitempty"`
	UpdatedAt         *time.Time `json:"updated_at,omitempty"`
}

type UserStats struct {
	TotalPosts     int
	TotalFollowers int
	TotalFollowing int
}

type UserResponse struct {
	Success bool
	Message string
	Data    map[string]interface{}
}
