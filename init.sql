CREATE USER L0_user WITH PASSWORD 'zxc';

CREATE DATABASE L0_database;

GRANT ALL PRIVILEGES ON DATABASE L0_database TO L0_user;

\c L0_database;

CREATE TABLE IF NOT EXISTS orders (
    order_uid VARCHAR(255) PRIMARY KEY CHECK (order_uid <> ''),
    track_number VARCHAR(255),
    entry VARCHAR(255),
    delivery JSONB,
    payment JSONB,
    items JSONB,
    locale VARCHAR(10),
    internal_signature VARCHAR(255),
    customer_id VARCHAR(255),
    delivery_service VARCHAR(255),
    shardkey VARCHAR(10),
    sm_id INTEGER,
    date_created TIMESTAMP,
    oof_shard VARCHAR(10)
);


GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO L0_user;