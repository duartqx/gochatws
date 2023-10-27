CREATE TABLE IF NOT EXISTS Message
(
  id INTEGER NOT NULL PRIMARY KEY,
  created_at DATE NOT NULL,
  text TEXT NOT NULL,
  -- FKS
  chat_id INTEGER NOT NULL,
  user_id INTEGER NOT NULL,
  FOREIGN KEY(chat_id) REFERENCES ChatRoom(id) ON DELETE CASCADE,
  FOREIGN KEY(user_id) REFERENCES User(id) ON DELETE CASCADE
);