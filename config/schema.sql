CREATE TABLE slides
(
    no         INTEGER PRIMARY KEY AUTOINCREMENT,
    slide_id   TEXT UNIQUE,
    org        INTEGER,
    agent      INTEGER,
    extension  TEXT,
    synced     INTEGER,
    created_at TIMESTAMP DEFAULT (DATETIME('now', 'localtime')) NOT NULL,
    updated_at TIMESTAMP DEFAULT (DATETIME('now', 'localtime')) NOT NULL
)

CREATE TABLE tiles
(
    no         INTEGER PRIMARY KEY AUTOINCREMENT,
    slide_no   INTEGER,
    level      INTEGER,
    coordinate TEXT,
    size       INTEGER,
    synced     INTEGER,
    created_at TIMESTAMP DEFAULT (DATETIME('now', 'localtime')) NOT NULL,
    updated_at TIMESTAMP DEFAULT (DATETIME('now', 'localtime')) NOT NULL
)