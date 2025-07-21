# DEPLOY.txt

## Project Overview

This project is a **bus ticketing API** built with Go (Gin framework).  
It handles buses, schedules, reservations, and tickets with a RESTful API structure.  
The system uses SQLite by default but can be adapted to PostgreSQL/MySQL for production.  
API documentation is generated with Swagger and served at `/swagger/index.html`.

---

## Requirements

- **Go**: >=1.20
- **OS**: Linux, macOS, or Windows 11 (tested on Windows 11)
- **Database**:
  - Default: SQLite

**Go Modules (from `go.mod`):**

- `github.com/gin-gonic/gin v1.10.1`
- `gorm.io/gorm v1.30.0`
- `gorm.io/driver/sqlite v1.6.0`
- `github.com/swaggo/gin-swagger v1.6.0`
- `github.com/swaggo/swag v1.16.5`
- `gorm.io/datatypes v1.2.6`
- `github.com/swaggo/files v1.0.1`

---

## How the System Runs

- **Configuration**
  - HTTP server runs on port **8000** by default.
  - SQLite database file is created in the project root as `ticketings.sqlite3`.
  - Swagger docs are available at `/swagger/index.html` when the server is running.

- **Environment**
  - Development:
    - Runs directly with `go run cmd/main.go`.
    - Hot reloads can be enabled using tools like [Air](https://github.com/cosmtrek/air).

- **Startup Behavior**
  - On startup, the app initializes the database and creates tables if they don’t exist.
  - The API routes are grouped under `/api/v1`.

---

## Additional Notes

- For API testing, import the included Bruno collection in the `api_collection` folder (`Bruno is just an lightweight HTTP client just like postman`).
- Recommended for development/testing, not high-traffic production environments.

> [!NOTE]
>
> There as a make file that includes the neccessary command to run and can be used for swagger `docs`, `run`, `live` run(`Air`) and to `get` the dependencies
