package database

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/js-codegamer/fs-sync/internal/models"
	"github.com/js-codegamer/fs-sync/internal/utils"
	"github.com/js-codegamer/fs-sync/pkg/logger"
)

func CreateUser(user models.User, txn *sql.Tx) (models.User, error) {
	// do all work in txn so if any fails it reverts as a whole
	var (
		err error
	)
	noTxn := txn == nil
	if noTxn {
		txn, err = db.Begin()
		if err != nil {
			logger.Sugar.Errorln(err.Error())
			return models.User{}, err
		}
	}

	// create user
	user.ID = uuid.New().String()
	query := `
	INSERT INTO users
			(id, username, password, email)
		VALUES
			(?, ?, ?, ?)
	`
	if _, err := txn.Exec(query, user.ID, user.Username, user.Password, user.Email); err != nil {
		logger.Sugar.Errorln(err.Error())
		txn.Rollback()
		return models.User{}, err
	}

	// create root dir for user
	if err := utils.CreateRootDir(user); err != nil {
		logger.Sugar.Errorln(err.Error())
		txn.Rollback()
		return models.User{}, err
	}

	// create database entry for root dir
	root_dir := models.Asset{
		OwnerID:  user.ID,
		Name:     user.Username,
		Path:     utils.GetRootDir(user),
		ParentID: "-1",
		Type:     models.FolderType,
	}
	err = CreateAsset(&root_dir, txn)
	if err != nil {
		logger.Sugar.Errorln(err.Error())
		utils.DestroyRootDir(user)
		txn.Rollback()
		return models.User{}, err
	}

	// update user database entry to include root dir
	query = `
	UPDATE users
		SET root_asset = ?
		WHERE id = ?
	`
	if _, err := txn.Exec(query, root_dir.ID, user.ID); err != nil {
		logger.Sugar.Errorln(err.Error())
		utils.DestroyRootDir(user)
		txn.Rollback()
		return models.User{}, err
	}

	if err = txn.Commit(); err != nil {
		logger.Sugar.Errorln(err.Error())
		utils.DestroyRootDir(user)
		return models.User{}, err
	}

	return user, nil
}

func FindUserByUsername(username string) (models.User, error) {
	var user models.User

	query := `
	SELECT id, username, email, password, root_asset
		FROM users
		WHERE username = ?
	`
	err := db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.RootDirID,
	)

	if err == sql.ErrNoRows {
		logger.Sugar.Errorln(err.Error())
		return user, errors.New("user not found")
	} else if err != nil {
		logger.Sugar.Errorln(err.Error())
		return user, err
	}

	return user, nil
}

func FindUserByID(id string) (models.User, error) {
	var user models.User

	query := `
	SELECT id, username, email, password, root_asset
		FROM users
		WHERE id = ?
	`
	err := db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.RootDirID,
	)

	if err == sql.ErrNoRows {
		logger.Sugar.Errorln(err.Error())
		return user, errors.New("user not found")
	} else if err != nil {
		logger.Sugar.Errorln(err.Error())
		return user, err
	}

	return user, nil
}

func UpdateUser(user models.User, txn *sql.Tx) error {
	query := `
	UPDATE users
		SET email = ?, password = ?
		WHERE id = ?
	`

	var err error
	if txn == nil {
		_, err = db.Exec(query, user.Email, user.Password, user.ID)
	} else {
		_, err = txn.Exec(query, user.Email, user.Password, user.ID)
	}

	return err
}

func DeleteUser(user models.User, txn *sql.Tx) error {
	query := `
	DELETE FROM users
		WHERE id = ?
	`

	var err error
	if txn == nil {
		_, err = db.Exec(query, user.Email, user.Password, user.ID)
	} else {
		_, err = txn.Exec(query, user.Email, user.Password, user.ID)
	}

	return err
}
