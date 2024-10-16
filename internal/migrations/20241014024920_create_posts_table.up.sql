-- Create posts table
CREATE TABLE IF NOT EXISTS posts (
                                     id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                                     user_id UUID NOT NULL REFERENCES users(id),
                                     title VARCHAR(255),
                                     content TEXT,
                                     image_url TEXT,
                                     audio_url TEXT,
                                     video_url TEXT,
                                     parent_id UUID REFERENCES posts(id),
                                     is_edited BOOLEAN DEFAULT FALSE,
                                     is_pinned BOOLEAN DEFAULT FALSE,
                                     is_draft BOOLEAN DEFAULT FALSE,
                                     view_count INT DEFAULT 0,
                                     likes INT NOT NULL DEFAULT 0,
                                     reposts INT NOT NULL DEFAULT 0,
                                     comment_count INT DEFAULT 0,
                                     created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                                     updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- CREATE TABLE IF NOT EXISTS posts (
--                                      id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
--                                      user_id UUID NOT NULL REFERENCES users(id),
--                                      title VARCHAR(255),
--                                      content TEXT,
--                                      image_url TEXT,
--                                      audio_url TEXT,
--                                      parent_id UUID REFERENCES posts(id),
--                                      created_at TIMESTAMP NOT NULL DEFAULT NOW(),
--                                      updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
--                                      likes INT NOT NULL DEFAULT 0,
--                                      reposts INT NOT NULL DEFAULT 0
-- );