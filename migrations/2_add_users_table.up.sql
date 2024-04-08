CREATE TABLE IF NOT EXISTS users (
                                     token TEXT NOT NULL,
                                     is_admin BOOLEAN NOT NULL DEFAULT false
);

CREATE INDEX idx_users_token ON users (token);