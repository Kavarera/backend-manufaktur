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

## Endpoints

The API provides main endpoints:

1. **POST `/login`** - LOGIN

Endpoint:
POST /login

Description:
Authenticates a user with username and password. Returns a JWT token if credentials are valid.

Request Body (JSON):
{
  "username": "string",
  "password": "string"
}
Success Response:

Status: 200 OK
