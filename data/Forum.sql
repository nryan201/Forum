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

-- Création des tables

-- Table Users
CREATE TABLE users (
                       id INTEGER PRIMARY KEY AUTOINCREMENT,
                       username TEXT NOT NULL UNIQUE,
                       password TEXT NOT NULL,
                       email TEXT UNIQUE,
                       role TEXT,  -- Ajout de la colonne role
                       created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Table Topics
CREATE TABLE topics (
                        id INTEGER PRIMARY KEY AUTOINCREMENT,
                        user_id INTEGER NOT NULL,
                        title TEXT NOT NULL,
                        description TEXT NOT NULL,
                        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
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

-- Insertion des données

-- Insertion dans la table Users
INSERT INTO users (username, password, email, role) VALUES
                                                        ('user1', 'password1', 'user1@example.com', 'admin'),
                                                        ('user2', 'password2', 'user2@example.com', 'user');

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
