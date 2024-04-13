CREATE TABLE IF NOT EXISTS users (
                                     token TEXT NOT NULL,
                                     is_admin BOOLEAN NOT NULL DEFAULT false
);

CREATE INDEX IF NOT EXISTS idx_users_token ON users (token);

INSERT INTO users (token, is_admin) VALUES
                                        ('admin_token', true),
                                        ('user_token', false);