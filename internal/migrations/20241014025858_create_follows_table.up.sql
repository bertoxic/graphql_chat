-- Create follows table
CREATE TABLE IF NOT EXISTS follows (
                                       follower_id UUID NOT NULL REFERENCES users(id),
                                       followed_id UUID NOT NULL REFERENCES users(id),
                                       created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                                       PRIMARY KEY (follower_id, followed_id)
);

