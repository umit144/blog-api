# Go Blog API

A powerful and flexible blog API built with Go, featuring authentication, post management, and more.

## Technology Stack

- **Language**: Go
- **Database**: PostgreSQL
- **Web Framework**: Fiber
- **Authentication**: JWT (JSON Web Tokens)
- **OAuth**: Google OAuth support
- **Documentation**: OpenAPI (Swagger)

## Features

- User registration and authentication
- Blog post CRUD operations
- Category management
- File upload functionality
- Google OAuth integration
- OpenAPI documentation

## Getting Started

### Prerequisites

- Go (latest version recommended)
- PostgreSQL
- Make (for using Makefile commands)

### Installation

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/go-blog.git
   ```

2. Navigate to the project directory:
   ```
   cd go-blog
   ```

3. Copy the example environment file and update it with your configurations:
   ```
   cp .env.example .env
   ```

4. Update the `.env` file with your specific configurations.

5. Install dependencies:
   ```
   go mod tidy
   ```

### Running the Application

Use the following Make commands to manage the application:

- Build the application:
  ```
  make build
  ```

- Run the application:
  ```
  make run
  ```

- Create DB container:
  ```
  make docker-run
  ```

- Shutdown DB container:
  ```
  make docker-down
  ```

- Live reload the application:
  ```
  make watch
  ```

- Run tests:
  ```
  make test
  ```

- Clean up binaries:
  ```
  make clean
  ```

- Run all make commands with clean tests:
  ```
  make all build
  ```

## API Endpoints

The API provides the following main route groups:

- `/api/auth`: Authentication routes
- `/api/users`: User management
- `/api/posts`: Blog post operations
- `/api/categories`: Category management
- `/api/files`: File upload and management

For a complete list of endpoints and their descriptions, refer to the OpenAPI documentation available at `/swagger` when the server is running.

## Authentication

The API uses JWT for authentication. Most endpoints require a valid JWT token, which should be included in the `access_token` cookie.

## Documentation

OpenAPI documentation is available. When the server is running, visit `/swagger` to explore the API documentation interactively.

## Environment Variables

Create a `.env` file in the root directory with the following variables:

```
PORT=8080
APP_ENV=local

DB_HOST=localhost
DB_PORT=5432
DB_DATABASE=postgres
DB_USERNAME=postgres
DB_PASSWORD=postgres
DB_SCHEMA=public

JWT_SECRET=""

GOOGLE_CLIENT_ID=""
GOOGLE_CLIENT_SECRET=""
OAUTH_REDIRECT_URL="/auth/google/callback"
CLIENT_URL="http://localhost:3000"
```

Adjust the values according to your setup.