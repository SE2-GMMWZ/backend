-- +goose Up
ALTER TABLE reviews
ADD COLUMN docking_spot_id UUID;

ALTER TABLE reviews
ADD CONSTRAINT fk_review_docking_spot
FOREIGN KEY (docking_spot_id) REFERENCES docking_spots (dock_id);

-- Assign random docking_spot_id to existing reviews
WITH random_docks AS (
    SELECT
        review_id,
        (SELECT dock_id FROM docking_spots ORDER BY random() LIMIT 1) AS random_dock_id
    FROM reviews
)
UPDATE reviews
SET docking_spot_id = random_docks.random_dock_id
FROM random_docks
WHERE reviews.review_id = random_docks.review_id;

ALTER TABLE reviews
ALTER COLUMN docking_spot_id SET NOT NULL;

-- +goose Down
ALTER TABLE reviews
DROP CONSTRAINT IF EXISTS fk_review_docking_spot;

ALTER TABLE reviews
DROP COLUMN IF EXISTS docking_spot_id;
