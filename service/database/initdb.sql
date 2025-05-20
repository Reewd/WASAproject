PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS "users" (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    photoId TEXT,
    FOREIGN KEY (photoId) REFERENCES images(uuid)
);

CREATE TABLE IF NOT EXISTS "conversations" (
    id INTEGER PRIMARY KEY,
    isGroup BOOLEAN NOT NULL
);

CREATE TABLE IF NOT EXISTS "messages" (
    id INTEGER PRIMARY KEY,
    senderId INTEGER NOT NULL,
    conversationId INTEGER NOT NULL,
    content TEXT NOT NULL,
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (senderId) REFERENCES users(id),
    FOREIGN KEY (conversationId) REFERENCES conversations(id)
);

CREATE TABLE IF NOT EXISTS "reactions" (
    id INTEGER PRIMARY KEY,
    messageId INTEGER NOT NULL,
    senderId INTEGER NOT NULL,
    content TEXT NOT NULL,
    FOREIGN KEY (messageId) REFERENCES messages(id),
    FOREIGN KEY (senderId) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS "participants" (
    userId INTEGER NOT NULL,
    conversationId INTEGER NOT NULL,
    FOREIGN KEY (userId) REFERENCES users(id),
    FOREIGN KEY (conversationId) REFERENCES conversations(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "images" (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    path TEXT NOT NULL,   
)

