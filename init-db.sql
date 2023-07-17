\c erc_721_checks;

CREATE TABLE minters (
    id SERIAL PRIMARY KEY,
    address VARCHAR(255) UNIQUE,
    status INT
);