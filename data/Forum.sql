-- SQLite

-- Supprimer les tables existantes
DROP TABLE IF EXISTS likes;
DROP TABLE IF EXISTS topic_categories;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS topic_hashtags;
DROP TABLE IF EXISTS hashtags;
DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS topics;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS reports;
DROP TABLE IF EXISTS messenger;
DROP TABLE IF EXISTS dislikes;
-- Création des tables

-- Table Users
DROP TABLE IF EXISTS users;

CREATE TABLE users (
    id TEXT PRIMARY KEY,
    username TEXT UNIQUE,
    name TEXT,
    birthday DATE,
    password TEXT ,
    email TEXT,
    profile_image TEXT,
    role TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);


-- Table Topics
CREATE TABLE topics (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id TEXT,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);
CREATE TABLE reports (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    topic_id INTEGER,
    comment_id INTEGER,
    user_id INTEGER NOT NULL,
    reason TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    status TEXT DEFAULT 'waiting for a response',
    FOREIGN KEY (topic_id) REFERENCES topics (id) ON DELETE CASCADE,
    FOREIGN KEY (comment_id) REFERENCES comments (id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);
-- Table Comments
CREATE TABLE comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    topic_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (topic_id) REFERENCES topics (id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

-- Table Hashtags
CREATE TABLE hashtags (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE
);
CREATE TABLE messenger (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    sender_id INTEGER NOT NULL,
    receiver_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (sender_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (receiver_id) REFERENCES users (id) ON DELETE CASCADE
);
-- Table Topic_Hashtags
CREATE TABLE topic_hashtags (
    topic_id INTEGER NOT NULL,
    hashtag_id INTEGER NOT NULL,
    PRIMARY KEY (topic_id, hashtag_id),
    FOREIGN KEY (topic_id) REFERENCES topics (id) ON DELETE CASCADE,
    FOREIGN KEY (hashtag_id) REFERENCES hashtags (id) ON DELETE CASCADE
);

-- Table Categories
CREATE TABLE categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE
);

-- Table Topic_Categories
CREATE TABLE topic_categories (
    topic_id INTEGER NOT NULL,
    category_id INTEGER NOT NULL,
    PRIMARY KEY (topic_id, category_id),
    FOREIGN KEY (topic_id) REFERENCES topics (id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES categories (id) ON DELETE CASCADE
);

-- Table Likes
CREATE TABLE likes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    topic_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (topic_id) REFERENCES topics (id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE dislikes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    topic_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (topic_id) REFERENCES topics (id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

-- Insertion des données

-- Insertion dans la table Users
INSERT INTO users (id,username,name, password, birthday, email, role) VALUES
(1,'user1','user1', '$2y$10$XU5kkL4flq14JpyK/wenDuSO4Jdqb.EMPgB8td3nW1PhacAtRjmg6', '1999-01-18', 'user1@example.com', 'admin'),
(2,'user2','user2', '$2y$10$3sisfsHwR92g6Jo9udd1TuIgoU3DPgL/9.Z0mkRsKGCxTd10Tlq5a', '1999-01-18', 'user2@example.com', 'moderator'),
(3,'user3','user3', '$2y$10$3sisfsHwR92g6Jo9udd1TuIgoU3DPgL/9.Z0mkRsKGCxTd10Tlq5a', '1999-01-18', 'user3@exemple.com', 'admin');
-- Insertion dans la table Topics
INSERT INTO topics (user_id, title, description) VALUES
(1, 'Topic Title 1', 'Description of topic 1'),
(2, 'Topic Title 2', 'Description of topic 2');

-- Insertion dans la table Comments
INSERT INTO comments (topic_id, user_id, content) VALUES
(1, 1, 'Comment on topic 1 by user 1'),
(1, 2, 'Comment on topic 1 by user 2'),
(2, 1, 'Comment on topic 2 by user 1');

-- Insertion dans la table Hashtags
INSERT INTO hashtags (name) VALUES
('hashtag1'),
('hashtag2');

-- Insertion dans la table Topic_Hashtags
INSERT INTO topic_hashtags (topic_id, hashtag_id) VALUES
(1, 1),
(1, 2),
(2, 1);

-- Insertion dans la table Categories
INSERT INTO categories (name) VALUES
('category1'),
('category2');

-- Insertion dans la table Topic_Categories
INSERT INTO topic_categories (topic_id, category_id) VALUES
(1, 1),
(1, 2),
(2, 1);

-- Insertion dans la table Likes
INSERT INTO likes (topic_id, user_id) VALUES
(1, 1),
(1, 2),
(2, 1);

-- Insertion dans la table Dislikes
INSERT INTO dislikes (topic_id, user_id) VALUES
(1, 1),
(1, 2),
(2, 1);