-- Bảng lưu quan hệ vợ/chồng
CREATE TABLE IF NOT EXISTS spouse (
    id SERIAL PRIMARY KEY,
    person_id INT8 NOT NULL REFERENCES person(id) ON DELETE CASCADE,
    spouse_id INT8 NOT NULL REFERENCES person(id) ON DELETE CASCADE,
    marriage_year INT,
    note TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(person_id, spouse_id)
);

-- Index để tìm nhanh
CREATE INDEX idx_spouse_person_id ON spouse(person_id);
CREATE INDEX idx_spouse_spouse_id ON spouse(spouse_id);

-- Đảm bảo không tự kết hôn với chính mình
ALTER TABLE spouse ADD CONSTRAINT check_not_self_spouse 
    CHECK (person_id != spouse_id);

-- Comment
COMMENT ON TABLE spouse IS 'Quan hệ vợ/chồng giữa các thành viên';
COMMENT ON COLUMN spouse.person_id IS 'ID của người trong dòng họ';
COMMENT ON COLUMN spouse.spouse_id IS 'ID của vợ/chồng (có thể là dâu/rể từ bên ngoài)';
COMMENT ON COLUMN spouse.marriage_year IS 'Năm kết hôn (optional)';
