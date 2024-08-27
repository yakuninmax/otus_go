CREATE TABLE IF NOT EXISTS events (
    id serial PRIMARY KEY,
    title TEXT NOT NULL,
    start_date timestamptz NOT NULL,
    end_date timestamptz NOT NULL,
    description TEXT,
    user_id int NOT NULL
);