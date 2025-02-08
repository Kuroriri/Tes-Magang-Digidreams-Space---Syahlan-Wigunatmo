# Cara Menjalankan API serta mengujinya dengan Postman (Windows)
Requirements
- Code Editor
- Docker Desktop
- Postman Desktop

1. Menjalankan API
    - Clone repository ini pada device anda
    -  Masuk ke direktori proyek dan buka terminal
    -  Buka docker desktop pada device anda
    -  Jalankan perintah "docker-compose up -d" untuk menjalankan PostgreSQL dengan Docker
    -  Jalankan perintah "go run main.go" untuk menjalankan aplikasi
      
2. Menguji API dengan Postman
    - Buka Postman Desktop pada device anda
    - Pada bagian collections, pilih import kemudian pilih file postman_collection.json yang tersedia pada folder proyek
    - jalankan 2 pengujian awal, yaitu Register User dan Login User. pada pengujian Login User, simpan JWT token yang di generate
    - Pilih ikon tiga titik (⋮) di sebelah nama Collection, lalu pilih Edit
    - Masuk ke Tab "Variables"
    - Masukkan token di kolom CURRENT VALUE untuk semua variable lalu Klik Save
    - Jalankan kembali pengujian API
   


      

  


      
