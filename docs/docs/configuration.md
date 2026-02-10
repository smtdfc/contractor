---
title: Cấu hình
---

# Hướng dẫn cấu hình

Để Contractor có thể hoạt động chính xác, bạn cần tạo một file cấu hình có tên là `contractor.config.json` đặt tại thư mục gốc của dự án. File này định nghĩa cách thức công cụ tìm kiếm schema và nơi xuất mã nguồn.

### Nội dung file cấu hình mẫu

```json
{
  "dir": ["./contract"],
  "output": "./generated",
  "lang": ["ts", "go"]
}
```

### Giải thích các tham số

| Tham số  | Kiểu dữ liệu    | Mô tả                                                                                                             |
| -------- | --------------- | ----------------------------------------------------------------------------------------------------------------- |
| `dir`    | `Array<string>` | Danh sách các thư mục chứa các file định nghĩa `.ctr`. Contractor sẽ quét toàn bộ các file trong các thư mục này. |
| `output` | `string`        | Đường dẫn đến thư mục sẽ chứa mã nguồn sau khi biên dịch.                                                         |
| `lang`   | `Array<string>` | Danh sách các ngôn ngữ đích mà bạn muốn sinh mã. Hiện tại hỗ trợ: `ts` (TypeScript) và `go` (Go).                 |

### Cấu trúc dự án khuyến nghị

Để quản lý code hiệu quả, bạn nên tổ chức thư mục như sau:

```text
my-project/
├── contract/               # Nơi viết các file định nghĩa (.ctr)
│   ├── auth.ctr
│   └── user.ctr
├── generated/              # Mã nguồn sinh ra sẽ nằm ở đây
│   ├── ts/
│   └── go/
├── contractor.config.json  # File cấu hình
└── package.json

```

### Lưu ý quan trọng

1. **Đường dẫn**: Các đường dẫn trong file cấu hình nên là đường dẫn tương đối tính từ vị trí của file `contractor.config.json`.
2. **Tự động tạo thư mục**: Nếu thư mục `output` chưa tồn tại, Contractor sẽ tự động tạo mới khi bạn chạy lệnh biên dịch.
