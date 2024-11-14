CREATE TABLE posts(
    id BIGSERIAL PRIMARY KEY,
    author VARCHAR(100) NOT NULL,
    title VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    comments_disabled BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE comments(
    id BIGSERIAL PRIMARY KEY,
    author VARCHAR(100) NOT NULL,
    post_id BIGINT NOT NULL,
    parent_id BIGINT NULL,
    content VARCHAR(2000) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (post_id) REFERENCES posts (id)
);