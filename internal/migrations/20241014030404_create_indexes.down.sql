-- Drop indexes to revert the migration
DROP INDEX IF EXISTS idx_posts_user_id;
DROP INDEX IF EXISTS idx_posts_parent_id;
DROP INDEX IF EXISTS idx_posts_created_at;
DROP INDEX IF EXISTS idx_follows_follower_id;
DROP INDEX IF EXISTS idx_follows_followed_id;
DROP INDEX IF EXISTS idx_post_likes_post_id;
DROP INDEX IF EXISTS idx_post_likes_user_id;
DROP INDEX IF EXISTS idx_bookmarks_user_id;
DROP INDEX IF EXISTS idx_bookmarks_post_id;
DROP INDEX IF EXISTS idx_post_analytics_post_id;
