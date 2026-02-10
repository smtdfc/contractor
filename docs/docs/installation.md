---
title: Cài đặt
---

# Hướng dẫn cài đặt

Để bắt đầu sử dụng **Contractor**, bạn cần cài đặt hai thành phần chính: **Bộ biên dịch CLI** (để sinh mã) và **Thư viện Runtime** tương ứng với ngôn ngữ lập trình bạn sử dụng (để xử lý logic validation và kiểu dữ liệu).

## 1. Cài đặt Contractor CLI (Compiler)

Bộ biên dịch dùng để chuyển đổi các file schema `.ctr` của bạn thành mã nguồn. Thành phần này được viết bằng Go và cần được cài đặt trên máy phát triển.

```bash
# Cài đặt phiên bản mới nhất của Contractor CLI
go install github.com/smtdfc/contractor@latest

```

_Lưu ý: Hãy đảm bảo thư mục `$GOPATH/bin` đã được thêm vào biến môi trường `PATH` của hệ thống để có thể gọi lệnh `contractor` ở bất cứ đâu._

## 2. Cài đặt Runtime cho các ngôn ngữ

Mỗi ngôn ngữ khi sử dụng code do Contractor sinh ra đều cần một thư viện runtime đi kèm để hỗ trợ các quy tắc kiểm tra dữ liệu (validation) và quản lý lỗi.

### Cho dự án TypeScript / JavaScript

Nếu bạn sinh code cho Frontend hoặc môi trường Node.js:

```bash
# Sử dụng npm
npm install @contractor/runtime

# Hoặc sử dụng yarn
yarn add @contractor/runtime

```

### Cho dự án Go

Nếu bạn sử dụng Contractor để đồng bộ kiểu dữ liệu cho các dịch vụ Backend bằng Go:

```bash
# Thêm gói contractor vào file go.mod
go get github.com/smtdfc/contractor@latest

```

## 3. Kiểm tra cài đặt

Sau khi cài đặt xong, bạn có thể kiểm tra xem bộ biên dịch đã sẵn sàng chưa bằng lệnh:

```bash
contractor --version

```
