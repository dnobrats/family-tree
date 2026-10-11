-- ==========================================================================
-- seed_demo.sql — demo family with ~5 generations (for local testing only).
-- Run with:  make seed     (or: psql "$DATABASE_URL" -f scripts/seed_demo.sql)
--
-- WARNING: this clears existing person/clan/spouse rows.
-- ==========================================================================

BEGIN;

TRUNCATE spouse, person, clan RESTART IDENTITY CASCADE;

INSERT INTO clan (id, name, description, root_person_id) VALUES
    (1, 'Chi Trần', 'Chi trưởng dòng họ Trần (demo)', NULL);

INSERT INTO person (id, full_name, gender, birth_year, death_year, is_alive, father_id, mother_id, clan_id, occupation, address) VALUES
    -- Đời 1 — thủy tổ
    (1, 'Trần Văn Tổ',      1, 1890, 1965, false, NULL, NULL, 1, 'Nông nghiệp', 'Nam Định'),
    (2, 'Lê Thị Tổ',        2, 1895, 1970, false, NULL, NULL, 1, NULL,          'Nam Định'),
    -- Đời 2
    (3, 'Trần Văn Bình',    1, 1920, 2005, false, 1, 2, 1, 'Giáo viên',    'Hà Nội'),
    (4, 'Nguyễn Thị Bình',  2, 1922, 2010, false, NULL, NULL, NULL, NULL,       'Hà Nội'),
    (5, 'Trần Thị Hoa',     2, 1925, 2015, false, 1, 2, 1, 'Buôn bán',     'Hải Phòng'),
    -- Đời 3
    (6, 'Trần Văn Cường',   1, 1950, NULL, true, 3, 4, 1, 'Kỹ sư',        'Hà Nội'),
    (7, 'Phạm Thị Cường',   2, 1952, NULL, true, NULL, NULL, NULL, 'Y tá',  'Hà Nội'),
    (8, 'Trần Thị Dung',    2, 1955, NULL, true, 3, 4, 1, 'Kế toán',      'Đà Nẵng'),
    -- Đời 4
    (9,  'Trần Văn Dũng',   1, 1980, NULL, true, 6, 7, 1, 'Lập trình viên', 'Hà Nội'),
    (10, 'Hoàng Thị Dũng',  2, 1982, NULL, true, NULL, NULL, NULL, 'Giáo viên', 'Hà Nội'),
    (11, 'Trần Thị Em',     2, 1985, NULL, true, 6, 7, 1, 'Bác sĩ',       'TP.HCM'),
    -- Đời 5
    (12, 'Trần Văn Phúc',   1, 2010, NULL, true, 9, 10, 1, 'Học sinh',    'Hà Nội');

INSERT INTO spouse (person_a_id, person_b_id, marriage_year, note) VALUES
    (1, 2,  1915, NULL),
    (3, 4,  1945, NULL),
    (6, 7,  1975, NULL),
    (9, 10, 2005, NULL);

UPDATE clan SET root_person_id = 1 WHERE id = 1;

-- Keep sequences ahead of the explicit IDs above.
SELECT setval('person_id_seq', (SELECT MAX(id) FROM person));
SELECT setval('clan_id_seq',   (SELECT MAX(id) FROM clan));
SELECT setval('spouse_id_seq', (SELECT MAX(id) FROM spouse));

COMMIT;
