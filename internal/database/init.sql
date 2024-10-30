CREATE TABLE IF NOT EXISTS posts (
                       id INTEGER PRIMARY KEY AUTOINCREMENT,
                       user_id INTEGER NOT NULL,
                       title TEXT,
                       content TEXT,
                       created DATETIME DEFAULT CURRENT_TIMESTAMP,
                       FOREIGN KEY (user_id) REFERENCES users(id)
);


CREATE TABLE IF NOT EXISTS users (
                                     id INTEGER PRIMARY KEY AUTOINCREMENT,
                                     username TEXT NOT NULL UNIQUE,
                                     email TEXT NOT NULL UNIQUE,
                                     password TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS categories (
                                          id INTEGER PRIMARY KEY AUTOINCREMENT,
                                          name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS post_categories (
                                               post_id INTEGER,
                                               category_id INTEGER,
                                               FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
                                               FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
);

INSERT INTO categories (name)
SELECT 'Technology' WHERE NOT EXISTS (SELECT 1 FROM categories WHERE name = 'Technology');
INSERT INTO categories (name)
SELECT 'Entertainment' WHERE NOT EXISTS (SELECT 1 FROM categories WHERE name = 'Entertainment');
INSERT INTO categories (name)
SELECT 'Sports' WHERE NOT EXISTS (SELECT 1 FROM categories WHERE name = 'Sports');
INSERT INTO categories (name)
SELECT 'Education' WHERE NOT EXISTS (SELECT 1 FROM categories WHERE name = 'Education');
INSERT INTO categories (name)
SELECT 'Health' WHERE NOT EXISTS (SELECT 1 FROM categories WHERE name = 'Health');

CREATE TABLE IF NOT EXISTS post_votes (
                                          post_id INTEGER,
                                          user_id INTEGER,
                                          vote_type INTEGER,
                                          PRIMARY KEY (post_id, user_id),
                                          FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
                                          FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);


