# JWT-auth-GO
# Features

- **JWT Authentication**: Use JSON Web Tokens (JWTs) for secure user authentication.
- **Access and Refresh Tokens**: Support for both access and refresh tokens to manage user sessions effectively.
- **Token Management**: Custom database management for tracking and invalidating tokens, including the `IsDead` flag for immediate token invalidation.
- **User Authentication and Status**: Endpoints for user signup, login, token refresh, and status checking.
- **Middleware**: Authentication middleware to secure routes and verify token validity.

## Overview

This project combines two strategies for managing JWTs:

1. **Standard JWT Strategy**:
   - Tokens are self-contained and include expiration and other claims.
   - JWTs are validated independently using the `ParseToken` function.

2. **Custom Database Management**:
   - Tokens are stored in a database with an `IsDead` field to handle token revocation.
   - Access and refresh tokens are tracked in the database for additional control and flexibility.
   - Allows immediate invalidation of tokens if needed, such as during user logout or security events.

### Why Save Tokens in the Database?

By default, JWTs can become problematic if they are not managed carefully. Tokens are often designed to expire after a set period, but once issued, they cannot be disabled or invalidated until they expire. This poses a security risk if a token is compromised.

For example, if a user logs out and then logs in again, the system issues new refresh and access tokens. If someone has obtained the old refresh or access token, they could potentially use it to generate new access tokens, effectively having a backdoor to the user's data.

To address this, we use a database to manage tokens:

- **Token Storage**: Every time a user logs in, the system checks if they already have a token in the database. If so, it deletes the old token and inserts a new one. This ensures that only one set of tokens (access and refresh) is valid at any time.
- **Immediate Invalidations**: When a user logs out, both the access and refresh tokens are deleted from the database. This prevents any further use of those tokens.
- **Custom Expiration**: We use a longer expiration time for refresh tokens and a shorter one for access tokens. This approach allows for additional control and security, as access tokens are short-lived and less valuable if compromised.

## Endpoints

- **POST /signup**: Create a new user account.
- **POST /login**: Authenticate a user and issue access and refresh tokens.
- **POST /refresh-token**: Generate a new access token using a refresh token.
- **POST /logout**: Invalidate the current tokens and mark them as dead.
- **GET /check-status**: Check if the user is online.
- **GET /user-info**: Retrieve user information using the access token.
- **GET /admin**: For users who have admin role.

## Token Management

- **Access Tokens**: Used for short-lived authentication.
- **Refresh Tokens**: Used to obtain new access tokens without re-authentication.
- **Database Management**: Tokens are stored in the database with an `IsDead` flag to manage and invalidate tokens effectively.


## To install the necessary dependencies for this project, run the following commands:
1. GO
2. 
```bash
- go get github.com/githubnemo/CompileDaemon
- go install github.com/githubnemo/CompileDaemon
- go get github.com/joho/godotenv
- go get -u github.com/gin-gonic/gin
- go get -u gorm.io/gorm
- go get -u gorm.io/driver/postgres
- go get github.com/jackc/pgx/v5
- go get -u golang.org/x/crypto/bcrypt
```
and for run use:
```bash
CompileDaemon -command="./JWTauth"
```