-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

DROP TABLE IF EXISTS users;
CREATE TABLE users
(
    user_id      UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name         VARCHAR(255) NOT NULL,
    surname      VARCHAR(255) NOT NULL,
    email        VARCHAR(255) NOT NULL UNIQUE,
    phone_number VARCHAR(20),
    password     VARCHAR(255) NOT NULL,
    role         TEXT         NOT NULL CHECK (role IN ('user', 'admin', 'dock_owner', 'sailor'))
);

CREATE TABLE port
(
    port_id       UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name          VARCHAR(255) NOT NULL,
    location      JSONB        NOT NULL, -- latitude, longitude, town
    description   TEXT,
    owner_id      UUID         NOT NULL,
    docking_spots JSONB,
    services      JSONB,
    is_approved   BOOLEAN      NOT NULL,
    CONSTRAINT fk_port_owner FOREIGN KEY (owner_id) REFERENCES users (user_id)
);

CREATE TABLE docking_spots
(
    dock_id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name             VARCHAR(255)   NOT NULL,
    location         JSONB          NOT NULL, -- latitude, longitude, town
    description      TEXT,
    owner_id         UUID           NOT NULL,
    services         TEXT,
    services_pricing NUMERIC(10, 2),
    price_per_night  NUMERIC(10, 2) NOT NULL,
    price_per_person NUMERIC(10, 2),
    availability     TEXT           NOT NULL CHECK (availability IN ('available', 'unavailable')),
    CONSTRAINT fk_docking_owner FOREIGN KEY (owner_id) REFERENCES users (user_id)
);

CREATE TABLE bookings
(
    booking_id     UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    sailor_id      UUID NOT NULL,
    dock_id        UUID NOT NULL,
    start_date     DATE NOT NULL,
    end_date       DATE NOT NULL,
    payment_method TEXT NOT NULL CHECK (payment_method IN ('online', 'in-person')),
    payment_status TEXT NOT NULL CHECK (payment_status IN ('paid', 'unpaid')),
    people         INT  NOT NULL,
    CONSTRAINT fk_booking_sailor FOREIGN KEY (sailor_id) REFERENCES users (user_id),
    CONSTRAINT fk_booking_dock FOREIGN KEY (dock_id) REFERENCES docking_spots (dock_id)
);

CREATE TABLE guides
(
    guide_id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title            VARCHAR(255) NOT NULL,
    content          TEXT         NOT NULL,
    author_id        UUID         NOT NULL,
    publication_date TIMESTAMP    NOT NULL,
    images           JSONB,
    links            JSONB,
    location         JSONB        NOT NULL, -- latitude and longitude
    is_approved      BOOLEAN      NOT NULL,
    CONSTRAINT fk_guide_author FOREIGN KEY (author_id) REFERENCES users (user_id)
);

CREATE TABLE comments
(
    comment_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    guide_id   UUID      NOT NULL,
    user_id    UUID      NOT NULL,
    content    TEXT      NOT NULL,
    timestamp  TIMESTAMP NOT NULL,
    CONSTRAINT fk_comment_guide FOREIGN KEY (guide_id) REFERENCES guides (guide_id),
    CONSTRAINT fk_comment_user FOREIGN KEY (user_id) REFERENCES users (user_id)
);

CREATE TABLE reviews
(
    review_id      UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    reviewer_id    UUID          NOT NULL,
    rating         NUMERIC(2, 1) NOT NULL CHECK (rating >= 1 AND rating <= 5),
    comment        TEXT,
    date_of_review DATE          NOT NULL,
    CONSTRAINT fk_review_reviewer FOREIGN KEY (reviewer_id) REFERENCES users (user_id)
);

CREATE TABLE notifications
(
    notification_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id         UUID      NOT NULL,
    message         TEXT      NOT NULL,
    timestamp       TIMESTAMP NOT NULL,
    CONSTRAINT fk_notification_user FOREIGN KEY (user_id) REFERENCES users (user_id)
);

-- +goose Down
DROP TABLE IF EXISTS notifications;
DROP TABLE IF EXISTS reviews;
DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS guides;
DROP TABLE IF EXISTS bookings;
DROP TABLE IF EXISTS docking_spots;
DROP TABLE IF EXISTS port;
DROP TABLE IF EXISTS users;

DROP EXTENSION IF EXISTS "uuid-ossp";