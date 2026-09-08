CREATE TABLE dreams (
  id TEXT PRIMARY KEY,
  title TEXT NOT NULL,
  created_at DATETIME NOT NULL,
  raw_text TEXT NOT NULL
)