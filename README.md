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
├── config/                   # Configuration files and utilities
├── internal/                 # Core application logic
│   ├── auth/                 # Authentication logic
│   ├── assets/               # Asset (file/folder) management logic
│   ├── database/             # Database management logic
│   ├── models/               # Internal modeling logic
│   └── routes/               # API routing logic
├── pkg/                      # Shared utility packages
│   ├── jwt/
│   ├── logger/
│   ├── password/
│   └── validator/
├── web/                      # Static files and frontend assets
│   ├── index.html            ## CURRENTLY UNUSED ( FOR COMPILED FRONTEND )
│   └── app.js
├── scripts/                  # Development / Management Scripts I Use
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
- **Availability**
  - `GET /ping`
- **Authentication**:
  - `POST /register`
  - `POST /login`
  - **Profile Management**
    - `GET /me`
    - `POST /me`
    - `DELETE /me`
- **File Management**:
  - `POST /a`                   -- Create file object / (upload) metadata
  - `GET /a/:fileID`            -- Get file (only works if you have permissions)
  - `PUT /a/:fileID`            -- Update file object / metadata
  - `PATCH /a/:fileID`          -- Upload file binary
  - `DELETE /files/:id`

## License
This project is licensed under the MIT License.
