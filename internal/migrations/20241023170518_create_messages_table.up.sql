-- Messages table to store all messages
CREATE TABLE messages (
                          id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                          from_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                          to_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                          content TEXT NOT NULL,
                          message_type VARCHAR(50) NOT NULL, -- 'text', 'image', 'file', etc.
                          is_read BOOLEAN DEFAULT FALSE,
                          is_deleted BOOLEAN DEFAULT FALSE,
                          created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                          updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    -- Add indexes for frequent queries
                          CONSTRAINT check_different_users CHECK (from_user_id != to_user_id)
);

-- Create indexes for messages table
CREATE INDEX idx_messages_from_user ON messages(from_user_id);
CREATE INDEX idx_messages_to_user ON messages(to_user_id);
CREATE INDEX idx_messages_created_at ON messages(created_at);
CREATE INDEX idx_messages_is_read ON messages(is_read);