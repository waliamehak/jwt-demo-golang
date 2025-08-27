# jwt-demo-golang

A simple Go server demonstrating **JWT authentication** with Docker.  
This project shows how to create, sign, and validate JWT tokens in Go.


## Features
- Generate JWT tokens for login
- Validate JWT tokens for protected routes
- Fully containerized with Docker
- JWT secret is passed securely at runtime

## Prerequisites
- Go 1.25+ (for local development)
- Docker

## Run with Docker

Build the Docker image:

```bash
docker build -t jwt-demo-go .
```
Run the container with your secret:

```bash
docker run -p 8080:8080 -e JWT_SECRET='superSecretFromDocker!' jwt-demo-go
```
The server will run at http://localhost:8080.

## Test the API
Login to get a token:

```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username":"waliamehak","password":"password123"}'
```

Use the token to access the protected welcome route:

```bash
TOKEN=<paste-your-token-here>
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/welcome
```
You should see:

```nginx
Welcome waliamehak!
```

