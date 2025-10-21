# Use Alpine Linux with timezone support
FROM alpine:latest

# Install timezone data untuk Asia/Jakarta
RUN apk add --no-cache tzdata

# Set default timezone ke Asia/Jakarta
ENV TZ=Asia/Jakarta

COPY . .

# Expose port 21999, tempat aplikasi akan berjalan di dalam kontainer
EXPOSE 21999

# Menetapkan perintah default untuk menjalankan aplikasi saat kontainer dimulai
CMD ["/server"]
