# Go Hexagon Project

## Deskripsi

Proyek ini adalah contoh implementasi Arsitektur Hexagonal dengan Golang menggunakan Go Fiber sebagai framework web dan Gorm untuk ORM dengan PostgreSQL.

## Instalasi Dependensi

1. Inisialisasi modul Go:

   ```sh
   go mod init go-hexagon
   ```

2. Instal Go Fiber:

   ```sh
   go get -u github.com/gofiber/fiber/v2
   ```

3. Instal GORM sebagai ORM untuk database SQL:

   ```sh
    go get -u gorm.io/gorm
    go get -u gorm.io/driver/postgres
   ```

4. Instal MongoDB driver:

   ```sh
    go get go.mongodb.org/mongo-driver/mongo
   ```

5. Mengelola dependensi:
   ```sh
   go mod tidy
   ```

## Cara Menjalankan Aplikasi

1. Pastikan PostgreSQL dan MongoDB berjalan.
2. Jalankan aplikasi:
   Ke PostgreSQL

   ```
   go run cmd\main.go --db=postgres
   ```

   Ke MongoDB

   ```
   go run cmd\main.go --db=mongodb
   ```

## API Endpoint

- GET /check-postgres - Cek koneksi ke PostgreSQL
- GET /check-mongo - Cek koneksi ke MongoDB
- GET /products - Mendapatkan daftar produk
- GET /products/:id - Mendapatkan detail produk berdasarkan ID
- POST /products - Membuat produk baru
- PUT /products/:id - Memperbarui produk berdasarkan ID
- DELETE /products/:id - Menghapus produk berdasarkan ID

## Penjelasan Arsitektur Hexagonal

Arsitektur Hexagonal (atau Ports and Adapters) adalah pola desain software yang memisahkan core logic dari infrastruktur dan framework. Ini memudahkan untuk mengubah komponen eksternal tanpa mengganggu logika inti.

## Keuntungan Menggunakan Arsitektur Hexagonal

1. Isolasi Logika Bisnis: Memisahkan logika bisnis dari detail teknis membuat aplikasi lebih mudah untuk dipelihara dan diubah.
2. Fleksibilitas: Memungkinkan penggunaan berbagai teknologi tanpa mempengaruhi logika inti aplikasi.
3. Pengujian Mudah: Memudahkan penulisan tes unit karena komponen-komponen yang berbeda dapat diuji secara terpisah.
4. Skalabilitas: Membantu dalam membangun aplikasi yang dapat diskalakan dengan memisahkan kekhawatiran antara komponen.

## Perbedaan dengan Arsitektur Lainnya

1. Monolitik:

- Karakteristik: Semua komponen aplikasi terintegrasi dalam satu kode besar.
- Kelemahan: Sulit untuk diskalakan dan dipelihara, terutama ketika aplikasi tumbuh besar.
- Perbandingan: Arsitektur Hexagonal lebih modular dan terisolasi, memudahkan pemeliharaan dan pengujian.

2. Layered (N-Tier) Architecture:

- Karakteristik: Memisahkan aplikasi menjadi lapisan-lapisan seperti presentasi, bisnis, dan data.
- Kelemahan: Ketergantungan antar lapisan bisa menjadi masalah, dan perubahan pada satu lapisan bisa mempengaruhi lapisan lain.
- Perbandingan: Arsitektur Hexagonal menawarkan isolasi yang lebih baik antara komponen dengan menggunakan port dan adapter.

3. Microservices:

- Karakteristik: Memecah aplikasi menjadi layanan-layanan kecil yang dapat dikelola dan diskalakan secara independen.
- Kelemahan: Kompleksitas dalam pengelolaan layanan-layanan yang terdistribusi.
- Perbandingan: Arsitektur Hexagonal dapat digunakan dalam konteks microservices untuk memastikan setiap layanan memiliki logika bisnis yang terisolasi.
