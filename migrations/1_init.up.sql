
CREATE TABLE IF NOT EXISTS banners (
                                       id SERIAL PRIMARY KEY,
                                       title TEXT NOT NULL,
                                       content TEXT NOT NULL,
                                       feature_id INT,
                                       is_active BOOLEAN NOT NULL DEFAULT true,
                                       created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                       updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_feature_id ON banners (feature_id);

CREATE TABLE IF NOT EXISTS tags (
                                    id SERIAL PRIMARY KEY,
                                    name TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_tag_name ON tags (name);

CREATE TABLE IF NOT EXISTS banner_tags (
    banner_id INT REFERENCES banners(id),
    tag_id INT REFERENCES tags(id),
    PRIMARY KEY (banner_id, tag_id)
    );

CREATE INDEX IF NOT EXISTS idx_banner_tags_banner_id ON banner_tags (banner_id);
CREATE INDEX IF NOT EXISTS idx_banner_tags_tag_id ON banner_tags (tag_id);
