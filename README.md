# Project: Google Drive Clone

## Description
A Google Drive-like application built using Golang for the backend with SQLite as the database. This project includes features such as user authentication, file management, and a web interface for interaction.

## Features
- User authentication using JWT.
- File upload, download, and management.
- SQLite database for storing user and file metadata.
- Extendable project structure for scalability.

## Project Structure
```
project-root/
├── cmd/
│   └── server/               # Entry point for the application
│       └── main.go
├── config/                   # Configuration files and utilities
│   └── config.go
├── internal/                 # Core application logic
│   ├── auth/                 # Authentication logic
│   │   ├── auth.go
│   │   ├── middleware.go
│   │   └── jwt.go
│   ├── file/                 # File management logic
│   │   ├── file.go
│   │   └── upload.go
│   └── user/                 # User management logic
│       ├── user.go
│       └── profile.go
├── pkg/                      # Shared utility packages
│   ├── database/
│   │   └── sqlite.go
│   └── logger/
│       └── logger.go
├── web/                      # Static files and frontend assets
│   ├── index.html
│   └── app.js
├── .env                      # Environment variables
├── go.mod                    # Go module file
├── go.sum                    # Dependencies checksum
└── README.md                 # Documentation
```

## Getting Started

### Prerequisites
- Go 1.20+
- SQLite 3

### Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/your-repo/google-drive-clone.git
   cd google-drive-clone
   ```
2. Install dependencies:
   ```bash
   go mod tidy
   ```
3. Set up environment variables in `.env` file.

### Running the Application
```bash
go run cmd/server/main.go
```

## API Endpoints
- **Authentication**:
  - `POST /login`
  - `POST /register`
- **File Management**:
  - `POST /upload`
  - `GET /files`
  - `DELETE /files/:id`

## License
This project is licensed under the MIT License.
