
CREATE TABLE users (
    i BIGSERTAL PRIMARY KEY, 
    email VARCHAR(255), 
    first_name VARCHAR(255), 
    last_name VARCHAR(255), 
    password VARCHAR(GO),
    user_active BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO users (
    email,
    first_name,
    last_name,
    password,
    user_active,
    created_at,
    updated_at
) VALUES (
'admin@example.com',
'Admin',
'User',
'$$2a$12$1zGLuYDDNvAThh4RA4avbKuheAMpb1svexSrzQm7up.bnpwQHs0jNe',
true,
'2026-05-11 00:00:00',
'2026-03-11 00:00:00');