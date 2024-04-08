CREATE TABLE IF NOT EXISTS users (
                                     token TEXT NOT NULL,
                                     is_admin BOOLEAN NOT NULL DEFAULT false
);

CREATE INDEX IF NOT EXISTS idx_users_token ON users (token);