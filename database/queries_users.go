package database

import (
	"log"

	"github.com/BrunoCardosoFS/NaxiServer/models"
)

func CreateUser(user, name, email, hashedPassword string, typeUser int) error {
	sql := `INSERT INTO users (user, name, email, password_hash, type) VALUES (?, ?, ?, ?, ?)`

	_, err := DB.Exec(sql, user, name, email, hashedPassword, typeUser)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return err
	}

	return nil
}

func DeleteUser(user string) error {
	sql := `DELETE FROM users WHERE user = ?`

	result, err := DB.Exec(sql, user)
	if err != nil {
		log.Printf("Error deleting user '%s': %v", user, err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		// return sql.ErrNoRows
	}

	log.Printf("User successfully deleted: %s", user)

	return nil
}

func GetUserByUsername(username string) (*models.User, error) {
	sql := `SELECT id, user, email, password_hash, type FROM users WHERE user = ?`

	row := DB.QueryRow(sql, username)

	var u models.User
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.Type)

	if err != nil {
		log.Printf("Error retrieving user: %v", err)
		return nil, err
	}

	return &u, nil
}
