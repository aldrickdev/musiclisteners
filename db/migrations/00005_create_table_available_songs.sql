-- +goose Up
-- +goose StatementBegin
CREATE TABLE production.available_songs(
    "id" SERIAL NOT NULL,
    "track_name" VARCHAR(255) NOT NULL,
    "artists_name" VARCHAR(255) NOT NULL,
    "released_year" BIGINT NOT NULL
);
ALTER TABLE
    production.available_songs ADD PRIMARY KEY("id");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE production.available_songs DROP CONSTRAINT available_songs_pkey;
DROP TABLE production.available_songs;
-- +goose StatementEnd
