-- Reverse the migration: Drop indexes first, then drop the table

-- Drop indexes for the messages table
DROP INDEX IF EXISTS idx_messages_from_user;
DROP INDEX IF EXISTS idx_messages_to_user;
DROP INDEX IF EXISTS idx_messages_created_at;
DROP INDEX IF EXISTS idx_messages_is_read;

-- Drop the messages table
DROP TABLE IF EXISTS messages;
