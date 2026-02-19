-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS appeals (
  user_id TEXT NOT NULL,
  name TEXT NOT NULL,
  picture TEXT NOT NULL,
  username TEXT NOT NULL,
  message TEXT NOT NULL,

  FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS subscribers (
  subscriberId TEXT NOT NULL,
  email TEXT NOT NULL,

  UNIQUE(subscriberId),
  UNIQUE(email)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE appeals;
DROP TABLE subscribers;
-- +goose StatementEnd
