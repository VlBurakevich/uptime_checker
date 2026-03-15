-- +goose Up
CREATE TABLE roles
(
    id   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT UNIQUE NOT NULL
);

CREATE TABLE users
(
    id         UUID PRIMARY KEY     DEFAULT gen_random_uuid(),
    username   TEXT UNIQUE NOT NULL,
    email      TEXT UNIQUE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE user_roles
(
    user_id UUID REFERENCES users (id) ON DELETE CASCADE,
    role_id UUID REFERENCES roles (id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, role_id)
);

CREATE TABLE credentials
(
    user_id       UUID PRIMARY KEY REFERENCES users (id) ON DELETE CASCADE,
    password_hash TEXT NOT NULL
);

CREATE TABLE sites
(
    id              UUID PRIMARY KEY     DEFAULT gen_random_uuid(),
    url             TEXT        NOT NULL,
    name            TEXT,
    interval_sec    INTEGER              DEFAULT 60,
    last_checked_at TIMESTAMPTZ,
    next_checked_at TIMESTAMPTZ,
    is_active       BOOLEAN              DEFAULT true,
    user_id         UUID        REFERENCES users (id) ON DELETE SET NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE TABLE site_checks
(
    id          UUID PRIMARY KEY     DEFAULT gen_random_uuid(),
    site_id     UUID REFERENCES sites (id) ON DELETE CASCADE,
    status_code INTEGER,
    latency_ms  BIGINT,
    is_up       BOOLEAN,
    checked_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_sites_active_next ON sites(is_active, next_checked_at);
CREATE INDEX idx_site_checks_site_id ON site_checks(site_id);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);

-- +goose Down
DROP TABLE IF EXISTS site_checks;
DROP TABLE IF EXISTS sites;
DROP TABLE IF EXISTS credentials;
DROP TABLE IF EXISTS user_roles;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS roles;