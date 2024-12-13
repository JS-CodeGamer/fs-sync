package files

import (
	"github.com/js-codegamer/fs-sync/internal/database"
	"github.com/js-codegamer/fs-sync/internal/models"
)

func CreateFile(file models.File) error {
	db := database.GetConnection()
	query := `INSERT INTO files 
        (user_id, filename, path, size, mime_type, is_public, created_at) 
        VALUES (?, ?, ?, ?, ?, ?, ?)`

	_, err := db.Exec(query,
		file.UserID,
		file.Filename,
		file.Path,
		file.Size,
		file.MimeType,
		file.IsPublic,
		file.CreatedAt,
	)
	return err
}

func GetUserFiles(userID int64) ([]models.File, error) {
	db := database.GetConnection()
	query := `SELECT id, filename, path, size, mime_type, is_public 
              FROM files WHERE user_id = ?`

	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []models.File
	for rows.Next() {
		var f models.File
		err := rows.Scan(&f.ID, &f.Filename, &f.Path, &f.Size, &f.MimeType, &f.IsPublic)
		if err != nil {
			return nil, err
		}
		files = append(files, f)
	}
	return files, nil
}

func DeleteFile(fileID int64, userID int64) error {
	db := database.GetConnection()
	query := `DELETE FROM files WHERE id = ? AND user_id = ?`

	_, err := db.Exec(query, fileID, userID)
	return err
}
