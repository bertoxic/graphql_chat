-- Drop indexes for post_tags
DROP INDEX IF EXISTS idx_post_tags_post_id;
DROP INDEX IF EXISTS idx_post_tags_tag_id;

-- Drop post_tags table
DROP TABLE IF EXISTS post_tags;
