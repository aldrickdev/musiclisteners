-- +goose Up
-- +goose StatementBegin
CREATE TABLE production.users(
    "id" SERIAL NOT NULL,
    "name" VARCHAR(255) NOT NULL,
    "avatar" VARCHAR(255) NOT NULL
);
ALTER TABLE production.users ADD PRIMARY KEY("id");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE production.users DROP CONSTRAINT users_pkey;
DROP TABLE production.users;
-- +goose StatementEnd
