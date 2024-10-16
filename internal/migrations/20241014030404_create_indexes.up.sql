-- Add indexes for better query performance
CREATE INDEX idx_posts_parent_id ON posts(parent_id);
CREATE INDEX idx_posts_created_at ON posts(created_at);
CREATE INDEX idx_follows_follower_id ON follows(follower_id);
CREATE INDEX idx_follows_followed_id ON follows(followed_id);
CREATE INDEX idx_post_likes_post_id ON post_likes(post_id);
CREATE INDEX idx_post_likes_user_id ON post_likes(user_id);
