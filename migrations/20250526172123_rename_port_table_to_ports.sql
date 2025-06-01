-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE "port" RENAME TO "ports";
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE "ports" RENAME TO "port";
-- +goose StatementEnd
