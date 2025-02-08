# Tes-Magang-Digidreams-Space---Syahlan-Wigunatmo

Requirements
  - Code Editor
  - Docker Desktop
  - Postman Desktop

1. Menjalankan API
   - Clone repository pada device anda
   - Buka Docker Desktop, kemudian pada terminal code editor anda, jalankan perintah "docker-compose up -d" untuk menjalankan Postgre SQL dengan Docker
   - jalankan perintah "go run main.go" untuk menjalankan aplikasi
     
2. Pengujian API
   - Pengujian 5 API awal
       - Buka Postman Desktop
       - Pada bagian collections, import postman_collection.json yang terdapat di folder proyek
       - Uji 2 API awal yaitu Register User dan Login User
       - simpan token JWT yang dihasilkan ketika menjalankan API Login User
       - Klik ikon tiga titik (⋮) di sebelah nama Collection, lalu pilih Edit
       - Masuk ke Tab "Variables"
       - Tambahkan nilai jwt_token dengan mengisi pada bagian current value, klik save
       - Uji 3 API berikutnya
   
   - Pengujian API Delete User
       - buka terminal pada code editor, gunakan CTRL+C untuk menghentikan aplikasi
       - Jalankan perintah "docker ps" untuk melihat container PostgreSQL yang berjalan
       - Jalankan perintah "docker exec -it <container_id> psql -U postgres" untuk masuk ke dalam container. ganti container_id dengan container ID yang berjalan
       - Jalankan perintah "\c users" untuk memilih database yang digunakan oleh aplikasi
       - Jalankan perintah "\dt" untuk melihat daftar tabel
       - Jalankan perintah "SELECT * FROM users;" untuk melihat semua data pada tabel
       - Jalankan perintah "UPDATE users SET role = 'admin' WHERE email = 'sofia@example.com';" untuk mengupdate role dari data tersebut
       - Jalankan perintah "\q" untuk keluar dari container
       - Jalankan kembali aplikasi dengan perintah "go run main.go"
       - Lakukan pengujian API terakhir di Postman

   
 
