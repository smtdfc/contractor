---
title: Hướng dẫn về Model
---

# Hướng dẫn về Model

Trong hệ sinh thái **Contractor**, Model đóng vai trò là đơn vị cấu trúc cơ bản dùng để định nghĩa hình dạng của dữ liệu. Đây là "hợp đồng" chung giữa các bên tham gia vào hệ thống phân tán, đảm bảo rằng dữ liệu được gửi đi và nhận về luôn nhất quán bất kể ngôn ngữ lập trình sử dụng.

## 1. Khai báo Model

Cú pháp khai báo sử dụng từ khóa `model`, theo sau là tên Model và khối các trường dữ liệu bên trong dấu ngoặc nhọn.

### Cú pháp tổng quát:

```java
model <ModelName> {
    <Type> <fieldName>
    <Type> <fieldName>
}

```

### Ví dụ thực tế:

```java
model User {
    String name
    String address
    Int age
}

```

## 2. Hệ thống kiểu dữ liệu (Data Types)

Contractor cung cấp một tập hợp các kiểu dữ liệu nguyên bản (primitive) và phức hợp (composite) để ánh xạ chính xác sang các ngôn ngữ đích:

| Loại          | Kiểu dữ liệu             | Mô tả                                            |
| ------------- | ------------------------ | ------------------------------------------------ |
| **Văn bản**   | `String`                 | Chuỗi ký tự UTF-8.                               |
| **Số học**    | `Number`, `Int`, `Float` | Các định dạng số nguyên và số thực.              |
| **Logic**     | `Bool`                   | Giá trị đúng/sai (`true`/`false`).               |
| **Danh sách** | `Array<T>`               | Mảng các phần tử có kiểu dữ liệu `T`.            |
| **Đặc biệt**  | `Null`, `Any`, `Unknown` | Dùng cho các trường hợp dữ liệu không định hình. |
| **Đối tượng** | `Object`                 | Đại diện cho một cấu trúc Key-Value linh hoạt.   |

## 3. Quy ước chuyển đổi ngôn ngữ (Mapping)

Khi thực thi lệnh `contractor gen`, các định nghĩa Model sẽ được biên dịch sang các cấu trúc dữ liệu tương ứng của ngôn ngữ đích nhằm tối ưu hóa hiệu suất và tính an toàn kiểu (Type-safety).

### Đối với TypeScript (Target: TS)

Model được chuyển đổi thành các **Class**. Việc sử dụng Class giúp tận dụng các tính năng như khởi tạo mặc định và phương thức đính kèm nếu cần.

```typescript
// Kết quả đầu ra trong generated/User.ts
class User {
  name: string;
  address: string;
  age: number;
}
```

### Đối với Go (Target: Go)

Model được chuyển đổi thành các **Struct**. Contractor tự động đính kèm các `json tags` để hỗ trợ quá trình Serialize/Deserialize dữ liệu một cách chuẩn xác.

```go
// Kết quả đầu ra trong generated/user.go
type User struct {
    Name    string `json:"name"`
    Address string `json:"address"`
    Age     int    `json:"age"`
}

```

## 4. Các quy tắc cần lưu ý

1. **Tên Model:** Phải viết hoa chữ cái đầu (PascalCase) để tương thích tốt nhất với các trình biên dịch ngôn ngữ đích.
2. **Tính nhất quán:** Các trường dữ liệu trong Model mặc định là bắt buộc (Required) trừ khi có các Decorator bổ trợ khác.
3. **Mở rộng:** Bạn có thể lồng các Model vào nhau để tạo ra cấu trúc dữ liệu phức tạp hơn.

```java
model Company {
    String companyName
    Array<User> employees
}

```
