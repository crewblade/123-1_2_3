
CREATE TABLE IF NOT EXISTS banners (
                                       id SERIAL PRIMARY KEY,
                                       content JSONB NOT NULL,
                                       feature_id INT NOT NULL,
                                       tag_ids INT[] NOT NULL,
                                       is_active BOOLEAN NOT NULL DEFAULT true,
                                       created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                       updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_feature_tag_id_pair ON banners (feature_id, tag_ids);




