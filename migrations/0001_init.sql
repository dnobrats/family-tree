-- ==========================================================================
-- 0001_init.sql — core schema for the family tree app.
-- One row per person; parents are optional so a founding ancestor (thủy tổ)
-- can exist with no father/mother.
-- ==========================================================================

CREATE TABLE IF NOT EXISTS clan (
    id             BIGSERIAL PRIMARY KEY,
    name           TEXT NOT NULL UNIQUE,
    description    TEXT,
    parent_clan_id BIGINT REFERENCES clan(id) ON DELETE SET NULL,
    root_person_id BIGINT,          -- FK added after person exists (circular ref)
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS person (
    id               BIGSERIAL PRIMARY KEY,
    full_name        TEXT NOT NULL,
    gender           SMALLINT NOT NULL DEFAULT 0 CHECK (gender IN (0, 1, 2)), -- 1 nam, 2 nữ, 0 khác

    birth_year       INT,
    birth_date_solar DATE,
    birth_date_lunar TEXT,          -- lunar date kept as free text (e.g. "15/08/Canh Tý")
    is_alive         BOOLEAN NOT NULL DEFAULT true,
    death_year       INT,
    death_date_solar DATE,
    death_date_lunar TEXT,

    father_id        BIGINT REFERENCES person(id) ON DELETE SET NULL,
    mother_id        BIGINT REFERENCES person(id) ON DELETE SET NULL,
    clan_id          BIGINT REFERENCES clan(id) ON DELETE SET NULL,

    address          TEXT,
    phone            TEXT,
    occupation       TEXT,
    avatar_url       TEXT,
    grave_location   TEXT,
    note             TEXT,

    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT chk_no_self_parent
        CHECK (father_id IS DISTINCT FROM id AND mother_id IS DISTINCT FROM id)
);

-- person -> clan FK is defined above; now close the clan -> person loop.
ALTER TABLE clan
    DROP CONSTRAINT IF EXISTS fk_clan_root_person;
ALTER TABLE clan
    ADD CONSTRAINT fk_clan_root_person
    FOREIGN KEY (root_person_id) REFERENCES person(id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_person_father ON person(father_id);
CREATE INDEX IF NOT EXISTS idx_person_mother ON person(mother_id);
CREATE INDEX IF NOT EXISTS idx_person_clan   ON person(clan_id);

CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE INDEX IF NOT EXISTS idx_person_name_trgm
    ON person USING gin (full_name gin_trgm_ops);

-- A couple is stored once per pair (person_a_id < person_b_id) so the
-- relationship is symmetric — querying either partner finds the same row.
CREATE TABLE IF NOT EXISTS spouse (
    id            BIGSERIAL PRIMARY KEY,
    person_a_id   BIGINT NOT NULL REFERENCES person(id) ON DELETE CASCADE,
    person_b_id   BIGINT NOT NULL REFERENCES person(id) ON DELETE CASCADE,
    marriage_year INT,
    note          TEXT,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT chk_spouse_order CHECK (person_a_id < person_b_id),
    CONSTRAINT uq_spouse_pair  UNIQUE (person_a_id, person_b_id)
);

CREATE INDEX IF NOT EXISTS idx_spouse_a ON spouse(person_a_id);
CREATE INDEX IF NOT EXISTS idx_spouse_b ON spouse(person_b_id);

CREATE TABLE IF NOT EXISTS admin_user (
    id            BIGSERIAL PRIMARY KEY,
    username      TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);
