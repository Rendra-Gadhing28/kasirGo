# PANDUAN DEPLOY GRATIS 100% TANPA KARTU KREDIT (VERCEL FULL-STACK)

Semua komponen berjalan di **Vercel** dan **TiDB Cloud**:
- **Database**: TiDB Cloud Serverless (Free MySQL 5 GB, Tanpa Kartu).
- **Backend (Go Fiber)**: Vercel Serverless Function (Free, Tanpa Kartu).
- **Frontend (SvelteKit)**: Vercel (Free, Tanpa Kartu).

---

## TAHAP 1: DEPLOY BACKEND GO DI VERCEL (PROJECT KE-2)

1. Pastikan repo `kasir-go` sudah di-push ke GitHub.
2. Buka [https://vercel.com](https://vercel.com) dan login via GitHub.
3. Klik **"Add New..."** -> **"Project"**.
4. Import repository `kasir-go`.
5. Atur konfigurasi:
   - **Project Name**: `kasirpro-api` (atau nama lain)
   - **Framework Preset**: `Other`
   - **Root Directory**: Klik **"Edit"**, pilih folder **`backend`**, lalu klik **"Continue"**.
6. Buka bagian **"Environment Variables"**, masukkan:
   | Key | Value |
   |---|---|
   | `DB_HOST` | *(host TiDB Cloud Anda)* |
   | `DB_PORT` | `4000` |
   | `DB_USER` | *(user TiDB Cloud Anda)* |
   | `DB_PASSWORD` | *(password TiDB Cloud Anda)* |
   | `DB_NAME` | `kasirpro` |
   | `DB_SSL` | `true` |
   | `JWT_SECRET` | `f8b4c29e71a0653d9e832104bf56dc7a18e034927f15ac6389ebd0412758bc6a` |
   | `CORS_ALLOWED_ORIGINS` | `*` |
   | `MIDTRANS_SERVER_KEY` | `SB-Mid-server-simulator-mock-key-12345` |
   | `MIDTRANS_CLIENT_KEY` | `SB-Mid-client-simulator-mock-key-12345` |
   | `MIDTRANS_IS_PRODUCTION` | `false` |

7. Klik **"Deploy"**.
8. Selesai! Vercel otomatis mengompilasi Go binary dan menghasilkan URL backend gratis, contoh:  
   👉 `https://kasirpro-api.vercel.app`
9. Uji di browser: Buka `https://kasirpro-api.vercel.app/api/v1/health` -> keluar JSON `{"success": true, ...}`.

---

## TAHAP 2: HUBUNGKAN KE FRONTEND VERCEL

1. Di dashboard Vercel, buka proyek **Frontend** KasirPro Anda.
2. Buka **Settings** -> **Environment Variables**.
3. Tambahkan / Update:
   - **Key**: `BACKEND_URL`
   - **Value**: `https://kasirpro-api.vercel.app` *(URL backend dari Tahap 1)*
4. Masuk ke tab **Deployments** -> Klik titik tiga deploy terbaru -> **Redeploy**.

---

## HASIL
- Frontend & Backend aktif 24/7 di Vercel.
- Database aman di TiDB Cloud.
- **100% Bebas Biaya & Tanpa Kartu Kredit.**
- Laptop bebas dimatikan.
