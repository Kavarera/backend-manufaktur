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
1. `/register` Super Admin only

### Responses 

#### 1. Login Super Admin Method
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
