# Mini Banking REST API

Proyek ini adalah implementasi layanan **RESTful API** untuk simulasi sistem perbankan sederhana. Dibangun menggunakan bahasa pemrograman **Golang** dengan arsitektur yang modular, aman, dan *scalable*.

Proyek ini dikembangkan sebagai bagian dari tugas mata kuliah **Pengembangan Platform Web (PPW)** di Institut Teknologi Del.

## 📖 Dokumentasi Lengkap (Wiki)

Untuk panduan lengkap mengenai cara penggunaan, struktur database, dan detail endpoint API, silakan kunjungi halaman Wiki kami:

👉 **[Klik disini untuk menuju Wiki Project](../../wiki)**

---

## 🛠️ Teknologi yang Digunakan

Proyek ini dibangun di atas *stack* teknologi berikut:

* **Language:** [Go (Golang)](https://go.dev/) (v1.25+)
* **Framework:** [Gin Web Framework](https://github.com/gin-gonic/gin) - Untuk routing dan HTTP handling yang cepat.
* **Database:** [PostgreSQL](https://www.postgresql.org/) - Relational Database Management System.
* **ORM:** [GORM](https://gorm.io/) - Untuk interaksi database dan migrasi otomatis.
* **Authentication:** [JWT (JSON Web Token)](https://github.com/golang-jwt/jwt) - Untuk keamanan endpoint.
* **Utilities:** `bcrypt` untuk hashing password & `godotenv` untuk manajemen konfigurasi.

## ✨ Fitur Utama

Berdasarkan implementasi kode, API ini mendukung fitur-fitur berikut:

1.  **Authentication**: Registrasi nasabah baru dan Login (mendapatkan Token JWT).
2.  **Account Management**: Pembuatan rekening otomatis saat registrasi dan pengecekan informasi akun.
3.  **Transactions**:
    * **Topup**: Menambah saldo rekening.
    * **Transfer**: Mengirim dana antar rekening nasabah dengan validasi saldo.
    * **Withdraw**: Penarikan saldo dari rekening.
4.  **History**: Melihat riwayat mutasi transaksi (uang masuk/keluar).
5.  **Security**: Middleware untuk memvalidasi token pada endpoint sensitif.

## 🚀 Cara Menjalankan (Installation)

Ikuti langkah-langkah berikut untuk menjalankan proyek di komputer lokal Anda:

### 1. Clone Repository
```bash
git clone https://github.com/JuliusSinaga/UAS-Praktikum-PPW-2025-11.git
cd UAS-Praktikum-PPW-2025-11
````

### 2\. Setup Database

Pastikan PostgreSQL sudah berjalan, lalu buat database baru (misal: `banking_db`).

### 3\. Konfigurasi Environment

Salin file `.env.example` menjadi `.env` dan sesuaikan dengan kredensial database Anda:

```bash
cp .env.example .env
```

Isi file `.env`:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password_anda
DB_NAME=banking_db
PORT=8080
```

### 4\. Install Dependencies & Run

Jalankan perintah berikut untuk mengunduh library dan menjalankan server:

```bash
go mod tidy
go run main.go
```

Server akan berjalan di `http://localhost:8080`. Database akan otomatis melakukan migrasi tabel (`users`, `accounts`, `transactions`) saat aplikasi pertama kali dijalankan.

## 🔗 Daftar Endpoint Singkat

| Method | Endpoint | Deskripsi | Auth |
| :--- | :--- | :--- | :--- |
| `POST` | `/api/register` | Mendaftarkan user baru | Public |
| `POST` | `/api/login` | Login & dapatkan Token | Public |
| `GET` | `/api/balance` | Cek saldo user login | **Token** |
| `POST` | `/api/topup` | Isi ulang saldo | **Token** |
| `POST` | `/api/transfer` | Transfer sesama user | **Token** |
| `POST` | `/api/withdraw` | Tarik saldo | **Token** |
| `GET` | `/api/mutations` | Lihat riwayat transaksi | **Token** |

*Catatan: Dokumentasi request/response body lengkap dapat dilihat di [MiniBank-Postman](https://juliussinaga-4799400.postman.co/workspace/Julius-Sinaga's-Workspace~f4fe4bb7-58ea-4c1f-8d85-13716d3d02f6/collection/50633353-f8ea6499-c165-41e4-9ba4-3514498a2549?action=share&source=copy-link&creator=50633353).*

-----

© 2025 **Kelompok 11 PPW**.

