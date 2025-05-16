CREATE TABLE IF NOT EXISTS subscriptions (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    city VARCHAR(255) NOT NULL,
    frequency VARCHAR(50) NOT NULL,
    token VARCHAR(255) NOT NULL UNIQUE,
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
); 


CREATE USER postgres with encrypted password 'postgres';

GRANT READ, WRITE, UPDATE, DELETE ON DATABASE weather_subscription TO postgres;