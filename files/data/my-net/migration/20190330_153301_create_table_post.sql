CREATE TABLE posts (
    id        SERIAL       PRIMARY KEY,
    user_id   BIGINT       NOT NULL    REFERENCES users (id),
    content   VARCHAR(140) NOT NULL,
    posted_at TIMESTAMP    NOT NULL
);

CREATE INDEX idx_posts_posted_at ON posts (posted_at);
