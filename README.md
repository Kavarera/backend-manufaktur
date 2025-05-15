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
1. `/regitser` Super Admin only

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
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InN1cGVyX2FkbWluX3VzZXIiLCJyb2xlIjoic3VwZXJfYWRtaW5cbiIsImV4cCI6MTczMzk5MTMzN30.xtsawJm2U2Q8RxutPPiECyhewPWCNQk0PgPT9c7Y8BE",
    "user": {
        "fullname": "superAdmin1",
        "role": "super_admin\n",
        "username": "super_admin_user"
    }
}
```
