-- +goose Up
-- +goose StatementBegin

-- Update dock names with random port name prefixes using VALUES clause
WITH port_names AS (
    SELECT * FROM (VALUES 
        ('Anchor Bay'),
        ('Sunset Harbor'),
        ('Crystal Cove'),
        ('Golden Gate Marina'),
        ('Silver Point'),
        ('Moonlight Bay'),
        ('Thunder Port'),
        ('Serenity Harbor'),
        ('Storm Haven'),
        ('Pearl Marina'),
        ('Diamond Dock'),
        ('Emerald Harbor'),
        ('Ruby Bay'),
        ('Sapphire Port'),
        ('Coral Reef Marina'),
        ('Starfish Cove'),
        ('Seahorse Harbor'),
        ('Dolphin Bay'),
        ('Whale Point'),
        ('Pelican Port'),
        ('Seagull Marina'),
        ('Albatross Harbor'),
        ('Mermaid Cove'),
        ('Neptune Bay'),
        ('Poseidon Port'),
        ('Atlantis Marina'),
        ('Treasure Island'),
        ('Pirate Cove'),
        ('Buccaneer Bay'),
        ('Adventure Harbor'),
        ('Discovery Port'),
        ('Explorer Marina'),
        ('Compass Point'),
        ('Navigator Bay'),
        ('Captain Harbor'),
        ('Admiral Port'),
        ('Seafarer Marina'),
        ('Mariner Bay'),
        ('Voyager Harbor'),
        ('Windward Port')
    ) AS port_names(name)
)
UPDATE docking_spots 
SET name = CONCAT(
    (SELECT name FROM port_names ORDER BY RANDOM() LIMIT 1),
    ' - ',
    docking_spots.name
);

-- Drop the ports table
DROP TABLE ports;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Recreate the ports table
CREATE TABLE ports
(
    port_id       UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name          VARCHAR(255) NOT NULL,
    location      JSONB        NOT NULL,
    description   TEXT,
    owner_id      UUID         NOT NULL,
    docking_spots JSONB,
    services      JSONB,
    is_approved   BOOLEAN      NOT NULL,
    CONSTRAINT fk_port_owner FOREIGN KEY (owner_id) REFERENCES users (user_id)
);

-- Remove the port name prefix from dock names
UPDATE docking_spots
SET name = SUBSTRING(name FROM POSITION(' - ' IN name) + 3)
WHERE name LIKE '% - %';

-- +goose StatementEnd
