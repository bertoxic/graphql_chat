-- Create post_tags table
CREATE TABLE IF NOT EXISTS post_tags (
                                         post_id UUID NOT NULL REFERENCES posts(id),
                                         tag_id UUID NOT NULL REFERENCES tags(id),
                                         PRIMARY KEY (post_id, tag_id)
);



-- Create indexes for post_tags
CREATE INDEX idx_post_tags_post_id ON post_tags(post_id);
CREATE INDEX idx_post_tags_tag_id ON post_tags(tag_id);