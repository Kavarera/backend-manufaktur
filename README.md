# API Manufaktur

This project is a RESTful API built using the [Gin](https://github.com/gin-gonic/gin) framework in Golang.

## Project Structure
```
testing-api/ 
│   main.go
|   go.mod
|   go.sum
└───db/
    │   db.go
└───handlers/
    │   auth.go
└───middlewares/
    │   auth.go
└───models/
└───utils/
    │   jwt.go
```

### Endpoints
1. POST `/register`
2. POST `/login`
3. GET `/users/:id` 
4. DELETE `/users/:id` 
5. POST `/barangMentah`
6. GET `/barangMentah`
7. PUT `/barangMentah/:id`
8. DELETE `/barangMentah/id`
9. POST `/barangProduksi`
10. GET `/barangProduksi`
11. GET `/barangProduksi/:id`
12. PUT `/barangProduksi/:id`
13. DELETE `/barangProduksi/:id`
14. POST `/gudang`
15. GET `/gudang`
16. GET `/gudang/:id`
17. PUT `/gudang/:id`
18. DELETE `/gudang/:id`
19. POST `/rencanaProduksi`
20. GET `/rencanaProduksi`
21. GET `/rencanaProduksi/:id`
22. PUT `/rencanaProduksi/:id`
23. DELETE `/rencanaProduksi/:id`

### Responses 

1. Register User

**HTTP Request:**
```
POST /register
```

Payload Body JSON: 
```
{
  "id": "admin1",
  "username": "Test2",
  "password": "12345678910",
  "hak_akses": 7
}
```

Response: 200 SUCCESS
```
{
    "message": "User registered successfully"
}
```

2. Login Super Admin Method

**HTTP Request:**
```
POST /login
```

Payload Body JSON: 
```
{
  "username": "Test4",
  "password": "abcdefghijk"
}
```

Response: 200 SUCCESS
```
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6IlRlc3Q0Iiwicm9sZSI6IlN1cGVyQWRtaW4iLCJleHAiOjE3NDc0MTQ5MjZ9.BeyfcFbcTK1zd1kPRceJvdK5C7AiiCPKun6F1ZnIMxk",
    "user": {
        "roleName": "SuperAdmin",
        "username": "Test4"
    }
}
```

3. Menampilkan List User

**HTTP Request:**
```
GET /users/1
```

Response: 200 SUCCESS
```
{
    "id": "1",
    "username": "Test1",
    "password": "12345678910",
    "hak_akses": "SuperAdmin",
    "hak_id": 7
}
```

4. Hapus User (Super Admin Only)

**HTTP Request:**
```
DELETE /users/1
```

Response: 200 SUCCESS
```
{
    "message": "User deleted successfully"
}
```

5. Tambah Barang Mentah

**HTTP Request:**
```
POST  /barangMentah
```
Payload Body JSON: 
```
{
  "nama": "Example Item",
  "kodeBarang": "EX123",
  "hargaStandar": 10000.5,
  "satuanId": 1,
  "stok": 50,
  "gudangId": 1
}
```

Response: 200 SUCCESS
```
{
    "id": 2,
    "nama": "Jamban Wangi",
    "kodeBarang": "EX124",
    "hargaStandar": 10000.5,
    "satuanId": 1,
    "satuanNama": "",
    "stok": 50,
    "gudangId": 1,
    "gudangNama": ""
}
```

6. List Barang Mentah

**HTTP Request:**
```
GET /barangMentah
```

Response: 200 SUCCESS
```
[
    {
        "id": 1,
        "nama": "Example Item",
        "kodeBarang": "EX123",
        "hargaStandar": 10000.5,
        "satuanId": 1,
        "satuanNama": "jamban turunan",
        "stok": 50,
        "gudangId": 1,
        "gudangNama": "gudang1"
    },
    {
        "id": 2,
        "nama": "Jamban Wangi",
        "kodeBarang": "EX124",
        "hargaStandar": 10000.5,
        "satuanId": 1,
        "satuanNama": "jamban turunan",
        "stok": 50,
        "gudangId": 1,
        "gudangNama": "gudang1"
    }
]
```

7. Update Barang Mentah

**HTTP Request:**
```
PUT /barangMentah/1
```

Payload Body JSON: 
```
{
  "nama": "Jamban Anjay",
  "kodeBarang": "EX124",
  "hargaStandar": 20000.5,
  "satuanId": 1,
  "stok": 100,
  "gudangId": 1
}
```

Response: 200 SUCCESS
```
{
    "message": "Barang updated successfully"
}
```

8. Hapus Barang Mentah

**HTTP Request:**
```
DELETE /barangMentah/1
```


Response: 200 SUCCESS
```
{
    "message": "Barang deleted successfully"
}
```

9. Tambah Barang Produksi

**HTTP Request:**
```
POST  /barangProduksi
```
Payload Body JSON: 
```
{
  "nama": "Produk Contoh",
  "kodeBarang": "PRD123",
  "hargaStandar": 20000.5,
  "hargaReal": 19000.75,
  "satuanId": 1,
  "stok": 100,
  "gudangId": 1
}

```

Response: 200 SUCCESS
```
{
    "id": 8,
    "nama": "Produk Contoh",
    "kodeBarang": "PRD123",
    "hargaStandar": 20000.5,
    "hargaReal": 19000.75,
    "satuanId": 1,
    "satuanNama": "",
    "stok": 100,
    "gudangId": 1,
    "gudangNama": ""
}
```

10. List Barang Produksi

**HTTP Request:**
```
GET /barangProduksi
```

Response: 200 SUCCESS
```
[
    {
        "id": 5,
        "nama": "Produk Contoh",
        "kodeBarang": "PRD123",
        "hargaStandar": 20000,
        "hargaReal": 19000,
        "satuanId": 1,
        "satuanNama": "jamban turunan",
        "stok": 100,
        "gudangId": 1,
        "gudangNama": "gudang1"
    },
    {
        "id": 8,
        "nama": "Produk Contoh 2",
        "kodeBarang": "PRD124",
        "hargaStandar": 20000,
        "hargaReal": 19000,
        "satuanId": 1,
        "satuanNama": "jamban turunan",
        "stok": 100,
        "gudangId": 1,
        "gudangNama": "gudang1"
    }
]
```

11. List Barang Produksi by ID

**HTTP Request:**
```
GET /barangProduksi/5
```

Response: 200 SUCCESS
```
{
    "id": 5,
    "nama": "Produk Contoh",
    "kodeBarang": "PRD123",
    "hargaStandar": 20000,
    "hargaReal": 19000,
    "satuanId": 1,
    "satuanNama": "jamban turunan",
    "stok": 100,
    "gudangId": 1,
    "gudangNama": "gudang1"
}
```

12. Update Barang Produksi

**HTTP Request:**
```
PUT /barangProduksi/8
```

Payload Body JSON: 
```
{
  "nama": "Produk Contoh",
  "kodeBarang": "PRD999,
  "hargaStandar": 1000000,
  "hargaReal": 1100000,
  "satuanId": 1,
  "stok": 10,
  "gudangId": 1
}

```

Response: 200 SUCCESS
```
{
    "message": "Barang produksi updated successfully"
}
```

13. Hapus Barang Mentah

**HTTP Request:**
```
DELETE /barangProduksi/8
```


Response: 200 SUCCESS
```
{
    "message": "Barang produksi deleted successfully"
}
```

14. Tambah Gudang

**HTTP Request:**
```
POST  /gudang
```
Payload Body JSON: 
```
{
  "nama": "gudang2",
}

```

Response: 200 SUCCESS
```
{
    "id": 2,
    "nama": "gudang2"
}
```

15. List Gudang

**HTTP Request:**
```
GET /gudang
```

Response: 200 SUCCESS
```
[
    {
        "id": 1,
        "nama": "gudang1"
    },
    {
        "id": 2,
        "nama": "gudang2"
    }
]
```

16. List Gudang by ID

**HTTP Request:**
```
GET /gudang/2
```

Response: 200 SUCCESS
```
{
    "id": 2,
    "nama": "gudang2"
}
```

17. Update Gudang

**HTTP Request:**
```
PUT /gudang/2
```

Payload Body JSON: 
```
{
  "nama": "gudang kedua",
}

```

Response: 200 SUCCESS
```
{
    "message": "Gudang updated successfully"
}

```

18. Hapus Gudang

**HTTP Request:**
```
DELETE /barangProduksi/8
```


Response: 200 SUCCESS
```
{
    "message": "Gudang deleted successfully"
}
```

19. Tambah Rencana Produksi

**HTTP Request:**
```
POST  /rencanaProduksi
```
Payload Body JSON: 
```
{
  "id": "RP001",
  "barangProduksiId": 5,
  "tanggalMulai": "2025-06-01",
  "waktuMulai": "08:00:00Z",
  "tanggalSelesai": "2025-06-15",
  "waktuSelesai": "17:00:00Z"
}

```

Response: 200 SUCCESS
```
{
    "barangProduksiId": 5,
    "id": "RP00111",
    "tanggalMulai": "2025-06-01T08:00:00Z",
    "tanggalSelesai": "2025-06-15T17:00:00Z"
}
```

20. List Rencana Produksi

**HTTP Request:**
```
GET /rencanaProduksi
```

Response: 200 SUCCESS
```
[
    {
        "id": "RP001",
        "barangProduksiId": 5,
        "tanggalMulai": "2025-06-01T08:00:00Z",
        "tanggalSelesai": "2025-06-15T17:00:00Z"
    },
    {
        "id": "RP0010",
        "barangProduksiId": 5,
        "tanggalMulai": "2025-06-01T08:00:00Z",
        "tanggalSelesai": "2025-06-15T17:00:00Z"
    },
    {
        "id": "RP00111",
        "barangProduksiId": 5,
        "tanggalMulai": "2025-06-01T08:00:00Z",
        "tanggalSelesai": "2025-06-15T17:00:00Z"
    }
]
```

21. List Rencana Produksi by ID

**HTTP Request:**
```
GET /gudang/RP00111
```

Response: 200 SUCCESS
```
{
    "id": "RP00111",
    "barangProduksiId": 5,
    "tanggalMulai": "2025-06-01T08:00:00Z",
    "tanggalSelesai": "2025-06-15T17:00:00Z"
}
```

22. Update Rencana Produksi

**HTTP Request:**
```
PUT /rencanaProduksi/RP00111
```

Payload Body JSON: (Dinamis)
```
{
  "tanggalMulai": "2025-06-10T09:00:00Z"
}

atau

{
  "tanggalSelesai": "2025-06-10T09:00:00Z"
}

```

Response: 200 SUCCESS
```
{
    "message": "Rencana produksi updated successfully"
}
```

23. Hapus Rencana Produksi

**HTTP Request:**
```
DELETE /barangProduksi/RP00111
```


Response: 200 SUCCESS
```
{
    "message": "Rencana produksi deleted successfully"
}
```




