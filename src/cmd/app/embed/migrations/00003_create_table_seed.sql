-- +goose Up
-- +goose StatementBegin
CREATE TABLE production.seed (
    "status" SMALLINT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE production.seed;
-- +goose StatementEnd
