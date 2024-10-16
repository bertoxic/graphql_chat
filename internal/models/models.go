package models

import "time"

type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]int
	Floatmap  map[string]float32
	Data      map[string]interface{}
	CSRFToken string
	Flash     string
	Warning   string
	Error     string
}

type InputDetails interface {
	Sanitize()
	Validate() error
}

//	type UserDetails struct {
//		ID       string `bson:"id" json:"id"`
//		UserName string `bson:"user_name" json:"user_name"`
//		Email    string `bson:"email" json:"email"`
//		Password string `bson:"password" json:"password"`
//	}
type RegistrationInput struct {
	Email    string `bson:"email" json:"email"`
	Username string `bson:"username" json:"username"`
	Password string `bson:"password" json:"password"`
}

type LoginInput struct {
	Email    string `bson:"email" json:"email"`
	Username string `bson:"username" json:"username"`
	Password string `bson:"password" json:"password"`
}

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
type UserPostStats struct {
	TotalPosts   int `json:"total_posts"`
	TotalLikes   int `json:"total_likes"`
	TotalReposts int `json:"total_reposts"`
}
