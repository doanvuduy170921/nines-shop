Tên Dự Án: Nineshop-BE
Một hệ thống Backend E-commerce hiệu năng cao được xây dựng bằng Go, tập trung vào kiến trúc sạch (Clean Architecture) và xử lý đồng thời (Concurrency).

🚀 Điểm Nổi Bật Kỹ Thuật (Key Technical Features)

Clean Architecture: Tổ chức code theo các lớp (Handler, Service, Repository, DTO) giúp dễ dàng mở rộng và bảo trì.

High-Performance Concurrency: - Áp dụng mô hình Producer-Consumer với Buffered Channels để xử lý gửi Email/OTP ngầm.

Sử dụng Worker Pool để kiểm soát tài nguyên hệ thống, tránh tràn RAM khi có lượng truy cập lớn.

Type-Safe Database: Sử dụng SQLC để generate code Go từ SQL thuần, đảm bảo hiệu năng và an toàn kiểu dữ liệu.

Optimized Caching: Tích hợp Redis để lưu trữ OTP (TTL) và cache danh sách sản phẩm, giảm tải 60-70% cho Database chính.

Graceful Shutdown: Đảm bảo hệ thống đóng các kết nối DB và hoàn thành các tác vụ ngầm trước khi tắt hoàn toàn.

🛠 Tech Stack
Language: Go (Golang)

Framework: Gin Gonic (HTTP Web Framework)

Database: PostgreSQL (với pgxpool để quản lý kết nối hiệu quả)

Caching: Redis

Tools: SQLC (Generate type-safe code), Bcrypt (Password hashing), Viper/Env (Config management).

📐 Sơ Đồ Kiến Trúc (Architecture Diagram)
<img width="722" height="714" alt="image" src="https://github.com/user-attachments/assets/8f926759-b880-42b4-9fb4-f49202a709d6" />


⚙️ Hướng Dẫn Cài Đặt
Clone dự án:
  
Bash

git clone https://github.com/username/nineshop-be.git
Cấu hình môi trường: Tạo file .env từ .env.example và điền các thông tin:

Đoạn mã

DB_HOST=localhost
DB_PORT=5432
REDIS_ADDRESS=localhost:6379
GMAIL_APP_PASSWORD=your_password
Chạy ứng dụng:

Bash

go run cmd/main.go


📝 Liên Hệ
Tên: Đoàn Vũ Duy

Email: doanvuduyndh@gmail.com

LinkedIn: https://www.linkedin.com/in/duydoanvu17092001/
