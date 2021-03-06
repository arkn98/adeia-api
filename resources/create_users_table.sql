CREATE TABLE users
(
    id           SERIAL PRIMARY KEY,
    employee_id  CITEXT UNIQUE       NOT NULL,
    name         text                NOT NULL,
    email        varchar(120) UNIQUE NOT NULL,
    password     varchar(128)        NOT NULL,
    designation  varchar(255)        NOT NULL,
    is_activated boolean DEFAULT FALSE
);
