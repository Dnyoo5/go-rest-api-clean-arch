# 1. Ambil bahan dasar: Go versi 1.23 (Sesuaikan versi go.mod kamu)
# Kita pakai versi 'alpine' yang ukurannya kecil banget (ringan)
FROM golang:1.25-alpine

# 2. Bikin folder kerja di dalam kontener
WORKDIR /app

# 3. Copy file dependency dulu (biar cache-nya awet)
COPY go.mod go.sum ./

# 4. Download semua library (kayak 'go mod download')
RUN go mod download

# 5. Copy sisa semua kodingan kamu ke dalam kontener
COPY . .

# 6. Build aplikasi jadi file binary bernama 'main'
RUN go build -o main ./cmd/api/main.go

# 7. Beri tahu bahwa aplikasi ini pakai port 8080
EXPOSE 8080

# 8. Perintah terakhir: Jalankan aplikasi!
CMD ["./main"]