package posts

import (
	"context"
	"fmt"
	errorx "github.com/bertoxic/graphqlChat/internal/error"
	"log"
	"strings"
	"time"
)

var (
	titleMinLength = 2
	bodyMinLength  = 2
)

// Input struct for creating/updating posts
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

type UserPostStats struct {
	TotalPosts   int `json:"total_posts"`
	TotalLikes   int `json:"total_likes"`
	TotalReposts int `json:"total_reposts"`
}

// PostService defines methods to handle posts, comments, and interactions bla..bla..bla
type PostService interface {
	// Post management
	CreatePost(ctx context.Context, input CreatePostInput, userID string, parentID *string) (*Post, error)
	GetPost(ctx context.Context, postID string) (*Post, error)
	UpdatePost(ctx context.Context, postID string, input CreatePostInput) (*Post, error)
	DeletePost(ctx context.Context, postID string) (PostResponse, error)
	GetAllUserPosts(ctx context.Context, userID string) ([]*Post, error)

	// Repost / Comment functionality
	Repost(ctx context.Context, postID string, userID string) (*Post, error)
	AddComment(ctx context.Context, postID string, input CreatePostInput, userID string) (*Post, error)
	GetPostComments(ctx context.Context, postID string) ([]*Post, error)

	// Like/Unlike a post
	LikePost(ctx context.Context, postID string, userID string) (PostResponse, error)
	UnlikePost(ctx context.Context, postID string, userID string) (PostResponse, error)
	GetUsersWhoLikedPost(ctx context.Context, postID string) ([]string, error)
	// Feed and tagging
	GetUserFeed(ctx context.Context, userID string) ([]*Post, error)
	TagUserInPost(ctx context.Context, postID string, taggedUserID string) (PostResponse, error)

	// search posts
	SearchPosts(ctx context.Context, query string) ([]*Post, error)
	GetTrendingPosts(ctx context.Context, limit int) ([]*Post, error)
	GetPostsByTag(ctx context.Context, tag string) ([]*Post, error)

	BookmarkPost(ctx context.Context, postID string, userID string) (PostResponse, error)
	RemoveBookmark(ctx context.Context, postID string, userID string) (PostResponse, error)
	GetUserBookmarkedPosts(ctx context.Context, userID string) ([]*Post, error)

	//SaveDraft(ctx context.Context, userID string, draftInput CreatePostInput) (*Post, error)
	GetDrafts(ctx context.Context, userID string) ([]*Post, error)
	GetPostAnalytics(ctx context.Context, postID string) (*PostAnalytics, error)
	GetUserPostStats(ctx context.Context, userID string) (*UserPostStats, error)
	HandleNullablePostFields(post *Post)
}

// PostServiceImpl implements the PostService interface
type PostServiceImpl struct {
	Repo *PostRepo
}

func NewPostServiceImpl(repo *PostRepo) *PostServiceImpl {

	return &PostServiceImpl{Repo: repo}
}

func (pr *PostServiceImpl) HandleNullablePostFields(post *Post) {
	if post == nil {
		return
	}
	if post.Title == nil {
		defaultTitle := "Untitled"
		post.Title = &defaultTitle
	}
	if post.ImageURL == nil {
		defaultImageURL := ""
		post.ImageURL = &defaultImageURL
	}
	if post.VideoURL == nil {
		defaultVideoURL := ""
		post.VideoURL = &defaultVideoURL
	}
	if post.AudioURL == nil {
		defaultAudioURL := ""
		post.AudioURL = &defaultAudioURL
	}
	if post.ParentID == nil {
		defaultParentID := ""
		post.ParentID = &defaultParentID
	}

	if post.IsEdited == nil {
		defaultIsEdited := false
		post.IsEdited = &defaultIsEdited
	}
	if post.IsDraft == nil {
		defaultIsDraft := true
		post.IsDraft = &defaultIsDraft
	}

}

func (pr *PostServiceImpl) SearchPosts(ctx context.Context, query string) ([]*Post, error) {
	//pr.lowercaseStringsZ()
	log.Printf("xxxxxxxxxxxxxxxxxxxxxxxxxxxx%s", query)
	posts, err := pr.Repo.SearchPosts(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to search posts: %w", err)
	}
	return posts, nil
}

func (pr *PostServiceImpl) GetTrendingPosts(ctx context.Context, limit int) ([]*Post, error) {
	trendingPosts, err := pr.Repo.GetTrendingPosts(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get trending posts: %w", err)
	}
	return trendingPosts, nil
}

func (pr *PostServiceImpl) GetPostsByTag(ctx context.Context, tag string) ([]*Post, error) {

	posts, err := pr.Repo.GetPostsByTag(ctx, tag)
	if err != nil {
		return nil, fmt.Errorf("failed to get posts by tag: %w", err)
	}
	return posts, nil
}

func (pr *PostServiceImpl) BookmarkPost(ctx context.Context, postID string, userID string) (PostResponse, error) {
	postresp, err := pr.Repo.BookmarkPost(ctx, postID, userID)
	if err != nil {
		return postresp, fmt.Errorf("failed to bookmark post: %w", err)
	}
	return postresp, nil
}

func (pr *PostServiceImpl) RemoveBookmark(ctx context.Context, postID string, userID string) (PostResponse, error) {
	postresp, err := pr.Repo.RemoveBookmark(ctx, postID, userID)
	if err != nil {
		return postresp, fmt.Errorf("failed to remove bookmark: %w", err)
	}
	return postresp, nil
}

func (pr *PostServiceImpl) GetUserBookmarkedPosts(ctx context.Context, userID string) ([]*Post, error) {
	bookmarkedPosts, err := pr.Repo.GetUserBookmarkedPosts(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get bookmarked posts: %w", err)
	}
	return bookmarkedPosts, nil
}

//func (pr *PostServiceImpl) SaveDraft(ctx context.Context, userID string, draftInput CreatePostInput) (*Post, error) {
//	draft, err := pr.Repo.SaveDraft(ctx, userID, draftInput)
//	if err != nil {
//		return nil, fmt.Errorf("failed to save draft: %w", err)
//	}
//	return draft, nil
//}

func (pr *PostServiceImpl) GetDrafts(ctx context.Context, userID string) ([]*Post, error) {
	drafts, err := pr.Repo.GetDrafts(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get drafts: %w", err)
	}
	return drafts, nil
}
func (pr *PostServiceImpl) GetPostAnalytics(ctx context.Context, postID string) (*PostAnalytics, error) {
	analytics, err := pr.Repo.GetPostAnalytics(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to get post analytics: %w", err)
	}
	return analytics, nil
}

func (pr *PostServiceImpl) GetUserPostStats(ctx context.Context, userID string) (*UserPostStats, error) {
	stats, err := pr.Repo.GetUserPostStats(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user post stats: %w", err)
	}
	return stats, nil
}

func (pr *PostServiceImpl) CreatePost(ctx context.Context, input CreatePostInput, userID string, parentID *string) (*Post, error) {
	post, err := pr.Repo.CreatePost(ctx, input, userID, parentID)
	if err != nil {
		return nil, errorx.New(errorx.ErrCodeDatabase, "failed to create post", err)
	}
	return post, nil
}
func (pr *PostServiceImpl) GetUsersWhoLikedPost(ctx context.Context, postID string) ([]string, error) {
	userIdList, err := pr.Repo.GetUsersWhoLikedPost(ctx, postID)
	if err != nil {
		return nil, err
	}
	return userIdList, nil
}

func (pr *PostServiceImpl) GetPost(ctx context.Context, postID string) (*Post, error) {
	post, err := pr.Repo.GetPost(ctx, postID)
	if err != nil {
		return nil, errorx.New(errorx.ErrCodeDatabase, "failed to get post", err)
	}
	return post, nil
}

func (pr *PostServiceImpl) UpdatePost(ctx context.Context, postID string, input CreatePostInput) (*Post, error) {
	post, err := pr.Repo.UpdatePost(ctx, postID, input)
	if err != nil {
		return nil, errorx.New(errorx.ErrCodeDatabase, "failed to update post", err)
	}
	return post, nil
}

func (pr *PostServiceImpl) DeletePost(ctx context.Context, postID string) (PostResponse, error) {
	postresp, err := pr.Repo.DeletePost(ctx, postID)
	if err != nil {
		return postresp, errorx.New(errorx.ErrCodeDatabase, "failed to delete post", err)
	}
	return postresp, nil
}

func (pr *PostServiceImpl) GetAllUserPosts(ctx context.Context, userID string) ([]*Post, error) {
	posts, err := pr.Repo.GetAllUserPosts(ctx, userID)
	if err != nil {
		return nil, errorx.New(errorx.ErrCodeDatabase, "failed to get user posts", err)
	}
	return posts, nil
}

func (pr *PostServiceImpl) Repost(ctx context.Context, postID string, userID string) (*Post, error) {
	repost, err := pr.Repo.Repost(ctx, postID, userID)
	if err != nil {
		return nil, errorx.New(errorx.ErrCodeDatabase, "failed to repost", err)
	}
	return repost, nil
}

func (pr *PostServiceImpl) AddComment(ctx context.Context, postID string, input CreatePostInput, userID string) (*Post, error) {
	comment, err := pr.Repo.AddComment(ctx, postID, input, userID)
	if err != nil {
		return nil, errorx.New(errorx.ErrCodeDatabase, "failed to add comment", err)
	}
	return comment, nil
}

func (pr *PostServiceImpl) GetPostComments(ctx context.Context, postID string) ([]*Post, error) {
	comments, err := pr.Repo.GetPostComments(ctx, postID)
	if err != nil {
		return nil, errorx.New(errorx.ErrCodeDatabase, "failed to get post comments", err)
	}
	return comments, nil
}

func (pr *PostServiceImpl) LikePost(ctx context.Context, postID string, userID string) (PostResponse, error) {
	postresp, err := pr.Repo.LikePost(ctx, postID, userID)
	if err != nil {
		return postresp, errorx.New(errorx.ErrCodeDatabase, "failed to like post", err)
	}
	return postresp, nil
}

func (pr *PostServiceImpl) UnlikePost(ctx context.Context, postID string, userID string) (PostResponse, error) {
	postresp, err := pr.Repo.UnlikePost(ctx, postID, userID)
	if err != nil {
		return postresp, errorx.New(errorx.ErrCodeDatabase, "failed to unlike post", err)
	}
	return postresp, nil
}

func (pr *PostServiceImpl) GetUserFeed(ctx context.Context, userID string) ([]*Post, error) {
	feed, err := pr.Repo.GetUserFeed(ctx, userID)
	if err != nil {
		return nil, errorx.New(errorx.ErrCodeDatabase, "failed to get feed", err)
	}
	return feed, nil
}

func (pr *PostServiceImpl) TagUserInPost(ctx context.Context, postID string, taggedUserID string) (PostResponse, error) {
	postResp, err := pr.Repo.TagUserInPost(ctx, postID, taggedUserID)
	if err != nil {
		return postResp, errorx.New(errorx.ErrCodeDatabase, "failed to tag user in post", err)
	}
	return postResp, nil
}
