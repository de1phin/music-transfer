CREATE TABLE IF NOT EXISTS Yandex (
    id INT PRIMARY KEY UNIQUE NOT NULL,
    yandex_login VARCHAR,
    yandex_id VARCHAR,
    cookies VARCHAR
);

CREATE TABLE IF NOT EXISTS Spotify (
    id INT PRIMARY KEY UNIQUE NOT NULL,
    access_token VARCHAR,
    refresh_token VARCHAR
);