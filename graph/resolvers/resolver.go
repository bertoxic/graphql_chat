package resolvers

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/bertoxic/graphqlChat/graph/model"
	"github.com/bertoxic/graphqlChat/internal/auth"
	"github.com/bertoxic/graphqlChat/internal/posts"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"net/http"
)

//go:generate go run github.com/99designs/gqlgen

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	AuthService auth.AuthService
	UserService auth.UserRepository
	PostService posts.PostService
}

func NewResolver(authService auth.AuthService, userService auth.UserRepository, postService posts.PostService) *Resolver {
	return &Resolver{
		AuthService: authService,
		UserService: userService,
		PostService: postService,
	}
}

func buildBadRequestError(ctx context.Context, err error) error {
	return &gqlerror.Error{
		Err:       err,
		Message:   err.Error(),
		Path:      graphql.GetPath(ctx),
		Locations: nil,
		Extensions: map[string]interface{}{
			"code": http.StatusBadRequest,
		},
	}
}

func convertToModelPost(post *posts.Post) *model.Post {
	if post == nil {
		return nil
	}

	childrenPosts := make([]*model.Post, len(post.Children))
	for i, child := range post.Children {
		childrenPosts[i] = convertToModelPost(child)
	}

	return &model.Post{
		ID:        post.ID,
		UserID:    post.UserID,
		Title:     post.Title,
		Content:   post.Content,
		ImageURL:  post.ImageURL,
		AudioURL:  post.AudioURL,
		ParentID:  post.ParentID,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
		Likes:     post.Likes,
		Reposts:   post.Reposts,
		Children:  childrenPosts,
	}
}
