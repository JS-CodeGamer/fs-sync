// internal/database/sqlite.go
package database

import (
	"database/sql"
	"log"
	"path/filepath"
	"sync"
	"time"

	"github.com/js-codegamer/fs-sync/config"

	_ "github.com/mattn/go-sqlite3"
)

var (
	db   *sql.DB
	once sync.Once
)

func InitDatabase() *sql.DB {
	conf := config.GetConfig()
	var err error
	once.Do(func() {
		dbPath := filepath.Join(conf.Storage.BasePath, conf.Storage.DbPath)
		db, err = sql.Open("sqlite3", dbPath)
		if err != nil {
			log.Fatalf("Failed to open database: %v", err)
		}

		db.SetMaxOpenConns(conf.Database.MaxConnections)
		db.SetMaxIdleConns(conf.Database.MaxIdleConnection)
		db.SetConnMaxLifetime(5 * time.Minute)

		if err = db.Ping(); err != nil {
			log.Fatalf("Failed to ping database: %v", err)
		}

		createTables(db)
		createTriggers(db)
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
            id TEXT PRIMARY KEY,
            username TEXT UNIQUE NOT NULL,
            email TEXT UNIQUE NOT NULL,
            password TEXT NOT NULL,
            root_asset TEXT,
            created_at DATETIME DEFAULT (strftime('%FT%TZ', 'now')),
            FOREIGN KEY(root_asset) REFERENCES assets(id) DEFERRABLE INITIALLY DEFERRED
        )`,
		`CREATE TABLE IF NOT EXISTS assets (
            id TEXT PRIMARY KEY,
            owner TEXT NOT NULL,
            name TEXT NOT NULL,
            parent TEXT NOT NULL,
            is_public BOOLEAN DEFAULT 0,
            type TEXT NOT NULL,
            path TEXT NOT NULL,
            created_at DATETIME DEFAULT (strftime('%FT%TZ', 'now')),
            updated_at DATETIME DEFAULT (strftime('%FT%TZ', 'now')),
            FOREIGN KEY(owner) REFERENCES users(id) ON DELETE CASCADE,
            FOREIGN KEY(parent) REFERENCES assets(id) ON DELETE CASCADE,
            CHECK(type IN ('folder', 'file'))
        )`,
		`CREATE TABLE IF NOT EXISTS files (
            id TEXT PRIMARY KEY,
            size INTEGER NOT NULL,
            version INTEGER NOT NULL,
            path TEXT DEFAULT "", -- null path means not uploaded
            asset TEXT NOT NULL,
            mime_type TEXT DEFAULT "",
            FOREIGN KEY(asset) REFERENCES assets(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS permissions (
            id TEXT PRIMARY KEY,
            file_id INTEGER NOT NULL,
            user_id INTEGER NOT NULL,
            can_read BOOLEAN DEFAULT 0,
            can_write BOOLEAN DEFAULT 0,
            FOREIGN KEY(file_id) REFERENCES assets(id) ON DELETE CASCADE,
            FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
            UNIQUE (file_id, user_id)
        )`,
	}

	for _, table := range tables {
		_, err := db.Exec(table)
		if err != nil {
			log.Fatalf("Error creating table: %v", err)
		}
	}
}

func createTriggers(db *sql.DB) {
	triggers := []string{
		`CREATE TRIGGER IF NOT EXISTS updated_at__assets
		UPDATE ON assets
		FOR EACH ROW BEGIN
			UPDATE assets
			SET updated_at = strftime('%FT%TZ', 'now')
			WHERE id = old.id;
		END`,
	}

	for _, trigger := range triggers {
		_, err := db.Exec(trigger)
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
