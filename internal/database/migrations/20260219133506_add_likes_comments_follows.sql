-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS likes (
  article_id TEXT NOT NULL,
  liker_user_id TEXT NOT NULL,

  UNIQUE (article_id, liker_user_id),
  FOREIGN KEY (article_id) REFERENCES articles(article_id) ON DELETE CASCADE,
  FOREIGN KEY (liker_user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS comments (
  comment_id TEXT NOT NULL,
  article_id TEXT,
  user_id TEXT NOT NULL,
  reply_id TEXT,
  replys INTEGER,
  content TEXT NOT NULL,

  UNIQUE (comment_id),
  FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS follows (
  follower_id TEXT NOT NULL,
  following_id TEXT NOT NULL,

  UNIQUE (follower_id, following_id),
  FOREIGN KEY (follower_id) REFERENCES users(user_id) ON DELETE CASCADE,
  FOREIGN KEY (following_id) REFERENCES users(user_id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE likes;
DROP TABLE comments;
DROP TABLE follows;
-- +goose StatementEnd
