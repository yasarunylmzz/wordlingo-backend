CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    surname TEXT NOT NULL,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    is_verified BOOLEAN DEFAULT FALSE
);

CREATE TABLE desk (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    image_link TEXT,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE card (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    language_1 TEXT NOT NULL,
    language_2 TEXT NOT NULL,
    desk_id UUID NOT NULL,
    importance_value INTEGER NOT NULL DEFAULT 0,
    FOREIGN KEY (desk_id) REFERENCES desk(id) ON DELETE CASCADE
); 

CREATE TABLE todo (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    todo TEXT NOT NULL,
    is_true BOOLEAN DEFAULT FALSE,
    description TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE progress (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    card_id UUID NOT NULL,
    progress_level INTEGER NOT NULL,
    date DATE NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (card_id) REFERENCES card(id) ON DELETE CASCADE
);

CREATE TABLE verification_codes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id),
    code VARCHAR(6) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    expires_at TIMESTAMP DEFAULT NOW() + INTERVAL '1 day'
);

CREATE VIEW user_desk_card_count AS
    SELECT 
    u.id AS user_id,
    u.name AS user_name,
    COUNT(DISTINCT d.id) AS desk_count,
    COUNT(DISTINCT c.id) AS card_count
    FROM users u
    LEFT JOIN desk d ON u.id = d.user_id
    LEFT JOIN card c ON d.id = c.desk_id
GROUP BY u.id;


CREATE INDEX idx_users_email ON users (email);
CREATE INDEX idx_verification_codes_user_id ON verification_codes (user_id);
CREATE INDEX idx_verification_codes_code ON verification_codes (code);

ALTER TABLE verification_codes ADD CONSTRAINT unique_code UNIQUE (code);
