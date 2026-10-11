# Nghiệp vụ & Dữ liệu

## Thuật ngữ gia phả

| Thuật ngữ | Nghĩa trong hệ thống |
|---|---|
| **Thuỷ tổ** | Người đầu tiên của dòng họ — không có cha/mẹ trong dữ liệu |
| **Chi** | Một nhánh của dòng họ (`clan`) |
| **Dâu/rể** | Người kết hôn vào dòng họ; vẫn là một `person` (thường không thuộc chi nào) |

## Thực thể

### Person (`person`)
| Trường | Ghi chú |
|---|---|
| `full_name` | Bắt buộc |
| `gender` | 1 = Nam, 2 = Nữ, 0 = Khác/không rõ |
| `birth_year`, `birth_date_solar` (DATE), `birth_date_lunar` (TEXT) | Chỉ ngày dương lưu dạng DATE; âm lịch lưu chuỗi tự do |
| `is_alive`, `death_year`, `death_date_solar`, `death_date_lunar`, `grave_location` | Nếu còn sống thì các trường ngày mất bị xoá |
| `father_id`, `mother_id` | **Có thể NULL** (tạo được thuỷ tổ) |
| `clan_id` | Có thể NULL |
| `address`, `phone`, `occupation`, `avatar_url`, `note` | Tuỳ chọn |

### Clan (`clan`)
`name` (duy nhất), `description`, `parent_clan_id`, `root_person_id` (người đứng đầu chi).

### Spouse (`spouse`)
Một dòng cho **mỗi cặp**, lưu với `person_a_id < person_b_id` (UNIQUE). Nhờ vậy quan hệ
đối xứng: tra vợ/chồng của A hay B đều ra cùng một dòng. `marriage_year`, `note` tuỳ chọn.

### Admin (`admin_user`)
`username` (duy nhất), `password_hash` (bcrypt).

## Quan hệ trong cây

- **Con của một người** = những người có `father_id` hoặc `mother_id` trỏ tới người đó.
- **Vợ/chồng** hiển thị cạnh nhau trong cây lấy từ bảng `spouse` (đối xứng).
- Cây dựng đệ quy từ người gốc; `visited` chống lặp vô hạn; giới hạn số đời `maxGen`.

## Quy tắc nghiệp vụ (và nơi cài đặt)

| Quy tắc | Nơi cài đặt |
|---|---|
| Họ tên bắt buộc, giới hạn độ dài; giới tính hợp lệ | `service/validate.go` |
| Cha/mẹ/chi **không bắt buộc** | `service/person.go` (validateRelations) |
| Cha phải là Nam, Mẹ phải là Nữ | `service/person.go` (validateParent) |
| Cha/mẹ không thể là chính mình | `service/person.go` |
| **Chống chu trình tổ tiên** (cha/mẹ không được là con cháu) | `service/person.go` (SubtreeContains) |
| Tuổi cha/mẹ tối thiểu `MinParentAge = 12` | `service/person.go` |
| Ngày mất không trước ngày sinh; người còn sống bỏ trống ngày mất | `service/validate.go` |
| **Không xoá người còn con** | `service/person.go` (DeletePerson) |
| Tên chi bắt buộc & duy nhất | `service/clan.go` + UNIQUE trong DB |
| Vợ/chồng: khác nhau, không trùng cặp, năm kết hôn ≥ năm sinh | `service/spouse.go` |
| Gốc cây mặc định = thuỷ tổ (không cha/mẹ, id nhỏ nhất) | `service/genealogy.go` + `repository/person.go` |

## Vì sao thiết kế như vậy

- **Cha/mẹ cho phép NULL**: nếu bắt buộc, không thể tạo thuỷ tổ — lỗi của bản cũ.
- **Quan hệ vợ/chồng một dòng/cặp**: tránh tình trạng một chiều (A thấy B nhưng B không thấy A).
- **Chọn cha/mẹ bằng ID (dropdown)**: tránh nhập tên tự do rồi resolve sai người trùng tên.
- **Âm lịch dạng chuỗi**: tránh sai ngữ nghĩa khi nhồi ngày âm vào kiểu DATE dương lịch.
