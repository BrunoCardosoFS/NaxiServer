package database

import (
	"database/sql"
	"log"
	"strings"

	"github.com/BrunoCardosoFS/NaxiServer/models"
)

func HasAdmin() bool {
	var count int

	sql := `SELECT COUNT(*) FROM users WHERE type = 0`
	err := DB.QueryRow(sql).Scan(&count)
	if err != nil {
		return true
	}

	return count > 0
}

func CreateUser(user, name, email, hashedPassword string, typeUser int) error {
	sql := `INSERT INTO users (user, name, email, password_hash, type) VALUES (?, ?, ?, ?, ?)`

	_, err := DB.Exec(sql, user, name, email, hashedPassword, typeUser)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return err
	}

	return nil
}

func AddPermissions(user string, permissions []string) error {
	if len(permissions) == 0 {
		return nil
	}

	_, err := DB.Exec("DELETE FROM user_permissions WHERE user = ?", user)
	if err != nil {
		return err
	}

	var valueStrings []string
	var valueArgs []interface{}

	for _, perm := range permissions {
		valueStrings = append(valueStrings, "(?, ?)")

		valueArgs = append(valueArgs, user)
		valueArgs = append(valueArgs, perm)
	}

	sql := "INSERT INTO user_permissions (user, permission) VALUES " + strings.Join(valueStrings, ",")

	_, err = DB.Exec(sql, valueArgs...)

	return err
}

func HasPermission(username string, permission string) bool {
	var exists bool

	query := `
        SELECT EXISTS(
            SELECT 1 
            FROM users u
            LEFT JOIN user_permissions up ON up.user = u.user 
            WHERE u.user = ? 
            AND (u.type = 0 OR up.permission = ?)
        )
    `

	err := DB.QueryRow(query, username, permission).Scan(&exists)
	if err != nil {
		return false
	}

	return exists
}

func DeleteUser(user string) error {
	query := `DELETE FROM users WHERE user = ?`

	result, err := DB.Exec(query, user)
	if err != nil {
		log.Printf("Error deleting user '%s': %v", user, err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	log.Printf("User successfully deleted: %s", user)

	return nil
}

func GetUserByUsername(username string) (*models.User, error) {
	sql := `SELECT id, user, name, email, password_hash, type FROM users WHERE user = ?`

	row := DB.QueryRow(sql, username)

	var u models.User
	err := row.Scan(&u.ID, &u.Username, &u.Name, &u.Email, &u.PasswordHash, &u.Type)

	if err != nil {
		log.Printf("Error retrieving user: %v", err)
		return nil, err
	}

	return &u, nil
}

func GetPermissionsByUsername(username string, userType uint) ([]string, error) {
	var permissions []string

	if userType == 0 {
		permissions = []string{"ADMIN"}
		return permissions, nil
	}

	sql := `SELECT permission FROM user_permissions WHERE user = ?`
	rows, err := DB.Query(sql, username)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p string
		if err := rows.Scan(&p); err == nil {
			permissions = append(permissions, p)
		}
	}

	if permissions == nil {
		permissions = []string{}
	}

	return permissions, nil
}
