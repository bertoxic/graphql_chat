CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create users table
CREATE TABLE IF NOT EXISTS users (
                                     id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                                     username VARCHAR(255) UNIQUE NOT NULL,
                                     email VARCHAR(255) UNIQUE NOT NULL,
                                     password TEXT NOT NULL,
                                     full_name VARCHAR(255),
                                     bio TEXT,
                                     date_of_birth DATE,
                                     profile_picture_url TEXT,
                                     cover_picture_url TEXT,
                                     location VARCHAR(255),
                                     website VARCHAR(255),
                                     is_verified BOOLEAN DEFAULT FALSE,
                                     is_private BOOLEAN DEFAULT FALSE,
                                     follower_count INT DEFAULT 0,
                                     following_count INT DEFAULT 0,
                                     post_count INT DEFAULT 0,
                                     last_login TIMESTAMP,
                                     created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                                     updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
--
-- CREATE TABLE IF NOT EXISTS users(
--                                     id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v1(),
--                                     username varchar(255) UNIQUE NOT NULL,
--                                     email varchar(255) UNIQUE NOT NULL,
--                                     password TEXT NOT NULL,
--                                     created_at TIMESTAMP NOT NULL DEFAULT NOW(),
--                                     updated_at TIMESTAMP NOT NULL DEFAULT NOW()
-- );
