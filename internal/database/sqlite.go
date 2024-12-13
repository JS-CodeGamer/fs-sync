// internal/database/sqlite.go
package database

import (
	"database/sql"
	"log"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var (
	db   *sql.DB
	once sync.Once
)

func InitDatabase(dbPath string) *sql.DB {
	var err error
	once.Do(func() {
		db, err = sql.Open("sqlite3", dbPath)
		if err != nil {
			log.Fatalf("Failed to open database: %v", err)
		}

		db.SetMaxOpenConns(25)
		db.SetMaxIdleConns(25)
		db.SetConnMaxLifetime(5 * time.Minute)

		if err = db.Ping(); err != nil {
			log.Fatalf("Failed to ping database: %v", err)
		}

		createTables(db)
	})
	return db
}

func GetConnection() *sql.DB {
	if db == nil {
		log.Fatal("Database not initialized")
	}
	return db
}

func createTables(db *sql.DB) {
	tables := []string{
		`CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            username TEXT UNIQUE NOT NULL,
            email TEXT UNIQUE NOT NULL,
            password TEXT NOT NULL,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP
        )`,
		`CREATE TABLE IF NOT EXISTS files (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            user_id INTEGER NOT NULL,
            filename TEXT NOT NULL,
            path TEXT NOT NULL,
            size INTEGER NOT NULL,
            mime_type TEXT,
            is_public BOOLEAN DEFAULT 0,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY(user_id) REFERENCES users(id)
        )`,
		`CREATE TABLE IF NOT EXISTS file_permissions (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            file_id INTEGER NOT NULL,
            user_id INTEGER NOT NULL,
            can_read BOOLEAN DEFAULT 0,
            can_write BOOLEAN DEFAULT 0,
            FOREIGN KEY(file_id) REFERENCES files(id),
            FOREIGN KEY(user_id) REFERENCES users(id)
        )`,
	}

	for _, table := range tables {
		_, err := db.Exec(table)
		if err != nil {
			log.Fatalf("Error creating table: %v", err)
		}
	}
}

func CloseDatabase() {
	if db != nil {
		db.Close()
	}
}
