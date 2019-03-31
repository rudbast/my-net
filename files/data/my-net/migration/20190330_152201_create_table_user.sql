CREATE TABLE users (
    id SERIAL PRIMARY KEY
);

CREATE TABLE follows (
    id          SERIAL PRIMARY KEY,
    user_id     BIGINT NOT NULL REFERENCES users (id),
    followed_id BIGINT NOT NULL REFERENCES users (id)
);

CREATE UNIQUE INDEX idx_unique_user_follows_user_followed ON follows (user_id, followed_id);
