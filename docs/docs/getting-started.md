---
title: Bắt đầu với Contractor
---

# Bắt đầu với Contractor

Contractor không chỉ đơn thuần là một ngôn ngữ định nghĩa giao diện (Interface Definition Language - IDL) mà còn là một bộ giải pháp toàn diện bao gồm cả môi trường thực thi (Runtime) hỗ trợ đa ngôn ngữ. Công cụ này được thiết kế với mục tiêu cốt lõi là giải quyết bài toán khó khăn nhất trong các hệ thống phân tán: đảm bảo tính nhất quán của dữ liệu (Data Consistency) và tính toàn vẹn của giao diện lập trình. Bằng cách cung cấp một lớp trừu tượng hóa chuẩn, Contractor giúp các dịch vụ được viết bằng các ngôn ngữ lập trình khác nhau có thể giao tiếp một cách chặt chẽ, giảm thiểu sai sót và tăng tốc độ phát triển hệ thống.

## Trước khi bắt đầu

Để đảm bảo quá trình cài đặt và vận hành **Contractor** diễn ra suôn sẻ, hãy chắc chắn rằng môi trường máy tính của bạn đã đáp ứng các yêu cầu tối thiểu sau:

### Yêu cầu hệ thống

- **Go:** Phiên bản `20.x` trở lên.
- **Node.js:** Phiên bản `20.x` trở lên.

::: tip KIỂM TRA PHIÊN BẢN
Bạn có thể kiểm tra nhanh phiên bản hiện tại bằng cách chạy các lệnh sau trong terminal:

```bash
go version
node -v
```

:::

## Cài đặt

Việc cài đặt Contractor rất đơn giản. Bạn chỉ cần thực hiện theo các bước dưới đây để thiết lập môi trường làm việc với Go.

### 1. Khởi tạo dự án

Trước tiên, hãy tạo một thư mục mới cho dự án của bạn và khởi tạo Go module:

```bash
mkdir my-contractor-project
cd my-contractor-project
go mod init module_name

```

### 2. Cài đặt Contractor CLI

Sử dụng lệnh `go install` để tải và cài đặt công cụ dòng lệnh (CLI) của Contractor trực tiếp từ GitHub:

```bash
go install github.com/smtdfc/contractor@latest

```

> **Lưu ý:** Hãy đảm bảo thư mục `$GOPATH/bin` đã được thêm vào biến môi trường `PATH` của bạn để có thể gọi lệnh `contractor` từ bất kỳ đâu.

---

### 3. Kiểm tra cài đặt

Sau khi cài đặt xong, hãy xác nhận rằng công cụ đã sẵn sàng hoạt động.

::: tip Kiểm tra phiên bản
Chạy lệnh sau trong terminal để kiểm tra phiên bản hiện tại:

```bash
contractor -v

```

Nếu màn hình hiển thị phiên bản (ví dụ: `contractor version 1.0.0`), bạn đã cài đặt thành công!
:::

## Khởi tạo dự án

Sau khi đã cài đặt xong CLI, bạn tiến hành khởi tạo cấu trúc dự án chuẩn của Contractor bằng lệnh sau:

```bash
contractor init

```

Lệnh này sẽ tự động tạo ra bộ khung dự án bao gồm các thành phần cốt lõi:

### Cấu trúc thư mục

- **`/contract`**: Đây là nơi chứa các file định nghĩa interface (`.contractor`). Bạn sẽ viết các thiết kế dữ liệu và dịch vụ của mình tại đây.
- **`/generated`**: Thư mục chứa mã nguồn đã được biên dịch sang các ngôn ngữ đích (như Go, TypeScript, v.v.). Bạn không nên sửa đổi trực tiếp các file trong này.
- **`contractor.config.json`**: File cấu hình chính của dự án. Tại đây, bạn có thể chỉ định ngôn ngữ đầu ra, các plugin đi kèm và các tùy chọn build khác.
  Phần hướng dẫn của bạn đã đi thẳng vào trọng tâm rất tốt. Để giúp người dùng mới dễ dàng nắm bắt logic của **Contractor** (từ lúc định nghĩa đến khi có sản phẩm thực tế), bạn có thể bổ sung thêm một chút giải thích về cú pháp và kết quả đầu ra.

## Viết file định nghĩa đầu tiên

Trong thư mục `/contract`, hãy tạo một file mới (ví dụ: `user.contractor`). Tại đây, bạn sử dụng cú pháp của Contractor để định nghĩa cấu trúc dữ liệu.

```java
// contract/user.contract

@CreateConstruction
model User {
    String userId
    String name
    String avatarUrl
    String displayName
}

```

> **Giải thích:** **`model`**: Từ khóa dùng để định nghĩa một cấu trúc dữ liệu (Object).
>
> - **`@CreateConstruction`**: Một Decorator chỉ định cho Contractor tạo ra các hàm khởi tạo (Constructor) tương ứng trong mã nguồn đích.

## Tạo mã nguồn (Code Generation)

Sau khi đã định nghĩa xong interface, bạn sử dụng lệnh `gen` để biên dịch file này sang ngôn ngữ lập trình mong muốn.

```bash
contractor gen . -lang=ts

```

**Trong đó:**

- **`gen .`**: Quét toàn bộ thư mục hiện tại để tìm các file định nghĩa.
- **`-lang=ts`**: Chỉ định ngôn ngữ đầu ra là **TypeScript**. Bạn cũng có thể đổi thành các ngôn ngữ khác được hỗ trợ như `go`, `java`, v.v.

### Kết quả thu được

Sau khi chạy lệnh, hãy kiểm tra thư mục `/generated`. Bạn sẽ thấy các file mã nguồn (như `User.ts`) đã được tự động tạo ra với đầy đủ các thuộc tính và kiểu dữ liệu mà bạn đã khai báo.

::: code-group

```typescript [TypeScript]
// Kết quả trong /generated/User.ts
export class User {
  userId: string;
  name: string;
  avatarUrl: string;
  displayName: string;
}
```

```go [Go]
// Kết quả trong /generated/user.go
type User struct {
    UserId      string
    Name        string
    AvatarUrl   string
    DisplayName string
}

```

:::
