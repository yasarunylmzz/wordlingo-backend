CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    surname TEXT NOT NULL,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    is_verified BOOLEAN DEFAULT FALSE
);

CREATE TABLE desk (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    image_link TEXT,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE card (
    id SERIAL PRIMARY KEY,
    language_1 TEXT NOT NULL,
    language_2 TEXT NOT NULL,
    description TEXT NOT NULL,
    desk_id INTEGER NOT NULL,
    importance_value INTEGER DEFAULT 0,
    FOREIGN KEY (desk_id) REFERENCES desk(id) ON DELETE CASCADE
);

CREATE TABLE todo (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    todo TEXT NOT NULL,
    isTrue BOOLEAN DEFAULT FALSE,
    description TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE progress (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    card_id INTEGER NOT NULL,
    progress_level INTEGER NOT NULL,
    date DATE NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (card_id) REFERENCES card(id) ON DELETE CASCADE
);

CREATE TABLE verification_codes (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    code VARCHAR(6),
    created_at TIMESTAMP DEFAULT NOW()
);


CREATE INDEX idx_users_email ON users (email);
