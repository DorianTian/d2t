# Go Backend Service

A Go-based backend service for the DeepSeek2Web project, providing API endpoints and core functionality.

## Architecture

The service follows a clean architecture pattern with the following structure:

```
be-server/
├── internal/            # Internal application code
│   ├── api/             # API server implementation
│   ├── config/          # Configuration management
│   ├── middleware/      # HTTP middleware components
│   ├── models/          # Data models
│   ├── routes/          # API route definitions
│   └── services/        # Business logic services
├── core/                # Core domain logic
├── utils/               # Utility functions
├── main.go              # Application entry point
├── Dockerfile           # Container definition
├── go.mod               # Go module definition
└── .env                 # Environment configuration
```

### Application Flow

1. The application starts in `main.go`, which loads configuration using the config package
2. The API server is initialized in the `api` package, setting up the Gin router
3. Middleware components are registered to handle cross-cutting concerns like authentication and logging
4. Routes are registered to define API endpoints and connect them to service handlers
5. Services implement the business logic and interact with models
6. Models represent the data structures and interact with the PostgreSQL database

The architecture follows dependency injection principles, with each layer only depending on the layers below it:

```
Routes → Services → Models → Database
```

## Dependencies

The service relies on the following main dependencies:

- [Gin](https://github.com/gin-gonic/gin) (v1.10.0) - HTTP web framework
- [godotenv](https://github.com/joho/godotenv) (v1.5.1) - Environment variable loading from .env files
- [lib/pq](https://github.com/lib/pq) (v1.10.9) - PostgreSQL driver for database connectivity

## Environment Variables

Configure the application using the following environment variables (preferably in a `.env` file):

| Variable | Description | Default |
|----------|-------------|---------|
| PORT | Server port | 8080 |
| DB_HOST | Database host | localhost |
| DB_PORT | Database port | 5432 |
| DB_USER | Database username | - |
| DB_PASSWORD | Database password | - |
| DB_NAME | Database name | - |
| DB_SSLMODE | Database SSL mode | disable |
| API_TIMEOUT_SECONDS | API timeout in seconds | 300 |

## Getting Started

### Prerequisites

- Go 1.23.4 or higher
- PostgreSQL database

### Database Setup

1. Install PostgreSQL if not already installed
2. Create a new database:
   ```sql
   CREATE DATABASE d2t_db;
   ```
3. Create a user with appropriate permissions:
   ```sql
   CREATE USER youruser WITH PASSWORD 'yourpassword';
   GRANT ALL PRIVILEGES ON DATABASE d2t_db TO youruser;
   ```

### Local Development

1. Clone the repository:
   ```bash
   git clone https://github.com/DorianTian/d2t.git
   cd d2t/be-server
   ```

2. Create a `.env` file with your configuration:
   ```
   PORT=8080
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=youruser
   DB_PASSWORD=yourpassword
   DB_NAME=d2t_db
   DB_SSLMODE=disable
   API_TIMEOUT_SECONDS=300
   ```

3. Install dependencies:
   ```bash
   go mod download
   ```

4. Run database migrations:
   ```bash
   ./migrate.sh
   ```

5. Run the server:
   ```bash
   go run main.go
   ```
   
   Alternatively, you can specify a custom .env file location:
   ```bash
   go run main.go -env=/path/to/.env
   ```

6. The server should now be running at http://localhost:8080

### Running with Docker

1. Build the Docker image:
   ```bash
   docker build -t deepseek2web-backend .
   ```

2. Run the container:
   ```bash
   docker run -p 8080:8080 --env-file .env deepseek2web-backend
   ```

3. For a complete setup with PostgreSQL, use Docker Compose (create a docker-compose.yml file):
   ```yaml
   version: '3'
   services:
     app:
       build: .
       ports:
         - "8080:8080"
       env_file:
         - .env
       depends_on:
         - db
     db:
       image: postgres:15
       environment:
         POSTGRES_USER: ${DB_USER}
         POSTGRES_PASSWORD: ${DB_PASSWORD}
         POSTGRES_DB: ${DB_NAME}
       ports:
         - "${DB_PORT}:5432"
       volumes:
         - postgres_data:/var/lib/postgresql/data
   volumes:
     postgres_data:
   ```

   Then run:
   ```bash
   docker-compose up -d
   ```

## API Endpoints

The service provides RESTful API endpoints with the following base URL format:

```
http://localhost:8080/api/v1/...
```

Common endpoints include:

- `GET /health` - Service health check
- `GET /api/v1/...` - API endpoints (see API documentation for details)

## Database Migration

To run database migrations:

```bash
./migrate.sh
```

This script will apply migrations to your database based on the configuration in your `.env` file.

## Troubleshooting

- **Database Connection Issues**: Ensure your PostgreSQL server is running and the credentials in `.env` are correct
- **Port Already in Use**: Change the `PORT` in `.env` if port 8080 is already used by another service
- **Permission Denied**: Ensure `migrate.sh` has execute permissions (`chmod +x migrate.sh`)

## License

MIT License


