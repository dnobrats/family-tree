CREATE TABLE clan (
    id              BIGSERIAL PRIMARY KEY,
    name            TEXT NOT NULL,
    description     TEXT,

    root_person_id  BIGINT,   -- trưởng chi / người đại diện

    created_at      TIMESTAMPTZ DEFAULT now(),
    updated_at      TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE person (
    id                  BIGSERIAL PRIMARY KEY,

    -- Thông tin cơ bản
    full_name           TEXT NOT NULL,
    gender              SMALLINT NOT NULL, -- 1: nam, 2: nữ, 0: khác/không rõ

    birth_date_solar    DATE,
    birth_year          INT,               -- fallback nếu không biết ngày
    is_alive            BOOLEAN NOT NULL DEFAULT true,

    -- Thông tin mất
    death_date_solar    DATE,
    death_date_lunar    TEXT,               -- để text, tránh xử lý lịch âm phức tạp

    -- Quan hệ huyết thống
    father_id           BIGINT,
    mother_id           BIGINT,

    -- Phân chi
    clan_id             BIGINT,

    -- Thông tin thêm
    address             TEXT,
    grave_location      TEXT,
    note                TEXT,

    created_at          TIMESTAMPTZ DEFAULT now(),
    updated_at          TIMESTAMPTZ DEFAULT now()
);

ALTER TABLE person
ADD CONSTRAINT fk_father
FOREIGN KEY (father_id) REFERENCES person(id);

ALTER TABLE person
ADD CONSTRAINT fk_mother
FOREIGN KEY (mother_id) REFERENCES person(id);

ALTER TABLE person
ADD CONSTRAINT fk_clan
FOREIGN KEY (clan_id) REFERENCES clan(id);

ALTER TABLE person
ADD CONSTRAINT chk_no_self_parent
CHECK (
    id IS NULL
    OR (father_id IS DISTINCT FROM id AND mother_id IS DISTINCT FROM id)
);

CREATE INDEX idx_person_father ON person(father_id);
CREATE INDEX idx_person_mother ON person(mother_id);
CREATE INDEX idx_person_clan   ON person(clan_id);

CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX idx_person_name_trgm
ON person
USING gin (full_name gin_trgm_ops);

CREATE TABLE person_import (
    full_name TEXT,
    gender TEXT,
    father_name TEXT,
    mother_name TEXT,
    birth_year INT,
    death_year INT,
    clan_name TEXT
);

ALTER TABLE clan
ADD COLUMN parent_clan_id BIGINT;

