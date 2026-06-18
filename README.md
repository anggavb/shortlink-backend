# ShortLink Backend

ShortLink Backend is a REST API for creating, managing, and resolving short URLs. It provides authentication, user-specific link management, Redis-backed link caching, click tracking, and Swagger API documentation.

## Technology Stack

![Go](https://img.shields.io/badge/Go-1.26.3-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Gin](https://img.shields.io/badge/Gin-1.12.0-008ECF?style=for-the-badge&logo=gin&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-Database-4169E1?style=for-the-badge&logo=postgresql&logoColor=white)
![Redis](https://img.shields.io/badge/Redis-Cache-DC382D?style=for-the-badge&logo=redis&logoColor=white)
![JWT](https://img.shields.io/badge/JWT-Auth-000000?style=for-the-badge&logo=jsonwebtokens&logoColor=white)
![Swagger](https://img.shields.io/badge/Swagger-OpenAPI-85EA2D?style=for-the-badge&logo=swagger&logoColor=black)
![MIT License](https://img.shields.io/badge/License-MIT-yellow?style=for-the-badge)

## Features

- User registration, login, and logout.
- JWT-based authentication for protected endpoints.
- Create short links from long URLs.
- Optional custom slugs or automatic random slug generation.
- List authenticated user's links with pagination and search.
- Soft delete links.
- Redirect short links to their original URLs.
- Track link clicks with IP address, user agent, and timestamp.
- Cache link resolution in Redis for faster redirects.
- Swagger API documentation at `/swagger/index.html`.

## Prerequisites

Before running this project locally, make sure you have:

- Git
- Go 1.26.3
- PostgreSQL
- Redis
- [golang-migrate](https://github.com/golang-migrate/migrate)
- `psql` command-line client for running seeders
- Optional: [Air](https://github.com/air-verse/air) for live reload
- Optional: [Swag CLI](https://github.com/swaggo/swag) for regenerating Swagger docs

## Setup Instruction

Clone the repository:

```bash
git clone https://github.com/anggavb/shortlink-backend.git
cd shortlink-backend
```

Copy the environment file and update the values:

```bash
cp .env.example .env
```

Configure these required values in `.env`:

```env
APP_HOST=localhost
APP_PORT=8080
DB_URL=postgresql://user:password@localhost:5432/database_name?sslmode=disable
RDB_HOST=localhost
RDB_PORT=6379
RDB_USER=
RDB_PASS=
RDB_PREFIX=shortlink
JWT_ISSUER=shortlink
JWT_SECRET=your-secret-key
```

Install dependencies:

```bash
go mod download
```

Run database migrations:

```bash
make migrate-up
```

Run seeders if needed:

```bash
make seed
```

Start the application:

```bash
go run ./cmd/main.go
```

The API will run at:

```text
http://localhost:8080
```

Init Swagger documentation:

```bash
swag init -g cmd/main.go -o docs
```

Open Swagger documentation at:

```text
http://localhost:8080/swagger/index.html
```

For live reload development, install Air and run:

Here for [.air.toml](https://gist.github.com/anggavb/11cc49709e1de97e1609f049bbf0569a) configuration file for Air live reload.

```bash
air
```

## Project Structure

```text
.
├── cmd/                 # Application entrypoint
├── db/                  # Database migrations and seeders
├── docs/                # Generated Swagger documentation
├── internal/            # Application modules
│   ├── binder/          # Request binding and validation helpers
│   ├── config/          # PostgreSQL and Redis connections
│   ├── controller/      # HTTP handlers
│   ├── dto/             # Request and response DTOs
│   ├── jwttoken/        # JWT helper utilities
│   ├── middleware/      # CORS and authentication middleware
│   ├── model/           # Domain models
│   ├── repository/      # PostgreSQL and Redis repositories
│   ├── response/        # Response helpers
│   ├── router/          # Route registration
│   └── service/         # Business logic
├── pkg/                 # Shared packages
├── public/              # Static assets
├── Makefile             # Migration and seeder commands
└── README.md
```

## How to Contribute

1. Fork this repository.
2. Create a new branch:

   ```bash
   git checkout -b feature/your-feature-name
   ```

3. Make your changes.
4. Format and test the project:

   ```bash
   gofmt -w .
   go test ./...
   ```

5. Commit your changes:

   ```bash
   git commit -m "feat: add your feature"
   ```

6. Push your branch and open a pull request.

## Related Project

- [ShortLink Frontend](https://github.com/anggavb/shortlink-frontend)

## License

This project is licensed under the [MIT License](LICENSE).
