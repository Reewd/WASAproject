PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS "users" (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    photoId TEXT,
    FOREIGN KEY (photoId) REFERENCES images(uuid)
);

CREATE TABLE IF NOT EXISTS "conversations" (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    isGroup BOOLEAN NOT NULL,
    photoId TEXT,
    FOREIGN KEY (photoId) REFERENCES images(uuid)
);

CREATE TABLE IF NOT EXISTS "messages" (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    senderId INTEGER NOT NULL,
    conversationId INTEGER NOT NULL,
    content TEXT,
    photoId TEXT,
    replyTo INTEGER,
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (senderId) REFERENCES users(id),
    FOREIGN KEY (conversationId) REFERENCES conversations(id),
    FOREIGN KEY (photoId) REFERENCES images(uuid),
    FOREIGN KEY (replyTo) REFERENCES messages(id) ON DELETE CASCADE
    CHECK (content IS NOT NULL OR photoId IS NOT NULL)
);

CREATE TABLE IF NOT EXISTS "message_status" (
    messageId INTEGER NOT NULL,
    conversationId INTEGER NOT NULL,
    recipientId INTEGER NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('sent', 'delivered', 'read')),
    deliveredAt DATETIME,
    readAt DATETIME,
    FOREIGN KEY (messageId) REFERENCES messages(id),
    FOREIGN KEY (conversationId) REFERENCES conversations(id),
    FOREIGN KEY (recipientId) REFERENCES users(id)
    
);

CREATE TABLE IF NOT EXISTS "reactions" (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    messageId INTEGER NOT NULL,
    senderId INTEGER NOT NULL,
    content TEXT NOT NULL,
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (messageId) REFERENCES messages(id),
    FOREIGN KEY (senderId) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS "participants" (
    userId INTEGER NOT NULL,
    conversationId INTEGER NOT NULL,
    firstMessageId INTEGER,
    FOREIGN KEY (userId) REFERENCES users(id),
    FOREIGN KEY (conversationId) REFERENCES conversations(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "images" (
    uuid TEXT PRIMARY KEY,
    path TEXT NOT NULL
);

