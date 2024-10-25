CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    name TEXT NOT NULL,
    verification_code VARCHAR(6),
    is_verified BOOLEAN DEFAULT FALSE
);

CREATE TABLE desk (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    user_id INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE card (
    id SERIAL PRIMARY KEY,
    language_1 TEXT NOT NULL,
    language_2 TEXT NOT NULL,
    description TEXT NOT NULL,
    desk_id INTEGER NOT NULL,
    FOREIGN KEY (desk_id) REFERENCES desk(id) ON DELETE CASCADE
);

CREATE INDEX idx_users_email ON users (email);