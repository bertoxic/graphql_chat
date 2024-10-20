package posts

import (
	"context"
	"fmt"
	"github.com/bertoxic/graphqlChat/internal/database"
	"github.com/bertoxic/graphqlChat/internal/database/postgres"
	errorx "github.com/bertoxic/graphqlChat/internal/error"
	"github.com/jackc/pgx/v4"
	"time"
)

type PostRepo struct {
	DB database.DatabaseRepo
}

func NewPostRepo(db database.DatabaseRepo) *PostRepo {
	return &PostRepo{
		DB: db,
	}
}

func (pr *PostRepo) CreatePost(ctx context.Context, input CreatePostInput, userID string, parentID *string) (*Post, error) {
	db, ok := pr.DB.(*postgres.PostgresDBRepo)
	if !ok {
		return nil, fmt.Errorf("pr.Repo does not implement database.Database")
	}

	pgDB := db.DB

	tx, err := pgDB.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	post, err := pr.createPostTx(ctx, tx, input, userID, parentID)
	if err != nil {
		return nil, fmt.Errorf("error creating post: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return post, nil
}

func (pr *PostRepo) createPostTx(ctx context.Context, tx pgx.Tx, input CreatePostInput, userID string, parentID *string) (*Post, error) {
	var parentIDValue interface{}
	if parentID != nil {
		parentIDValue = *parentID
	} else {
		parentIDValue = nil
	}

	query := `
    INSERT INTO posts (user_id, title, content, image_url, audio_url, parent_id, is_edited, created_at, updated_at)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8,$9)
    RETURNING id, user_id, title, content, image_url, audio_url, parent_id, created_at, updated_at, likes, reposts
`
	var post Post
	err := tx.QueryRow(ctx, query,
		userID, input.Title, input.Content, input.ImageURL, input.AudioURL, parentIDValue,
		false, time.Now(), time.Now(),
	).Scan(
		&post.ID, &post.UserID, &post.Title, &post.Content, &post.ImageURL, &post.AudioURL,
		&post.ParentID, &post.CreatedAt, &post.UpdatedAt, &post.Likes, &post.Reposts,
	)

	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	for _, tag := range post.Tags {
		_, err = pr.addTag(ctx, tx, post.ID, tag)
		if err != nil {
			return nil, fmt.Errorf("error adding tag: %w", err)
		}
	}

	return &post, nil
}

func (pr *PostRepo) GetPost(ctx context.Context, postID string) (*Post, error) {
	db, ok := pr.DB.(*postgres.PostgresDBRepo)
	if !ok {
		return nil, fmt.Errorf("pr.Repo does not implement database.Database")
	}

	pgDB := db.DB

	query := `
        WITH RECURSIVE post_tree AS (
            SELECT id, user_id, title, content, image_url, audio_url, parent_id, created_at, updated_at, likes, reposts, 0 AS depth
            FROM posts
            WHERE id = $1
            
            UNION ALL
            
            SELECT p.id, p.user_id, p.title, p.content, p.image_url, p.audio_url, p.parent_id, p.created_at, p.updated_at, p.likes, p.reposts, pt.depth + 1
            FROM posts p
            JOIN post_tree pt ON p.parent_id = pt.id
        )
        SELECT id, user_id, title, content, image_url, audio_url, parent_id, created_at, updated_at, likes, reposts, depth
        FROM post_tree
        ORDER BY depth, created_at DESC
    `

	rows, err := pgDB.Query(ctx, query, postID)
	if err != nil {
		return nil, fmt.Errorf("error fetching post and descendants: %w", err)
	}
	defer rows.Close()

	var rootPost *Post
	postMap := make(map[string]*Post)

	for rows.Next() {
		var post Post
		var depth int
		err := rows.Scan(
			&post.ID, &post.UserID, &post.Title, &post.Content, &post.ImageURL, &post.AudioURL,
			&post.ParentID, &post.CreatedAt, &post.UpdatedAt, &post.Likes, &post.Reposts, &depth,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning post: %w", err)
		}

		postMap[post.ID] = &post

		if depth == 0 {
			rootPost = &post
		} else {
			parent := postMap[*post.ParentID]
			parent.Children = append(parent.Children, &post)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating posts: %w", err)
	}

	if rootPost == nil {
		return nil, errorx.New(errorx.ErrCodeNotFound, "post not found", nil)
	}

	return rootPost, nil
}

func (pr *PostRepo) UpdatePost(ctx context.Context, postID string, input CreatePostInput) (*Post, error) {
	db, ok := pr.DB.(*postgres.PostgresDBRepo)
	if !ok {
		return nil, fmt.Errorf("pr.Repo does not implement database.Database")
	}

	pgDB := db.DB

	tx, err := pgDB.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `
		UPDATE posts
		SET title = $1, content = $2, image_url = $3, audio_url = $4, updated_at = $5
		WHERE id = $6
		RETURNING id, user_id, title, content, image_url, audio_url, parent_id, created_at, updated_at, likes, reposts
	`

	var post Post
	err = tx.QueryRow(ctx, query,
		input.Title, input.Content, input.ImageURL, input.AudioURL, time.Now(), postID,
	).Scan(
		&post.ID, &post.UserID, &post.Title, &post.Content, &post.ImageURL, &post.AudioURL,
		&post.ParentID, &post.CreatedAt, &post.UpdatedAt, &post.Likes, &post.Reposts,
	)

	if err != nil {
		return nil, fmt.Errorf("error updating post: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &post, nil
}

func (pr *PostRepo) DeletePost(ctx context.Context, postID string) (PostResponse, error) {
	db, ok := pr.DB.(*postgres.PostgresDBRepo)
	if !ok {
		return PostResponse{
			Message: "unable to delete post",
			Success: false,
		}, fmt.Errorf("pr.DB does not implement the expected database interface")
	}

	pgDB := db.DB

	// SQL query to delete a post by its ID
	query := `DELETE FROM posts WHERE id = $1`

	_, err := pgDB.Exec(ctx, query, postID)
	if err != nil {
		return PostResponse{
			Message: "error deleting post",
			Success: false,
		}, fmt.Errorf("error deleting post: %w", err)
	}

	return PostResponse{
		Message: "post successfully deleted",
		Success: true,
	}, nil
}

func (pr *PostRepo) LikePost(ctx context.Context, postID string, userID string) (PostResponse, error) {
	db, ok := pr.DB.(*postgres.PostgresDBRepo)
	if !ok {
		return PostResponse{
			Message: "unable to like post",
			Success: false,
		}, fmt.Errorf("pr.DB does not implement the expected database interface")
	}

	pgDB := db.DB

	tx, err := pgDB.Begin(ctx)
	if err != nil {
		return PostResponse{
			Message: "failed to begin transaction",
			Success: false,
		}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Check if the user has already liked the post
	checkQuery := `SELECT EXISTS(SELECT 1 FROM post_likes WHERE post_id = $1 AND user_id = $2)`
	var exists bool
	err = tx.QueryRow(ctx, checkQuery, postID, userID).Scan(&exists)
	if err != nil {
		return PostResponse{
			Message: "error checking existing like",
			Success: false,
		}, fmt.Errorf("error checking existing like: %w", err)
	}

	if exists {
		return PostResponse{
			Message: "user has already liked this post",
			Success: false,
		}, errorx.New(errorx.ErrCodeConflict, "user has already liked this post", nil)
	}

	// If not, add the like
	insertQuery := `INSERT INTO post_likes (post_id, user_id) VALUES ($1, $2)`
	_, err = tx.Exec(ctx, insertQuery, postID, userID)
	if err != nil {
		return PostResponse{
			Message: "error adding like",
			Success: false,
		}, fmt.Errorf("error adding like: %w", err)
	}

	updateQuery := `UPDATE posts SET likes = likes + 1 WHERE id = $1`
	_, err = tx.Exec(ctx, updateQuery, postID)
	if err != nil {
		return PostResponse{
			Message: "error updating post likes count",
			Success: false,
		}, fmt.Errorf("error updating post likes count: %w", err)
	}

	// Commit the transaction if all is successful...htis
	if err = tx.Commit(ctx); err != nil {
		return PostResponse{
			Message: "failed to commit transaction",
			Success: false,
		}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return PostResponse{
		Message: "post successfully liked",
		Success: true,
	}, nil
}

func (pr *PostRepo) UnlikePost(ctx context.Context, postID string, userID string) (PostResponse, error) {
	db, ok := pr.DB.(*postgres.PostgresDBRepo)
	if !ok {
		return PostResponse{
			Message: "unable to unlike post",
			Success: false,
		}, fmt.Errorf("pr.DB does not implement the expected database interface")
	}

	pgDB := db.DB

	tx, err := pgDB.Begin(ctx)
	if err != nil {
		return PostResponse{
			Message: "failed to begin transaction",
			Success: false,
		}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	checkQuery := `SELECT EXISTS(SELECT 1 FROM post_likes WHERE post_id = $1 AND user_id = $2)`
	var exists bool
	err = tx.QueryRow(ctx, checkQuery, postID, userID).Scan(&exists)
	if err != nil {
		return PostResponse{
			Message: "error checking existing like",
			Success: false,
		}, fmt.Errorf("error checking existing like: %w", err)
	}

	if !exists {
		return PostResponse{
			Message: "user has not liked this post",
			Success: false,
		}, errorx.New(errorx.ErrCodeNotFound, "user has not liked this post", nil)
	}

	deleteQuery := `DELETE FROM post_likes WHERE post_id = $1 AND user_id = $2`
	_, err = tx.Exec(ctx, deleteQuery, postID, userID)
	if err != nil {
		return PostResponse{
			Message: "error removing like",
			Success: false,
		}, fmt.Errorf("error removing like: %w", err)
	}

	updateQuery := `UPDATE posts SET likes = likes - 1 WHERE id = $1`
	_, err = tx.Exec(ctx, updateQuery, postID)
	if err != nil {
		return PostResponse{
			Message: "error updating post likes count",
			Success: false,
		}, fmt.Errorf("error updating post likes count: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return PostResponse{
			Message: "failed to commit transaction",
			Success: false,
		}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return PostResponse{
		Message: "successfully unliked post",
		Success: true,
	}, nil
}

func (pr *PostRepo) addTag(ctx context.Context, tx pgx.Tx, postID string, tagName string) (PostResponse, error) {
	_, err := tx.Exec(ctx, "INSERT INTO tags (name) VALUES ($1) ON CONFLICT (name) DO NOTHING", tagName)
	if err != nil {
		return PostResponse{
			Message: "error inserting tag",
			Success: false,
		}, fmt.Errorf("error inserting tag: %w", err)
	}

	_, err = tx.Exec(ctx, `
        INSERT INTO post_tags (post_id, tag_id)
        SELECT $1, id FROM tags WHERE name = $2
        ON CONFLICT (post_id, tag_id) DO NOTHING
    `, postID, tagName)
	if err != nil {
		return PostResponse{
			Message: "error linking tag to post",
			Success: false,
		}, fmt.Errorf("error linking tag to post: %w", err)
	}

	return PostResponse{
		Message: "successfully added tag to post",
		Success: true,
	}, nil
}

func (pr *PostRepo) GetUsersWhoLikedPost(ctx context.Context, postID string) ([]string, error) {
	db, ok := pr.DB.(*postgres.PostgresDBRepo)
	if !ok {
		return nil, fmt.Errorf("pr.Repo does not implement database.Database")
	}

	pgDB := db.DB

	query := `SELECT user_id FROM post_likes WHERE post_id = $1`
	rows, err := pgDB.Query(ctx, query, postID)
	if err != nil {
		return nil, fmt.Errorf("error querying users who liked post: %w", err)
	}
	defer rows.Close()

	var userIDs []string
	for rows.Next() {
		var userID string
		if err := rows.Scan(&userID); err != nil {
			return nil, fmt.Errorf("error scanning user ID: %w", err)
		}
		userIDs = append(userIDs, userID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return userIDs, nil
}

func (pr *PostRepo) GetAllUserPosts(ctx context.Context, userID string) ([]*Post, error) {
	db, ok := pr.DB.(*postgres.PostgresDBRepo)
	if !ok {
		return nil, fmt.Errorf("pr.Repo does not implement database.Database")
	}
	pgDB := db.DB
	query := `
		SELECT id, user_id, title, content, image_url, audio_url, parent_id, created_at, updated_at, likes, reposts
		FROM posts
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := pgDB.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("error fetching user posts: %w", err)
	}
	defer rows.Close()

	var userPosts []*Post
	for rows.Next() {
		var post Post
		err := rows.Scan(
			&post.ID, &post.UserID, &post.Title, &post.Content, &post.ImageURL, &post.AudioURL,
			&post.ParentID, &post.CreatedAt, &post.UpdatedAt, &post.Likes, &post.Reposts,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning post row: %w", err)
		}
		userPosts = append(userPosts, &post)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return userPosts, nil
}

func (pr *PostRepo) Repost(ctx context.Context, postID string, userID string) (*Post, error) {
	db, ok := pr.DB.(*postgres.PostgresDBRepo)
	if !ok {
		return nil, fmt.Errorf("pr.Repo does not implement database.Database")
	}

	pgDB := db.DB

	tx, err := pgDB.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	originalPost, err := pr.getPostTx(ctx, tx, postID)
	if err != nil {
		return nil, fmt.Errorf("error fetching Original post: %w", err)
	}

	repostQuery := `
		INSERT INTO posts (user_id, title, content, image_url, audio_url, parent_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, user_id, title, content, image_url, audio_url, parent_id, created_at, updated_at, likes, reposts
	`

	var repost Post
	err = tx.QueryRow(ctx, repostQuery,
		userID, originalPost.Title, originalPost.Content, originalPost.ImageURL, originalPost.AudioURL, postID,
		time.Now(), time.Now(),
	).Scan(
		&repost.ID, &repost.UserID, &repost.Title, &repost.Content, &repost.ImageURL, &repost.AudioURL,
		&repost.ParentID, &repost.CreatedAt, &repost.UpdatedAt, &repost.Likes, &repost.Reposts,
	)

	if err != nil {
		return nil, fmt.Errorf("error creating repost: %w", err)
	}

	updateQuery := `UPDATE posts SET reposts = reposts + 1 WHERE id = $1`
	_, err = tx.Exec(ctx, updateQuery, postID)
	if err != nil {
		return nil, fmt.Errorf("error updating Original post repost count: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &repost, nil
}

func (pr *PostRepo) AddComment(ctx context.Context, postID string, input CreatePostInput, userID string) (*Post, error) {
	return pr.CreatePost(ctx, input, userID, &postID)
}

func (pr *PostRepo) GetPostComments(ctx context.Context, postID string) ([]*Post, error) {
	db, ok := pr.DB.(*postgres.PostgresDBRepo)
	if !ok {
		return nil, fmt.Errorf("pr.Repo does not implement database.Database")
	}
	pgDB := db.DB

	query := `
		SELECT id, user_id, title, content, image_url, audio_url, parent_id, created_at, updated_at, likes, reposts
		FROM posts
		WHERE parent_id = $1
		ORDER BY created_at ASC
	`

	rows, err := pgDB.Query(ctx, query, postID)
	if err != nil {
		return nil, fmt.Errorf("error fetching post comments: %w", err)
	}
	defer rows.Close()

	var comments []*Post
	for rows.Next() {
		var comment Post
		err := rows.Scan(
			&comment.ID, &comment.UserID, &comment.Title, &comment.Content, &comment.ImageURL, &comment.AudioURL,
			&comment.ParentID, &comment.CreatedAt, &comment.UpdatedAt, &comment.Likes, &comment.Reposts,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning comment row: %w", err)
		}
		comments = append(comments, &comment)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return comments, nil
}

func (pr *PostRepo) GetUserFeed(ctx context.Context, userID string) ([]*Post, error) {
	db, ok := pr.DB.(*postgres.PostgresDBRepo)
	if !ok {
		return nil, fmt.Errorf("pr.Repo does not implement database.Database")
	}

	pgDB := db.DB
	query := `
		SELECT p.id, p.user_id, p.title, p.content, p.image_url, p.audio_url, p.parent_id, p.created_at, p.updated_at, p.likes, p.reposts
		FROM posts p
		JOIN follows f ON p.user_id = f.followed_id
		WHERE f.follower_id = $1
		ORDER BY p.created_at DESC
		LIMIT 50
	`

	rows, err := pgDB.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("error fetching feed: %w", err)
	}
	defer rows.Close()

	var feed []*Post
	for rows.Next() {
		var post Post
		err := rows.Scan(
			&post.ID, &post.UserID, &post.Title, &post.Content, &post.ImageURL, &post.AudioURL,
			&post.ParentID, &post.CreatedAt, &post.UpdatedAt, &post.Likes, &post.Reposts,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning feed post row: %w", err)
		}
		feed = append(feed, &post)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return feed, nil
}

func (pr *PostRepo) TagUserInPost(ctx context.Context, postID string, taggedUserID string) (PostResponse, error) {
	db, ok := pr.DB.(*postgres.PostgresDBRepo)
	if !ok {
		return PostResponse{
			Success: ok,
			Message: "unable to tag  user in post",
		}, fmt.Errorf("pr.Repo does not implement database.Database")
	}

	pgDB := db.DB

	query := `
        INSERT INTO post_tags (post_id, tag_id)
        VALUES ($1, $2)
        ON CONFLICT (post_id, tag_id) DO NOTHING
    `

	_, err := pgDB.Exec(ctx, query, postID, taggedUserID)
	if err != nil {
		return PostResponse{
			Success: false,
			Message: "unable to tag  user in post",
		}, fmt.Errorf("error tagging user in post: %w", err)
	}

	return PostResponse{
		Success: true,
		Message: "success",
	}, nil
}

func (pr *PostRepo) getPostTx(ctx context.Context, tx pgx.Tx, postID string) (*Post, error) {
	query := `
		SELECT id, user_id, title, content, image_url, audio_url, parent_id, created_at, updated_at, likes, reposts
		FROM posts
		WHERE id = $1
	`
	var post Post
	err := tx.QueryRow(ctx, query, postID).Scan(
		&post.ID, &post.UserID, &post.Title, &post.Content, &post.ImageURL, &post.AudioURL,
		&post.ParentID, &post.CreatedAt, &post.UpdatedAt, &post.Likes, &post.Reposts,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errorx.New(errorx.ErrCodeNotFound, "post not found", err)
		}
		return nil, fmt.Errorf("error fetching post: %w", err)
	}

	return &post, nil
}

//////////////////////----------------------NEW important funcs that might break the code ________________________________________//////////////////////////////

func (pr *PostRepo) SearchPosts(ctx context.Context, query string) ([]*Post, error) {
	db, ok := pr.DB.(*postgres.PostgresDBRepo)
	if !ok {
		return nil, fmt.Errorf("pr.DB does not implement database.Database")
	}

	sqlQuery := `
        SELECT id, user_id, title, content, image_url, audio_url, parent_id, created_at, updated_at, likes, reposts
        FROM posts
        WHERE to_tsvector('english', coalesce(title, '') || ' ' || content) @@ plainto_tsquery('english', $1)
        ORDER BY created_at DESC
    `

	rows, err := db.DB.Query(ctx, sqlQuery, query)
	if err != nil {
		return nil, fmt.Errorf("error executing search query: %w", err)
	}
	defer rows.Close()

	var posts []*Post
	for rows.Next() {
		var post Post
		err := rows.Scan(
			&post.ID, &post.UserID, &post.Title, &post.Content, &post.ImageURL, &post.AudioURL,
			&post.ParentID, &post.CreatedAt, &post.UpdatedAt, &post.Likes, &post.Reposts,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning post: %w", err)
		}
		posts = append(posts, &post)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating posts: %w", err)
	}

	return posts, nil
}
func (pr *PostRepo) GetTrendingPosts(ctx context.Context, limit int) ([]*Post, error) {
	db, ok := pr.DB.(*postgres.PostgresDBRepo)
	if !ok {
		return nil, fmt.Errorf("pr.DB does not implement database.Database")
	}

	query := `
        SELECT id, user_id, title, content, image_url, audio_url, parent_id, created_at, updated_at, likes, reposts
        FROM posts
        ORDER BY (likes + reposts) DESC, created_at DESC
        LIMIT $1
    `

	rows, err := db.DB.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("error fetching trending posts: %w", err)
	}
	defer rows.Close()

	var posts []*Post
	for rows.Next() {
		var post Post
		err := rows.Scan(
			&post.ID, &post.UserID, &post.Title, &post.Content, &post.ImageURL, &post.AudioURL,
			&post.ParentID, &post.CreatedAt, &post.UpdatedAt, &post.Likes, &post.Reposts,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning post: %w", err)
		}
		posts = append(posts, &post)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating posts: %w", err)
	}

	return posts, nil
}

func (pr *PostRepo) GetPostsByTag(ctx context.Context, tagName string) ([]*Post, error) {
	db, ok := pr.DB.(*postgres.PostgresDBRepo)
	if !ok {
		return nil, fmt.Errorf("pr.DB does not implement database.Database")
	}
	query := `
        SELECT p.id, p.user_id, p.title, p.content, p.image_url, p.audio_url, p.parent_id, p.created_at, p.updated_at
        FROM posts p
        JOIN post_tags pt ON p.id = pt.post_id
        JOIN tags t ON pt.tag_id = t.id
        WHERE t.name = $1 AND p.is_draft = FALSE
        ORDER BY p.created_at DESC
    `
	rows, err := db.DB.Query(ctx, query, tagName)
	if err != nil {
		return nil, fmt.Errorf("error querying posts by tag: %w", err)
	}
	defer rows.Close()

	var posts []*Post
	for rows.Next() {
		post := &Post{}
		err := rows.Scan(
			&post.ID, &post.UserID, &post.Title, &post.Content, &post.ImageURL, &post.AudioURL,
			&post.ParentID, &post.CreatedAt, &post.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning post: %w", err)
		}

		// Fetch tags for the post
		tags, err := pr.getPostTags(ctx, post.ID)
		if err != nil {
			return nil, fmt.Errorf("error fetching tags for post: %w", err)
		}
		post.Tags = tags

		posts = append(posts, post)
	}

	return posts, nil
}

func (pr *PostRepo) getPostTags(ctx context.Context, postID string) ([]string, error) {
	db, ok := pr.DB.(*postgres.PostgresDBRepo)
	if !ok {
		return nil, fmt.Errorf("pr.DB does not implement database.Database")
	}
	query := `
        SELECT t.name
        FROM tags t
        JOIN post_tags pt ON t.id = pt.tag_id
        WHERE pt.post_id = $1
    `
	rows, err := db.DB.Query(ctx, query, postID)
	if err != nil {
		return nil, fmt.Errorf("error querying post tags: %w", err)
	}
	defer rows.Close()

	var tags []string
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, fmt.Errorf("error scanning tag: %w", err)
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func (pr *PostRepo) GetDrafts(ctx context.Context, userID string) ([]*Post, error) {
	db, ok := pr.DB.(*postgres.PostgresDBRepo)
	if !ok {
		return nil, fmt.Errorf("pr.DB does not implement database.Database")
	}
	query := `
        SELECT id, title, content, image_url, audio_url, parent_id, created_at, updated_at
        FROM posts
        WHERE user_id = $1 AND is_draft = TRUE
        ORDER BY updated_at DESC
    `
	rows, err := db.DB.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("error querying drafts: %w", err)
	}
	defer rows.Close()

	var drafts []*Post
	for rows.Next() {
		draft := &Post{}
		err := rows.Scan(
			&draft.ID, &draft.Title, &draft.Content, &draft.ImageURL, &draft.AudioURL,
			&draft.ParentID, &draft.CreatedAt, &draft.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning draft: %w", err)
		}
		d := true
		draft.IsDraft = &d
		draft.UserID = userID

		// Fetch tags for the draft
		tags, err := pr.getPostTags(ctx, draft.ID)
		if err != nil {
			return nil, fmt.Errorf("error fetching tags for draft: %w", err)
		}
		draft.Tags = tags

		drafts = append(drafts, draft)
	}

	return drafts, nil
}

func (pr *PostRepo) PublishDraft(ctx context.Context, postID string) (PostResponse, error) {
	db, ok := pr.DB.(*postgres.PostgresDBRepo)
	if !ok {
		return PostResponse{Message: "unable to publish draft", Success: ok}, fmt.Errorf("pr.DB does not implement database.Database")
	}
	query := `UPDATE posts SET is_draft = FALSE WHERE id = $1`
	_, err := db.DB.Exec(ctx, query, postID)
	if err != nil {
		return PostResponse{Message: "unable to publish draft", Success: false}, fmt.Errorf("error publishing draft: %w", err)
	}
	return PostResponse{Message: "successfully published draft", Success: true}, nil
}

func (pr *PostRepo) BookmarkPost(ctx context.Context, postID string, userID string) (PostResponse, error) {
	db, ok := pr.DB.(*postgres.PostgresDBRepo)
	if !ok {
		return PostResponse{
			Message: "unable to bookmark post",
			Success: false,
		}, fmt.Errorf("pr.DB does not implement the expected database interface")
	}

	query := `
        INSERT INTO bookmarks (user_id, post_id, created_at)
        VALUES ($1, $2, $3)
        ON CONFLICT (user_id, post_id) DO NOTHING
    `

	_, err := db.DB.Exec(ctx, query, userID, postID, time.Now())
	if err != nil {
		return PostResponse{
			Message: "unable to bookmark post",
			Success: false,
		}, fmt.Errorf("error bookmarking post: %w", err)
	}

	// Return success message and nil error if everything goes well
	return PostResponse{
		Message: "successfully bookmarked post",
		Success: true,
	}, nil
}

func (pr *PostRepo) RemoveBookmark(ctx context.Context, postID string, userID string) (PostResponse, error) {
	db, ok := pr.DB.(*postgres.PostgresDBRepo)
	if !ok {
		return PostResponse{
			Message: "unable to remove bookmark",
			Success: false,
		}, fmt.Errorf("pr.DB does not implement the expected database interface")
	}

	query := `
        DELETE FROM bookmarks
        WHERE user_id = $1 AND post_id = $2
    `

	_, err := db.DB.Exec(ctx, query, userID, postID)
	if err != nil {
		return PostResponse{
			Message: "unable to remove bookmark",
			Success: false,
		}, fmt.Errorf("error removing bookmark: %w", err)
	}

	return PostResponse{
		Message: "successfully removed bookmark",
		Success: true,
	}, nil
}

func (pr *PostRepo) GetUserBookmarkedPosts(ctx context.Context, userID string) ([]*Post, error) {
	db, ok := pr.DB.(*postgres.PostgresDBRepo)
	if !ok {
		return nil, fmt.Errorf("pr.DB does not implement database.Database")
	}

	query := `
        SELECT p.id, p.user_id, p.title, p.content, p.image_url, p.audio_url, p.parent_id, p.created_at, p.updated_at, p.likes, p.reposts
        FROM posts p
        JOIN bookmarks b ON p.id = b.post_id
        WHERE b.user_id = $1
        ORDER BY b.created_at DESC
    `

	rows, err := db.DB.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("error fetching bookmarked posts: %w", err)
	}
	defer rows.Close()

	var posts []*Post
	for rows.Next() {
		var post Post
		err := rows.Scan(
			&post.ID, &post.UserID, &post.Title, &post.Content, &post.ImageURL, &post.AudioURL,
			&post.ParentID, &post.CreatedAt, &post.UpdatedAt, &post.Likes, &post.Reposts,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning post: %w", err)
		}
		posts = append(posts, &post)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating posts: %w", err)
	}

	return posts, nil
}

//func (pr *PostRepo) SaveDraft(ctx context.Context, userID string, draftInput CreatePostInput) (*Post, error) {
//	db, ok := pr.DB.(*postgres.PostgresDBRepo)
//	if !ok {
//		return nil, fmt.Errorf("pr.DB does not implement database.Database")
//	}
//
//	query := `
//        INSERT INTO drafts (user_id, title, content, image_url, audio_url, created_at, updated_at)
//        VALUES ($1, $2, $3, $4, $5, $6, $7)
//        RETURNING id, user_id, title, content, image_url, audio_url, created_at, updated_at
//    `
//
//	var draft Post
//	err := db.DB.QueryRow(ctx, query,
//		userID, draftInput.Title, draftInput.Content, draftInput.ImageURL, draftInput.AudioURL,
//		time.Now(), time.Now(),
//	).Scan(
//		&draft.ID, &draft.UserID, &draft.Title, &draft.Content, &draft.ImageURL, &draft.AudioURL,
//		&draft.CreatedAt, &draft.UpdatedAt,
//	)
//
//	if err != nil {
//		return nil, fmt.Errorf("error saving draft: %w", err)
//	}
//
//	return &draft, nil
//}

func (pr *PostRepo) GetPostAnalytics(ctx context.Context, postID string) (*PostAnalytics, error) {
	db, ok := pr.DB.(*postgres.PostgresDBRepo)
	if !ok {
		return nil, fmt.Errorf("pr.DB does not implement database.Database")
	}

	query := `
        SELECT views, engagement_rate, unique_views, shares
        FROM post_analytics
        WHERE post_id = $1
    `

	var analytics PostAnalytics
	err := db.DB.QueryRow(ctx, query, postID).Scan(
		&analytics.Views, &analytics.Reach, &analytics.CommentsCount, &analytics.Shares,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errorx.New(errorx.ErrCodeNotFound, "post analytics not found", nil)
		}
		return nil, fmt.Errorf("error fetching post analytics: %w", err)
	}

	return &analytics, nil
}

func (pr *PostRepo) GetUserPostStats(ctx context.Context, userID string) (*UserPostStats, error) {
	db, ok := pr.DB.(*postgres.PostgresDBRepo)
	if !ok {
		return nil, fmt.Errorf("pr.DB does not implement database.Database")
	}

	query := `
        SELECT 
            COUNT(*) as total_posts,
            COALESCE(SUM(likes), 0) as total_likes,
            COALESCE(SUM(reposts), 0) as total_reposts
        FROM posts
        WHERE user_id = $1
    `

	var stats UserPostStats
	err := db.DB.QueryRow(ctx, query, userID).Scan(
		&stats.TotalPosts, &stats.TotalLikes, &stats.TotalReposts,
	)

	if err != nil {
		return nil, fmt.Errorf("error fetching user post stats: %w", err)
	}

	return &stats, nil
}
