-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA production;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP SCHEMA production CASCADE;
-- +goose StatementEnd
