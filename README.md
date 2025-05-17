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
Endpoints
1. POST `/register`
2. POST `/login`
3. GET `/users/:id` 
4. DELETE `/users/:id` 
5. POST `/barangMentah`
6. GET `/barangMentah`
7. PUT `/barangMentah/:id`
8. DELETE `/barangMentah:id`

### Responses 

1. Register New User
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

3. Login Super Admin Method
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

4. Login Super Admin Method
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

