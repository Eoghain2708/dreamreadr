CREATE TABLE dream_embeddings (
    dream_id TEXT PRIMARY KEY,
    model TEXT NOT NULL,
    dimensions INTEGER NOT NULL,
    embedding BLOB NOT NULL,

    FOREIGN KEY (dream_id)
        REFERENCES dreams(id)
        ON DELETE CASCADE
);