# Authentication Back-End
This project focuses on all the core mechanics of user authentication. It is a back-end API that stores data on a local SQLite database built in Go.

## Features
### 🔐 API Key Protection
Each request needs to include an API key in the header or the request get's rejected by the middleware.

### 🪙 Token Authentication
Supports JWT-baed authentication with seperate tokens for both short lived access tokens as well as longer-lives refresh tokens.
Ensures token integrity with validation and a secure refresh endpoint for seamless front-end authentication processes.

### 👤 User Management
Creates users with a secure BCrypt hashed password. Can validate with the credentials or the token authentication.

### 🔥 Error Handling
Consistent JSON error responses that are descriptive without giving any sensitive information away. Also logs
internal server errors for further debugging.

## Endpoints
### /api/auth/validate-refresh (POST)
Verifies whether a refresh token is valid.

### /api/auth/refresh-token (POST)
Regenerates a new JWT from the refresh token.

### /api/auth/users (POST)
Creates a new user object and returns a JWT.

### /api/auth/verify (POST)
Validates either a user's credentials or a JWT.
