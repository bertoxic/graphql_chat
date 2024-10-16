-- Create bookmarks table
CREATE TABLE IF NOT EXISTS bookmarks (
                                         id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                                         user_id UUID NOT NULL REFERENCES users(id),
                                         post_id UUID NOT NULL REFERENCES posts(id),
                                         created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                                         UNIQUE (user_id, post_id)
);
