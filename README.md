# Asset Dairy Go API

This project is a Go-based RESTful API for managing assets, users, and related operations. It is designed to be modular, secure, and easy to deploy.

## Features
- User authentication (JWT-based)
- Asset management
- Database migrations
- Dockerized development environment
- Configurable via environment variables

## Project Structure
```
server/
├── db/                # Database connection and logic
├── handlers/          # HTTP handlers (controllers)
├── models/            # Data models
├── migrations/        # SQL migration files
├── main.go            # Entry point
├── docker-compose.yml # Docker configuration
├── .env.example       # Example environment variables
```

## Getting Started

### Prerequisites
- Go 1.18+
- Docker & Docker Compose (optional, for containerized setup)
- PostgreSQL (or your chosen DB)

### Setup
1. Clone the repository:
   ```bash
   git clone <repo-url>
   cd asset-dairy/server
   ```
2. Copy and configure environment variables:
   ```bash
   cp .env.example .env
   # Edit .env as needed
   ```
3. Run database migrations (if applicable):
   ```bash
   # Using migrate tool or your preferred method
   ```
4. Start the API:
   ```bash
   go run main.go
   # or with Docker
   docker-compose up --build
   ```

## API Endpoints
- `/api/auth/login` - User login
- `/api/auth/register` - User registration
- `/api/assets` - CRUD for assets

## Development
- Code is organized by feature (handlers, models, db)
- Use Go modules for dependency management
- Lint and test before pushing changes

## License
MIT

## Author
- [Your Name]

---
Feel free to customize this README to better fit your exact project details and team!
