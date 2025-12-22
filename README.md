🛒 NineShop Backend

📌 Overview

NineShop Backend là hệ thống backend cho website thương mại điện tử chuyên bán các sản phẩm công nghệ.
Dự án được xây dựng nhằm mô phỏng một hệ thống e-commerce thực tế, cung cấp các RESTful APIs phục vụ cho frontend, xử lý nghiệp vụ người dùng, sản phẩm và đơn hàng.

Backend được phát triển với Golang, tập trung vào kiến trúc rõ ràng, dễ mở rộng và dễ bảo trì.

🚀 Features

👤 User

   Đăng ký tài khoản

   Xác thực tài khoản bằng OTP qua email

   Đăng nhập với JWT Authentication

   Xem thông tin cá nhân

Xem lịch sử đơn hàng

🛒 Order

Tạo đơn hàng

Theo dõi trạng thái đơn hàng

Xem chi tiết đơn hàng

Xử lý nhiều người dùng đặt hàng đồng thời

📦 Product & Category

Lấy danh sách sản phẩm

Tìm kiếm, lọc sản phẩm theo danh mục, thương hiệu và giá cả

Quản lý danh mục sản phẩm

🛠️ Admin

Quản lý sản phẩm (CRUD)

Quản lý đơn hàng & trạng thái

Quản lý người dùng

## 📐 Sơ Đồ Kiến Trúc (Architecture Diagram)

<img width="722" alt="Architecture Diagram" src="https://github.com/user-attachments/assets/f57efd0b-4dce-43e1-a8b9-213badacfb36" />

---

## ⚙️ Hướng Dẫn Cài Đặt

### 1. Clone dự án
```bash
git clone https://github.com/doanvuduy170921/nineshop-be.git
cd nineshop-be
```
### 2. Cấu hình môi trường
Tạo file `.env` tại thư mục gốc và cấu hình các thông số sau:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=nineshop-be
REDIS_ADDRESS=localhost:6379
GMAIL_APP_PASSWORD=your_gmail_app_password
SERVER_ADDRESS=:8080
```
### 3. Chạy ứng dụng
```bash
go run cmd/main.go
```
### 4. Tải các thư viện phụ thuộc
Hệ thống sẽ tự động tải các package cần thiết dựa trên file go.mod:
```bash
go mod tidy
```
📝 Liên Hệ

Họ tên: Đoàn Vũ Duy

Email: doanvuduyndh@gmail.com

LinkedIn: [duydoanvu17092001](https://www.linkedin.com/in/duydoanvu17092001/)
