-- +goose Up
ALTER TABLE users
DROP CONSTRAINT IF EXISTS users_role_check;

UPDATE users
SET
    role = 'editor'
WHERE
    role = 'user';

ALTER TABLE users ADD CONSTRAINT users_role_check CHECK (
    role IN ('editor', 'admin', 'dock_owner', 'sailor')
);

-- +goose Down
ALTER TABLE users
DROP CONSTRAINT IF EXISTS users_role_check;

UPDATE users
SET
    role = 'user'
WHERE
    role = 'editor';

ALTER TABLE users ADD CONSTRAINT users_role_check CHECK (role IN ('user', 'admin', 'dock_owner', 'sailor'));
