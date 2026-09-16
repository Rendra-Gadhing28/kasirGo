# KasirPro — Sistem Kasir Modern Indonesia (Full-Stack)

Sistem Kasir (Point of Sale) modern, responsif, scalable, dan aman yang dirancang khusus untuk ekosistem retail dan pasar modern Indonesia. Mendukung format Rupiah, perhitungan PPN 11%, barcode scanner, struk thermal format ESC/POS, dan pembayaran dinamis **QRIS via Midtrans Sandbox / Simulator**.

Tampilan antarmuka mengusung gaya estetika **Neo-Brutalism UI** (border tebal solid, hard drop shadow, warna kontras tinggi, typography bold) dengan dukungan penuh **Dark Mode & Light Mode**.

---

## 1. Arsitektur & Tech Stack

| Layer | Teknologi |
|---|---|
| **Backend** | Go (Golang) 1.22+ dengan **Fiber v2**, **GORM**, dan **golang-jwt/v5** |
| **Frontend** | **SvelteKit 2** + **TypeScript**, **Tailwind CSS**, dan custom Neo-Brutalism Design Tokens |
| **Database** | **MySQL 8** (InnoDB, utf8mb4, parameterized query penuh, user privilege terpisah) |
| **Payment Gateway** | **Midtrans Simulator / Sandbox** (QRIS Dynamic Payment & Status Confirmation) |
| **Reverse Proxy & LB** | **Nginx** (Load balancing `least_conn`, IP rate limiting, security headers) |
| **Containerization** | **Docker & Docker Compose** multi-container orchestration |

---

## 2. Standar Keamanan Enterprise yang Diterapkan

1. **XSS Prevention**:
   - Security headers otomatis melalui middleware Fiber Helmet (CSP, `X-Content-Type-Options: nosniff`, `X-Frame-Options: DENY`).
   - Token JWT disimpan aman di cookie `HttpOnly`, `SameSite=Lax`, dan flag `Secure`.
   - Larangan rendering raw HTML tanpa sanitasi ketat di frontend Svelte.
2. **SQL Injection Prevention**:
   - 100% Parameterized query via GORM (`db.Where("outlet_id = ? AND id = ?", ...)`).
   - Larangan konkatenasi string untuk query database.
   - User database non-root (`kasirpro_app`) hanya diberikan hak akses DML (`SELECT`, `INSERT`, `UPDATE`, `DELETE`).
3. **RBAC (Role-Based Access Control)**:
   - Pembagian role hierarkis:
     - `owner`: Akses penuh semua modul, manajemen user, dan otorisasi diskon besar.
     - `admin`: Manajemen katalog produk, kategori, pelanggan, supplier, dan retur/void.
     - `kasir`: Akses operasional kasir (POS), pencarian produk read-only, dan riwayat transaksi shift pribadi.
4. **ABAC (Attribute-Based Access Control)**:
   - Tenant isolation wajib via `outlet_id` di seluruh query data.
   - Kasir dibatasi hanya dapat membaca histori transaksi miliknya (`user_id = current_user`).
   - Kebijakan diskon: Diskon transaksi di atas 20% membutuhkan role `owner`.
5. **Sanitasi Input**:
   - Sanitasi string berbahaya / script injection menggunakan library `bluemonday` di Go.
   - Pembersihan nama file upload untuk mencegah ancaman path traversal.
6. **Double Validation**:
   - Frontend: Validasi form interaktif real-time.
   - Backend: Validasi struct menyeluruh dengan `go-playground/validator/v10`.
7. **Rate Limiting**:
   - Nginx layer: Pembatasan per IP (`auth`: 5 req/menit, `api`: 100 req/detik).
   - Fiber layer: Auth rate limiter (5 req/menit), POS checkout limiter (30 req/menit), API umum (60 req/menit).
8. **Load Balancing & Stateless Architecture**:
   - Nginx upstream `least_conn` siap horizontal scaling instance backend.
   - Backend stateless berbasis validasi tanda tangan JWT.
9. **Caching Strategy**:
   - In-memory cache layer (`go-cache`) untuk data kategori dan pengaturan toko.
   - Cache invalidation otomatis saat terjadi mutasi data (create/update/delete).

---

## 3. Fitur Utama (Fase 1 - 3)

### Point of Sale (POS / Kasir)
- Antarmuka kasir 2-kolom responsif (Katalog di kiri, Keranjang di kanan; switch tab di mobile).
- Pencarian produk instan berdasarkan Nama, SKU, atau pemindaian Barcode (EAN-13/UPC).
- Filter produk cepat berdasarkan Kategori barang.
- Keranjang belanja interaktif: ubah kuantitas, batasan stok real-time, diskon per item (% atau nominal Rp).
- Diskon transaksi global dan switch opsional kalkulasi pajak PPN 11%.
- Multi-metode pembayaran:
  - **Tunai (Cash)**: Input nominal diterima, tombol cepat uang pas / pecahan Rp20k - Rp200k, kalkulasi kembalian otomatis.
  - **QRIS (Midtrans Sandbox / Simulator)**: Dynamic QR code generation, simulasi pelunasan langsung.
  - **Transfer Bank**.
- Integrasi pelanggan & program poin loyalitas (+1 poin tiap kelipatan Rp 10.000).
- Struk thermal format 58mm/80mm siap cetak langsung via browser (`window.print`).

### Fitur Pendaftaran Mandiri (Self-Registration)
- **Pendaftaran Mandiri Karyawan / Kasir**:
  - Calon kasir/karyawan membuka halaman `/register` dan memilih tab **"Gabung Toko (Kasir / Staff)"**.
  - Masukkan **Kode Toko** (contoh: `KASIR-DEMO-001`), nama, email, dan password.
  - Akun langsung terhubung ke outlet toko terkait dengan role `kasir` tanpa perlu di-input manual oleh Owner.
- **Pendaftaran Member Mandiri oleh Pelanggan (QR Meja Kasir)**:
  - Tersedia tombol **"📱 CETAK QR MEMBER"** di halaman POS dan Pelanggan.
  - Owner/Kasir mencetak standee QR dan menempelkannya di meja kasir.
  - Pembeli memindai QR via kamera HP, membuka tautan publik `https://fish-warming-logos-lots.trycloudflare.com/join-member/KASIR-DEMO-001`.
  - Pembeli mengisi nama dan no. WhatsApp, instan mendapatkan **10 Poin Bonus** dan kartu member digital dengan Barcode nomor HP yang siap di-scan kasir.

### Master Data Management
- **Produk**: Nama, SKU, Barcode, Kategori, Harga Beli, Harga Jual, Margin Laba, Stok, Satuan, status aktif. Fitur auto-generate barcode & SKU.
- **Kategori**: Pengelompokan jenis barang toko dengan in-memory caching.
- **Pelanggan**: Kontak WhatsApp/telepon, email, alamat, dan akumulasi poin belanja.
- **Supplier**: Data distributor, kontak person (sales/PIC), telepon, dan alamat gudang.

### Laporan & Riwayat Transaksi
- Dashboard analytics dengan metrik Omset Hari Ini, Jumlah Transaksi, Total SKU, Peringatan Stok Menipis.
- Grafik garis/bar tren pendapatan 7 hari terakhir.
- Kontribusi kategori terlaris dan rincian transaksi per metode pembayaran.
- Riwayat transaksi dengan filter tanggal, pencarian no faktur (INV), dan status bayar.
- Fitur **Void / Retur Transaksi** khusus Admin & Owner yang secara otomatis mengembalikan stok barang ke inventaris.

---

## 4. Akun Demo Bawaan (Default Seeder)

Database otomatis diisi data demo saat inisialisasi awal:

| Role | Email | Password |
|---|---|---|
| **Owner** | `owner@kasirpro.id` | `password123` |
| **Admin** | `admin@kasirpro.id` | `password123` |
| **Kasir** | `kasir@kasirpro.id` | `password123` |

Tersedia tombol pintas "Akun Demo Cepat" pada halaman login.

---

## 5. Cara Menjalankan dengan Docker Compose

Pastikan Docker Desktop / Docker Engine telah terpasang dan aktif.

### Jalankan Seluruh Sistem:
```bash
docker compose up -d --build
```

Setelah kontainer berjalan:
- **Aplikasi Kasir (Web)**: Buka browser di [https://fish-warming-logos-lots.trycloudflare.com](https://fish-warming-logos-lots.trycloudflare.com) (atau [http://localhost](http://localhost))
- **Backend API Direct**: [http://localhost:8080/api/v1/health](http://localhost:8080/api/v1/health)
- **Database MySQL**: `localhost:3306` (Database: `kasirpro`, User: `root` / `kasirpro_app`)

### Menghentikan Sistem:
```bash
docker compose down
```

---

## 6. Menjalankan untuk Pengembangan Lokal (Local Development)

Jika ingin menjalankan service secara terpisah di host lokal:

### 1. Database MySQL
Pastikan MySQL 8 berjalan di port 3306 dengan konfigurasi di `.env`.

### 2. Backend Go
```bash
cd backend
go run cmd/server/main.go
```
Backend akan berjalan di `http://localhost:8080`.

### 3. Frontend SvelteKit
```bash
cd frontend
npm install
npm run dev
```
Frontend Vite dev server akan berjalan di `http://localhost:5173`.
Proxy API diarahkan otomatis ke `http://localhost:8080`.

---

## 7. Struktur Proyek

```
kasir-go/
├── docker-compose.yml       # Definisi multi-container (MySQL, Backend, Frontend, Nginx)
├── Makefile                 # Shortcut make up, down, restart, logs
├── .env.example             # Template konfigurasi environment
├── nginx/
│   └── nginx.conf           # Load balancer, rate limit, dan security proxy
├── mysql/
│   └── init.sql             # Skrip inisialisasi user DB least-privilege
├── backend/
│   ├── cmd/server/main.go   # Entrypoint HTTP server Fiber
│   ├── internal/
│   │   ├── config/          # Loader variabel konfigurasi
│   │   ├── database/        # Koneksi GORM & Seeder demo
│   │   ├── dto/             # Request & Response Data Transfer Objects
│   │   ├── handler/         # Controller HTTP Fiber
│   │   ├── middleware/      # Security, JWT Auth, RBAC, ABAC, Rate Limiting
│   │   ├── model/           # Schema database GORM
│   │   ├── repository/      # Data access layer (parameterized queries)
│   │   ├── service/         # Business logic, transaksi DB, Midtrans QRIS
│   │   └── utils/           # Hash bcrypt, JWT, bluemonday sanitizer, go-cache
│   └── Dockerfile
└── frontend/
    ├── src/
    │   ├── app.css          # Neo-Brutalism design tokens & print styles
    │   ├── lib/
    │   │   ├── api/client.ts      # HTTP client fetch wrapper
    │   │   ├── components/
    │   │   │   ├── layout/        # Sidebar & Navbar
    │   │   │   ├── pos/           # CheckoutModal, ThermalReceipt
    │   │   │   └── ui/            # Button, Input, Modal, Badge, Toast
    │   │   ├── stores/            # Auth store, POS cart store, Theme store, Toast
    │   │   ├── types/             # TypeScript models & interfaces
    │   │   └── utils/             # Format Rupiah & Date
    │   └── routes/
    │       ├── login/             # Login dengan demo switcher
    │       ├── register/          # Buka toko & register owner
    │       ├── dashboard/         # Metrik, grafik 7 hari, sales breakdown
    │       ├── pos/               # Halaman kasir, katalog, cart, bayar
    │       ├── products/          # Manajemen inventaris produk
    │       ├── categories/        # Manajemen kategori
    │       ├── customers/         # Member & poin loyalitas
    │       ├── suppliers/         # Data distributor
    │       └── transactions/      # Riwayat transaksi & fitur void
    └── Dockerfile
```
