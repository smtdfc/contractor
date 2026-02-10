---
title: Biên dịch
---

# Hướng dẫn biên dịch (Generation)

Sau khi đã định nghĩa các schema và hoàn tất file cấu hình `contractor.config.json`, bạn có thể tiến hành sinh mã nguồn (Code Generation) bằng công cụ CLI của Contractor.

### 1. Lệnh biên dịch cơ bản

Sử dụng lệnh sau tại thư mục gốc của dự án (nơi đặt file cấu hình):

```bash
contractor generate

```

Trong đó, dấu `.` đại diện cho thư mục hiện tại. Contractor sẽ tự động tìm kiếm file `contractor.config.json` để thực hiện việc quét các schema và đẻ code vào thư mục `output` đã cấu hình.

### 2. Các tùy chọn nâng cao (CLI Flags)

Trong trường hợp bạn muốn ghi đè cấu hình hoặc chỉ định file cấu hình khác, Contractor hỗ trợ các tham số sau:

- **Chỉ định file cấu hình khác:**

```bash
contractor gen --config ./configs/my-custom-config.json

```

- **Chế độ Watch (Theo dõi thay đổi):**
  Tự động biên dịch lại mỗi khi bạn nhấn `Ctrl + S` lưu file `.ctr`:

```bash
contractor generate . --watch

```

### 3. Quy trình xử lý của Contractor

Khi bạn thực hiện lệnh `gen`, Contractor sẽ thực hiện quy trình 3 bước khép kín:

1. **Parsing**: Quét toàn bộ thư mục trong `dir`, đọc các file `.ctr` và phân tích cú pháp để xây dựng cây thực thể (Abstract Syntax Tree - AST).
2. **Validating**: Kiểm tra tính hợp lệ của schema (ví dụ: các kiểu dữ liệu có tồn tại không, các Annotation có đúng cú pháp không).
3. **Emitting**: Sử dụng các Generator tương ứng (TypeScript, Go,...) để chuyển đổi AST thành các class/struct thực tế với đầy đủ validation logic và getter/setter.

### 4. Tích hợp vào quy trình phát triển

Để tối ưu hóa năng suất, bạn nên thêm lệnh này vào các script quản lý dự án:

**Trong `package.json` (TypeScript):**

```json
{
  "scripts": {
    "contract": "contractor gen .",
    "dev": "npm run contract && vite"
  }
}
```

**Trong file Go (Sử dụng `go generate`):**
Chèn dòng sau vào file `main.go` của bạn:

```go
//go:generate contractor gen .

```

Sau đó chỉ cần chạy: `go generate ./...`
