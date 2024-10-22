package models

import (
	errorx "github.com/bertoxic/graphqlChat/internal/error"
	"strings"
	"time"
)

var (
	titleMinLength = 2
	bodyMinLength  = 2
)

type CreatePostInput struct {
	Title    *string `json:"title,omitempty"`
	Content  string  `json:"content"`
	ImageURL *string `json:"image_url,omitempty"`
	AudioURL *string `json:"audio_url,omitempty"`
}
type PostResponse struct {
	Success bool   `json:"success,omitempty"`
	Message string `json:"message,omitempty"`
}

func (in *CreatePostInput) Sanitize() {
	in.Content = strings.TrimSpace(in.Content)
	if in.Title != nil {
		trimmedTitle := strings.TrimSpace(*in.Title)
		in.Title = &trimmedTitle
	}
}
func (in *CreatePostInput) Validate() error {
	if len(in.Content) < bodyMinLength {
		return errorx.New(errorx.ErrCodeBadRequest, "length of post body is too short", nil)
	}
	return nil
}

type Post struct {
	ID        string         `json:"id"`
	UserID    string         `json:"user_id"`
	Title     *string        `json:"title,omitempty"`     // Optional title for posts
	Content   string         `json:"content"`             // The main body content
	ImageURL  *string        `json:"image_url,omitempty"` // Optional image (URL)
	VideoURL  *string        `json:"video_url,omitempty"` // Optional video URL
	AudioURL  *string        `json:"audio_url,omitempty"` // Optional audio (URL)
	IsEdited  *bool          `json:"is_edited"`           // Indicates if the post was edited
	IsDraft   *bool          `json:"is_draft"`            // Indicates if the post was edited
	ParentID  *string        `json:"parent_id,omitempty"` // Parent post ID for comments or reposts
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	Likes     int            `json:"likes"`
	Reposts   int            `json:"reposts"`
	Tags      []string       `json:"tags"`
	Children  []*Post        `json:"children"`            // Comments or reposts (children posts)
	Analytics *PostAnalytics `json:"analytics,omitempty"` //  field for analytics

}
type PostAnalytics struct {
	Views         int `json:"views"`
	Reach         int `json:"reach"`
	CommentsCount int `json:"comments"`
	Shares        int `json:"shares"`
}
