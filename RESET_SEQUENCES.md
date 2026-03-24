# Reset Database Sequences

Sau khi xóa nhiều records, ID sequence sẽ không tự động reset. Điều này dẫn đến ID mới sẽ nhảy cóc (ví dụ: 1, 2, 3 → xóa 2, 3 → thêm mới → ID = 4 thay vì 2).

## Tại sao cần reset?

PostgreSQL sử dụng SEQUENCE để tự động tăng ID. Khi bạn xóa records, sequence không tự động giảm xuống. Bạn cần reset thủ công.

## Cách 1: Sử dụng SQL Script (Khuyến nghị)

### Reset chỉ person table:
```bash
psql -U your_user -d your_database -f migrations/reset_person_id_sequence.sql
```

### Reset tất cả sequences:
```bash
psql -U your_user -d your_database -f migrations/reset_all_sequences.sql
```

## Cách 2: Sử dụng Makefile

```bash
make reset-sequences
```

**Lưu ý:** Cần set environment variables trước:
```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=your_user
export DB_NAME=your_database
```

## Cách 3: Chạy trực tiếp SQL

Kết nối vào database và chạy:

```sql
-- Reset person sequence
SELECT setval('person_id_seq', COALESCE((SELECT MAX(id) FROM person), 1), true);

-- Reset clan sequence
SELECT setval('clan_id_seq', COALESCE((SELECT MAX(id) FROM clan), 1), true);

-- Reset spouse sequence
SELECT setval('spouse_id_seq', COALESCE((SELECT MAX(id) FROM spouse), 1), true);
```

## Cách 4: Sử dụng Bash Script

```bash
chmod +x scripts/reset_sequences.sh
./scripts/reset_sequences.sh
```

## Kiểm tra kết quả

Sau khi reset, kiểm tra sequence hiện tại:

```sql
SELECT 
    'person' as table_name,
    COALESCE((SELECT MAX(id) FROM person), 0) as max_id,
    currval('person_id_seq') + 1 as next_id
UNION ALL
SELECT 
    'clan' as table_name,
    COALESCE((SELECT MAX(id) FROM clan), 0) as max_id,
    currval('clan_id_seq') + 1 as next_id
UNION ALL
SELECT 
    'spouse' as table_name,
    COALESCE((SELECT MAX(id) FROM spouse), 0) as max_id,
    currval('spouse_id_seq') + 1 as next_id;
```

## Ví dụ

**Trước khi reset:**
- Person IDs: 1, 2, 5, 8 (đã xóa 3, 4, 6, 7)
- Next ID: 15 (sequence tiếp tục từ 15)

**Sau khi reset:**
- Person IDs: 1, 2, 5, 8
- Next ID: 9 (sequence reset về 8 + 1)

## Lưu ý quan trọng

⚠️ **Chỉ reset sequence khi:**
- Bạn đã xóa nhiều records và muốn ID gọn gàng
- Không có foreign key references đến các ID đã xóa
- Bạn hiểu rõ hệ thống của mình

⚠️ **KHÔNG reset sequence khi:**
- Hệ thống đang production với nhiều users
- Có logs hoặc audit trails tham chiếu đến ID cũ
- Có external systems tham chiếu đến ID

## Tự động hóa

Nếu bạn muốn tự động reset sau mỗi lần xóa, có thể thêm trigger:

```sql
CREATE OR REPLACE FUNCTION reset_person_sequence()
RETURNS TRIGGER AS $$
BEGIN
    PERFORM setval('person_id_seq', COALESCE((SELECT MAX(id) FROM person), 1), true);
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER after_person_delete
AFTER DELETE ON person
FOR EACH STATEMENT
EXECUTE FUNCTION reset_person_sequence();
```

**Lưu ý:** Trigger này có thể ảnh hưởng performance nếu xóa nhiều records cùng lúc.
