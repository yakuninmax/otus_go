CREATE TABLE IF NOT EXISTS events (
    event_id serial PRIMARY KEY,
    title TEXT NOT NULL,
    date timestamptz NOT NULL,
    duration bigint NOT NULL,
    description TEXT,
    user_id int NOT NULL,
);