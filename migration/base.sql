CREATE TABLE users (
    id serial UNIQUE,
    email varchar(254) NOT NULL,
    username varchar(32) NOT NULL,
    password varchar NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);

CREATE TYPE oauth_service AS ENUM ('google')

CREATE TABLE oauth_tokens (
    id serial UNIQUE,
    user_id integer,
    service oauth_service NOT NULL,
    expires_at timestamp without time zone,
    scope varchar,
    token text,
    refresh_token text,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id)
        REFERENCES users(id)
);