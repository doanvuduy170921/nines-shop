# Nineshop-BE

> **Một hệ thống Backend E-commerce hiệu năng cao được xây dựng bằng Go, tập trung vào kiến trúc sạch (Clean Architecture) và xử lý đồng thời (Concurrency).**

---

## 🚀 Điểm Nổi Bật Kỹ Thuật (Key Technical Features)

* **Clean Architecture:** Tổ chức code theo các lớp tách biệt (Handler, Service, Repository, DTO) giúp dự án dễ dàng mở rộng, kiểm thử và bảo trì.
* **High-Performance Concurrency:** * Áp dụng mô hình **Producer-Consumer** với **Buffered Channels** để xử lý các tác vụ nặng như gửi Email/OTP chạy ngầm.
    * Sử dụng **Worker Pool** để giới hạn số lượng Goroutine, kiểm soát tài nguyên hệ thống và tránh tình trạng tràn RAM khi có đột biến truy cập.
* **Type-Safe Database:** Sử dụng **SQLC** để biên dịch SQL thuần thành code Go, đảm bảo hiệu suất tối đa và tránh các lỗi sai kiểu dữ liệu ngay từ khi biên dịch.
* **Optimized Caching:** Tích hợp **Redis** để quản lý OTP (với TTL) và cache danh sách sản phẩm, giúp giảm tải từ 60-70% cho Database chính.
* **Graceful Shutdown:** Hệ thống có cơ chế lắng nghe tín hiệu tắt từ OS để đóng các kết nối DB và hoàn thành các tác vụ Worker đang dang dở trước khi dừng hẳn.

---

## 🛠 Tech Stack

* **Language:** Go (Golang)
* **Web Framework:** Gin Gonic
* **Database:** PostgreSQL (Sử dụng `pgxpool` để quản lý kết nối hiệu quả)
* **Caching:** Redis
* **Tools:** * `SQLC` (Type-safe SQL)
    * `Bcrypt` (Hashing password)
    * `Viper/Env` (Quản lý cấu hình)

---

## 📐 Sơ Đồ Kiến Trúc (Architecture Diagram)

<img width="722" height="714" alt="image" src="https://github.com/user-attachments/assets/f57efd0b-4dce-43e1-a8b9-213badacfb36" />


⚙️ Hướng Dẫn Cài Đặt
1. Clone dự án
Bash

git clone [https://github.com/username/nineshop-be.git](https://github.com/username/nineshop-be.git)

cd nineshop-be

2. Cấu hình môi trường

Tạo file .env tại thư mục gốc và cấu hình các thông số sau:

Đoạn mã

DB_HOST=localhost

DB_PORT=5432

DB_USER=root

DB_PASSWORD=your_password

DB_NAME=nineshop-be

REDIS_ADDRESS=localhost:6379

GMAIL_APP_PASSWORD=your_gmail_app_password

SERVER_ADDRESS=:8080

3. Chạy ứng dụng

Bash

go run cmd/main.go

📝 Liên Hệ

Họ tên: Đoàn Vũ Duy

Email: doanvuduyndh@gmail.com

LinkedIn: duydoanvu17092001

Dự án đang trong quá trình phát triển và hoàn thiện các tính năng mới.
