CREATE TABLE dream_analyses (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  dream_id TEXT NOT NULL,
  summary TEXT NOT NULL,
  
  FOREIGN KEY (dream_id)
    REFERENCES dreams(id)
    ON DELETE CASCADE
);

CREATE TABLE dream_themes (
  dream_id TEXT NOT NULL,
  theme TEXT NOT NULL,

  PRIMARY KEY (dream_id, theme)

  FOREIGN KEY (dream_id)
    REFERENCES dreams(id)
    ON DELETE CASCADE
);

CREATE TABLE dream_emotions (
  dream_id TEXT NOT NULL,
  emotion TEXT NOT NULL,
  
  PRIMARY KEY (dream_id, emotion)

  FOREIGN KEY (dream_id)
    REFERENCES dreams(id)
    ON DELETE CASCADE
);

CREATE TABLE dream_locations (
    dream_id TEXT NOT NULL,
    location TEXT NOT NULL,

    PRIMARY KEY (dream_id, location),

    FOREIGN KEY (dream_id)
        REFERENCES dreams(id)
        ON DELETE CASCADE
);

CREATE TABLE dream_people (
    dream_id TEXT NOT NULL,
    person TEXT NOT NULL,

    PRIMARY KEY (dream_id, person),

    FOREIGN KEY (dream_id)
        REFERENCES dreams(id)
        ON DELETE CASCADE
);

CREATE TABLE dream_symbols (
    dream_id TEXT NOT NULL,
    symbol TEXT NOT NULL,

    PRIMARY KEY (dream_id, symbol),

    FOREIGN KEY (dream_id)
        REFERENCES dreams(id)
        ON DELETE CASCADE
);

