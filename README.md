<img width="4427" height="6950" alt="Go Order Processing Architecture-2026-10-31-134835" src="https://github.com/user-attachments/assets/4ff5bd86-00dd-459a-85df-eeec3b3e44be" />
# 🛒 Nine-Shop Backend — High-Performance E-commerce Platform

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-4169E1?style=flat&logo=postgresql)](https://www.postgresql.org/)
[![Redis](https://img.shields.io/badge/Redis-7-DC382D?style=flat&logo=redis)](https://redis.io/)
[![RabbitMQ](https://img.shields.io/badge/RabbitMQ-3-FF6600?style=flat&logo=rabbitmq)](https://www.rabbitmq.com/)
[![Docker](https://img.shields.io/badge/Docker-Compose-2496ED?style=flat&logo=docker)](https://www.docker.com/)

A production-ready e-commerce backend built with **Go** using a **Modular Monolith** architecture. Designed for high concurrency, reliable database transaction management, and asynchronous task execution under heavy load.

🔗 **Live Demo:** [http://46.225.219.255:3000/](http://46.225.219.255:3000/)  
📄 **API Documentation (Swagger):** [http://46.225.219.255:3000/swagger/index.html](http://46.225.219.255:3000/swagger/index.html)

---

## 🛠 Tech Stack & Tools

* **Core Language & Framework:** Go (Golang), Gin Framework
* **Database & ORM/CodeGen:** PostgreSQL 16, `sqlc` (Type-safe SQL query generator), `pgx/v5`
* **Caching & Rate Limiting:** Redis (Cache-Aside pattern)
* **Message Broker & Async Tasks:** RabbitMQ, Worker Pool (Producer-Consumer pattern)
* **Authentication & Authorization:** JWT (RS256 / RBAC)
* **DevOps & Containerization:** Docker, Docker Compose, GitHub Actions
* **API Documentation:** Swagger / OpenAPI 3.0

---

## 💡 Engineering Highlights & Key Solutions

### 1. Concurrency Control & Stock Overselling Prevention
* **Problem:** Under high concurrent checkout requests, multiple users attempting to buy limited stock simultaneously caused stock overselling (race conditions).
* **Solution:** Implemented **Pessimistic Locking** (`SELECT ... FOR UPDATE`) in PostgreSQL transaction blocks during order creation to ensure strict data consistency.

### 2. Database Deadlock Resolution
* **Problem:** Concurrent multi-item orders caused database deadlocks due to inconsistent row-locking orders across different transactions.
* **Solution:** Enforced **Deterministic Resource Sorting** (sorting product IDs numerically prior to acquiring locks) within the Unit-of-Work transaction pattern.

### 3. Asynchronous Task Offloading (Producer-Consumer Pattern)
* **Problem:** Sending verification OTP emails and order notifications synchronously bloated HTTP response times.
* **Solution:** Integrated **RabbitMQ** with a custom Go **Worker Pool** to handle background jobs asynchronously, cutting API response times down to sub-second levels.

### 4. High-Performance Read Caching
* **Problem:** High database query load when fetching popular catalog items and brand filters.
* **Solution:** Applied the **Cache-Aside Pattern** using Redis with configurable TTLs, reducing primary database load significantly during peak traffic.

---

## 📐 Architecture Diagram

<img width="722" alt="Architecture Diagram" src="https://github.com/user-attachments/assets/f57efd0b-4dce-43e1-a8b9-213badacfb36" />

---

## 🚀 Quick Start with Docker Compose

### Prerequisites
* [Docker](https://www.docker.com/) & [Docker Compose](https://docs.docker.com/compose/) installed.

### 1. Clone the repository
```bash
git clone [https://github.com/doanvuduy170921/nineshop-be.git](https://github.com/doanvuduy170921/nineshop-be.git)
cd nineshop-be
