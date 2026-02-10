---
title: Field Annotations
---

# Field Annotations

Field Annotations được sử dụng để thiết lập các quy tắc ràng buộc (constraints) và kiểm tra tính hợp lệ của dữ liệu (validation) cho từng trường trong Model. Khi biên dịch, các chỉ thị này sẽ được chuyển hóa thành các bộ logic kiểm tra tương ứng ở ngôn ngữ đích.

> **Lưu ý quan trọng:** Để đảm bảo tính minh bạch, hầu hết các annotation kiểm tra dữ liệu trong Contractor đều bắt buộc phải có tham số `msg` (thông báo lỗi trả về khi vi phạm).

---

## 1. Nhóm kiểm tra độ dài và kích thước

Nhóm này kiểm tra số lượng ký tự của `String` hoặc số lượng phần tử của `Array`. Cấu trúc yêu cầu tham số định mức và thông báo lỗi.

- **`@MinLength(n, msg)`**: Độ dài chuỗi tối thiểu là `n`.
- **`@MaxLength(n, msg)`**: Độ dài chuỗi tối đa là `n`.
- **`@ArrayMinSize(n, msg)`**: Mảng phải có ít nhất `n` phần tử.
- **`@ArrayMaxSize(n, msg)`**: Mảng không được vượt quá `n` phần tử.

```java
model Post {
    @MinLength(10, "Tiêu đề phải có ít nhất 10 ký tự")
    String title

    @ArrayMinSize(1, "Bài viết phải có ít nhất 1 thẻ tag")
    Array<String> tags
}

```

---

## 2. Nhóm kiểm tra giá trị số

Dùng cho các kiểu dữ liệu `Int`, `Float`, hoặc `Number` để giới hạn biên độ giá trị.

- **`@Min(value, msg)`**: Giá trị tối thiểu.
- **`@Max(value, msg)`**: Giá trị tối đa.
- **`@IsInt(msg)`**: Xác định trường này phải là số nguyên (Không cần `msg`).
- **`@IsFloat(msg)`**: Xác định trường này phải là số thực (Không cần `msg`).

```java
model Product {
    @IsInt("")
    @Min(1, "Số lượng trong kho không được nhỏ hơn 1")
    Int stock
}

```

---

## 3. Nhóm kiểm tra định dạng dữ liệu

Đảm bảo chuỗi nhập vào tuân thủ các định dạng tiêu chuẩn. Tham số bắt buộc duy nhất là `msg`.

- **`@IsEmail(msg)`**: Kiểm tra định dạng Email hợp lệ.
- **`@IsUrl(msg)`**: Kiểm tra định dạng liên kết (URL).
- **`@IsPhoneNumber(msg)`**: Kiểm tra định dạng số điện thoại.
- **`@IsDateString(msg)`**: Kiểm tra chuỗi ngày tháng theo chuẩn ISO 8601.

```java
model Contact {
    @IsEmail("Định dạng email không chính xác")
    String email

    @IsUrl("Avatar phải là một đường dẫn hợp lệ")
    String avatarUrl
}

```

## 4. Nhóm điều khiển logic

- **`@Optional`**: Đánh dấu trường này không bắt buộc. Trong mã nguồn đích, trường này sẽ được chuyển đổi sang kiểu dữ liệu có thể chứa giá trị rỗng (`optional` trong TS hoặc `pointer` trong Go).

## 5. Ví dụ tổng quát

```java
model Account {
    @MinLength(6, "Mật khẩu phải từ 6 ký tự")
    String password

    @IsEmail("Email không hợp lệ")
    String email

    @Optional
    @IsDateString("Ngày sinh phải đúng định dạng ISO")
    String birthday

    @Min(0, "Số dư không được là số âm")
    @IsFloat("")
    Float balance
}

```
