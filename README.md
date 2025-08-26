# jwt-demo-golang

Tiny Go demo showing JWT authentication (login and protected endpoint).

## What it is
- POST `/login` with `{ "username": "admin", "password": "password" }` → returns a JWT token.
- GET `/protected` with header `Authorization: Bearer <token>` → returns a protected message.

## Quick start
I’ll add runnable instructions, Dockerfile, CI, Postman collection and demo GIF in subsequent commits.
