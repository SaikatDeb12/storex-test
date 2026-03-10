BEGIN;

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TYPE asset_type AS ENUM (
    'laptop',
    'keyboard',
    'mouse',
    'mobile'
);

CREATE TYPE asset_status AS ENUM (
    'available',
    'assigned',
    'in_service',
    'under_repair',
    'damaged'
);

CREATE TYPE user_role AS ENUM (
    'admin',
    'employee',
    'project_manager',
    'asset_manager',
    'employee_manager'
);

CREATE TYPE employment_type AS ENUM (
    'full_time',
    'intern',
    'freelancer'
);

CREATE TYPE asset_owner_type AS ENUM (
    'client',
    'remotestate'
);

CREATE TYPE connection_type AS ENUM (
    'wired',
    'wireless'
);

CREATE TABLE users (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name          TEXT NOT NULL,
    email         TEXT NOT NULL,
    phone_number  TEXT NOT NULL,
    password      TEXT NOT NULL,
    role          user_role DEFAULT 'employee',
    employment    employment_type DEFAULT 'full_time',
    created_at    TIMESTAMPTZ DEFAULT now(),
    archived_at   TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_users_email_active
    ON users (email)
    WHERE archived_at IS NULL;

CREATE TABLE assets (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    brand           TEXT NOT NULL,
    model           TEXT NOT NULL,
    serial_number   TEXT UNIQUE NOT NULL,
    asset_type      asset_type NOT NULL,
    status          asset_status DEFAULT 'available',
    owner_type      asset_owner_type DEFAULT 'remotestate',
    assigned_by_id  UUID REFERENCES users(id),
    assigned_to_id  UUID REFERENCES users(id),
    assigned_at     TIMESTAMPTZ,
    warranty_start  TIMESTAMPTZ NOT NULL,
    warranty_end    TIMESTAMPTZ NOT NULL,
    service_start   TIMESTAMPTZ,
    service_end     TIMESTAMPTZ,
    returned_at     TIMESTAMPTZ,
    created_at      TIMESTAMPTZ DEFAULT now(),
    updated_at      TIMESTAMPTZ,
    archived_at     TIMESTAMPTZ,
    archived_by_id  UUID REFERENCES users(id)
);

CREATE INDEX idx_assets_assigned_to
    ON assets(assigned_to_id);

CREATE TABLE user_sessions (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users(id),
    created_at  TIMESTAMPTZ DEFAULT now(),
    archived_at TIMESTAMPTZ
);

CREATE TABLE laptops (
    id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    asset_id           UUID UNIQUE NOT NULL REFERENCES assets(id),
    processor          TEXT NOT NULL,
    ram                TEXT NOT NULL,
    storage            TEXT NOT NULL,
    operating_system   TEXT NOT NULL,
    charger            TEXT,
    device_password    TEXT NOT NULL
);

CREATE TABLE keyboards (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    asset_id     UUID UNIQUE NOT NULL REFERENCES assets(id),
    layout       TEXT,
    connectivity connection_type NOT NULL
);

CREATE TABLE mice (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    asset_id     UUID UNIQUE NOT NULL REFERENCES assets(id),
    dpi          INT,
    connectivity connection_type NOT NULL
);

CREATE TABLE mobiles (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    asset_id         UUID UNIQUE NOT NULL REFERENCES assets(id),
    operating_system TEXT NOT NULL,
    ram              TEXT NOT NULL,
    storage          TEXT NOT NULL,
    charger          TEXT,
    device_password  TEXT NOT NULL
);

COMMIT;
