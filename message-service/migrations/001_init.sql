CREATE TABLE IF NOT EXISTS channels (
                                        id INTEGER PRIMARY KEY AUTOINCREMENT,
                                        name TEXT NOT NULL UNIQUE,
                                        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS messages (
                                        id INTEGER PRIMARY KEY AUTOINCREMENT,
                                        channel_id INTEGER NOT NULL,
                                        user_id INTEGER NOT NULL,
                                        content TEXT NOT NULL,
                                        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                                        FOREIGN KEY(channel_id) REFERENCES channels(id)
    );
