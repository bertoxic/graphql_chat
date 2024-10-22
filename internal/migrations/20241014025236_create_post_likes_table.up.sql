-- Create post_likes table
CREATE TABLE IF NOT EXISTS post_likes (
                                          post_id UUID NOT NULL REFERENCES posts(id),
                                          user_id UUID NOT NULL REFERENCES users(id),
                                          created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                                          PRIMARY KEY (post_id, user_id)
);