-- +goose Up
-- +goose StatementBegin
CREATE TABLE production.song_history(
    "id" SERIAL NOT NULL,
    "user_id" BIGINT NOT NULL,
    "song_id" BIGINT NOT NULL,
    "timestamp" DATE NOT NULL
);
ALTER TABLE production.song_history ADD PRIMARY KEY("id");

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE production.song_history DROP CONSTRAINT song_history_pkey;
DROP TABLE production.song_history;
-- +goose StatementEnd
