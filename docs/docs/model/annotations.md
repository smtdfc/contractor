---
title: Model Annotations
---

# Model Annotations

Annotations trong **Contractor** là các chỉ thị (decorators) được đặt phía trên khai báo Model để điều khiển hành vi tạo mã nguồn (code generation). Chúng giúp tự động hóa việc tạo ra các hàm tiện ích, giảm thiểu mã lặp (boilerplate code) trong các ngôn ngữ đích.

## 1. @CreateConstructor

Annotation này yêu cầu Contractor tạo ra các hàm khởi tạo đối tượng một cách tự động.

- **TypeScript:** Tạo ra một `constructor` bên trong class với đầy đủ các tham số truyền vào.
- **Go:** Tạo ra một hàm khởi tạo theo định dạng `NewModelName()`, giúp việc khởi tạo struct trở nên an toàn và nhất quán.

```java
@CreateConstructor
model User {
    String id
    String username
}

```

## 2. @Data

Annotation `@Data` là một công cụ mạnh mẽ giúp tự động hóa việc truy cập và đóng gói dữ liệu. Khi sử dụng annotation này, Contractor sẽ sinh ra các phương thức truy xuất cho từng trường (field) trong Model.

- **Hành vi:** Tự động tạo các phương thức **Getter** (lấy giá trị) và **Setter** (gán giá trị).
- **Lợi ích:** Đảm bảo tính đóng gói (encapsulation) theo tiêu chuẩn của các ngôn ngữ hướng đối tượng như Java hoặc TypeScript.

```java
@Data
model Product {
    String sku
    Number price
}

```

## 3. @Mapper

`@Mapper` được sử dụng để tạo ra các logic chuyển đổi dữ liệu giữa các cấu trúc khác nhau. Điều này đặc biệt hữu ích khi bạn cần chuyển đổi từ Data Transfer Object (DTO) sang Model nghiệp vụ hoặc ngược lại.

- **Chức năng:** Tự động sinh ra các hàm `toObject` và `fromObject`
- **Ứng dụng:** Hỗ trợ ánh xạ dữ liệu nhanh chóng giữa các tầng (layers) trong hệ thống, đảm bảo tính thống nhất khi dữ liệu đi qua các ranh giới dịch vụ.

```java
@Mapper
model Profile {
    String bio
    String avatarUrl
}

```

## 4. Cách sử dụng kết hợp

Bạn có thể áp dụng đồng thời nhiều annotation lên cùng một Model để tận dụng tối đa khả năng tạo mã tự động:

```java
@CreateConstructor
@Data
@Mapper
model Account {
    Int accountId
    String email
    Bool isActive
}

```
