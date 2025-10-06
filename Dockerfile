# Stage 1: Build the Go application
# Menggunakan image Go berbasis Alpine Linux yang ringan sebagai builder
FROM golang:1.24-alpine AS builder

# Menetapkan variabel lingkungan untuk build statis yang kompatibel dengan Linux
ENV CGO_ENABLED=0 GOOS=linux
WORKDIR /app

# Menyalin file dependensi terlebih dahulu untuk memanfaatkan cache Docker
COPY go.mod go.sum ./

# Mengunduh dependensi
RUN go mod download

# Menyalin seluruh kode sumber aplikasi
COPY . .

# Meng-compile aplikasi Go menjadi satu file biner statis.
# Flag -ldflags="-s -w" menghapus informasi debug untuk memperkecil ukuran file biner.
RUN go build -o /server -ldflags="-s -w" .

# ---

# Stage 2: Create a minimal final image for production
# Memulai dari image 'scratch' yang benar-benar kosong untuk ukuran minimal dan keamanan maksimal
FROM scratch

# Menyalin hanya file biner yang telah di-compile dari stage builder
COPY --from=builder /server /server

# Mengekspos port 21999, tempat aplikasi akan berjalan di dalam kontainer
EXPOSE 21999

# Menetapkan perintah default untuk menjalankan aplikasi saat kontainer dimulai
CMD ["/server"]
