-- Reset sequence của person ID về giá trị tiếp theo sau ID lớn nhất hiện có
-- Chạy script này sau khi xóa người để ID tiếp tục đếm từ số gần nhất

-- Lấy ID lớn nhất hiện tại và set sequence
SELECT setval('person_id_seq', COALESCE((SELECT MAX(id) FROM person), 1), true);

-- Kiểm tra kết quả
SELECT 
    'Current max ID: ' || COALESCE(MAX(id), 0) as max_id,
    'Next ID will be: ' || nextval('person_id_seq') as next_id
FROM person;

-- Rollback nextval để không tăng sequence
SELECT setval('person_id_seq', COALESCE((SELECT MAX(id) FROM person), 1), true);

-- Hiển thị thông tin
SELECT 
    'Sequence has been reset!' as status,
    'Max ID: ' || COALESCE((SELECT MAX(id) FROM person), 0) as current_max,
    'Next ID: ' || currval('person_id_seq') + 1 as next_id;
