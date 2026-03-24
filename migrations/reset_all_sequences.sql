-- Reset tất cả sequences về giá trị tiếp theo sau ID lớn nhất
-- Hữu ích sau khi xóa nhiều records

-- Reset person sequence
SELECT setval('person_id_seq', COALESCE((SELECT MAX(id) FROM person), 1), true);

-- Reset clan sequence (nếu có)
SELECT setval('clan_id_seq', COALESCE((SELECT MAX(id) FROM clan), 1), true);

-- Reset spouse sequence
SELECT setval('spouse_id_seq', COALESCE((SELECT MAX(id) FROM spouse), 1), true);

-- Hiển thị kết quả
SELECT 
    'person' as table_name,
    COALESCE((SELECT MAX(id) FROM person), 0) as max_id,
    currval('person_id_seq') as current_sequence,
    currval('person_id_seq') + 1 as next_id
UNION ALL
SELECT 
    'clan' as table_name,
    COALESCE((SELECT MAX(id) FROM clan), 0) as max_id,
    currval('clan_id_seq') as current_sequence,
    currval('clan_id_seq') + 1 as next_id
UNION ALL
SELECT 
    'spouse' as table_name,
    COALESCE((SELECT MAX(id) FROM spouse), 0) as max_id,
    currval('spouse_id_seq') as current_sequence,
    currval('spouse_id_seq') + 1 as next_id;
