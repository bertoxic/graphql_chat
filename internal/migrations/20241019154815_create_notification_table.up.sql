CREATE TABLE notifications (
                               id          VARCHAR(255) PRIMARY KEY,
                               user_id     UUID NOT NULL,
                               type        VARCHAR(255),
                               title       TEXT,
                               content     TEXT,
                               is_read     BOOLEAN NOT NULL DEFAULT FALSE,
                               created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
                               FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Index on user_id for faster lookups of notifications for a specific user
CREATE INDEX idx_notifications_user_id ON notifications(user_id);

-- Index on is_read for faster lookups of read or unread notifications
CREATE INDEX idx_notifications_is_read ON notifications(is_read);