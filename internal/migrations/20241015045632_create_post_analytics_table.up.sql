-- Create post_analytics table
CREATE TABLE IF NOT EXISTS post_analytics (
                                              id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                                              post_id UUID NOT NULL REFERENCES posts(id),
                                              views INT NOT NULL DEFAULT 0,
                                              unique_views INT NOT NULL DEFAULT 0,
                                              shares INT NOT NULL DEFAULT 0,
                                              engagement_rate FLOAT,
                                              created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                                              updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_bookmarks_user_id ON bookmarks(user_id);
CREATE INDEX idx_bookmarks_post_id ON bookmarks(post_id);
CREATE INDEX idx_post_analytics_post_id ON post_analytics(post_id);