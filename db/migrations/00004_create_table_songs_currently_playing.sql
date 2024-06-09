-- +goose Up
-- +goose StatementBegin
CREATE TABLE production.songs_currently_playing(
    "id" SERIAL NOT NULL,
    "user_id" BIGINT NOT NULL,
    "song_id" BIGINT NOT NULL
);
ALTER TABLE production.songs_currently_playing ADD PRIMARY KEY("id");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE production.songs_currently_playing DROP CONSTRAINT songs_currently_playing_pkey;
DROP TABLE production.songs_currently_playing;
-- +goose StatementEnd
