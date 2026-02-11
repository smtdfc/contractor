---
title: Cấu hình
---

# Hướng dẫn cấu hình

Để Contractor có thể hoạt động chính xác, bạn cần tạo một file cấu hình có tên là `contractor.config.json` đặt tại thư mục gốc của dự án. File này định nghĩa cách thức công cụ tìm kiếm schema và nơi xuất mã nguồn.

### Nội dung file cấu hình mẫu

```json
{
  "entries": [
    {
      "source": "./contract",
      "output": "./generated/ts",
      "lang": "ts",
      "packageName": "contract"
    },
    {
      "source": "./contract",
      "output": "./generated/go",
      "lang": "go",
      "packageName": "contract"
    }
  ]
}
```

### Giải thích các tham số

| Tham số       | Kiểu dữ liệu | Mô tả                                                                                                    |
| ------------- | ------------ | -------------------------------------------------------------------------------------------------------- |
| `source`      | `string`     | Thư mục chứa các file định nghĩa `.contract`. Contractor sẽ quét toàn bộ các file trong các thư mục này. |
| `output`      | `string`     | Đường dẫn đến thư mục sẽ chứa mã nguồn sau khi biên dịch.                                                |
| `lang`        | `string`     | Ngôn ngữ đích mà bạn muốn sinh mã. Hiện tại hỗ trợ: `ts` (TypeScript) và `go` (Go).                      |
| `packageName` | `string`     | Tên package                                                                                              |

### Cấu trúc dự án khuyến nghị

Để quản lý code hiệu quả, bạn nên tổ chức thư mục như sau:

```text
my-project/
├── contract/               # Nơi viết các file định nghĩa (.ctr)
│   ├── auth.contract
│   └── user.contract
├── generated/              # Mã nguồn sinh ra sẽ nằm ở đây
│   ├── ts/
│   └── go/
├── contractor.config.json  # File cấu hình
└── package.json

```

### Lưu ý:

1. **Đường dẫn**: Các đường dẫn trong file cấu hình nên là đường dẫn tương đối tính từ vị trí của file `contractor.config.json`.
2. **Tự động tạo thư mục**: Nếu thư mục `output` chưa tồn tại, Contractor sẽ tự động tạo mới khi bạn chạy lệnh biên dịch.
