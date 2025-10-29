package database

import (
	"log"

	"github.com/BrunoCardosoFS/NaxiServer/models"
)

func GetCatalog() ([]models.Folder, error) {
	rows, err := DB.Query("SELECT id, title, path, type FROM catalog ORDER BY type ASC, title ASC")

	if err != nil {
		log.Printf("Erro ao consultar catalogo: %v", err)
		return nil, err
	}
	defer rows.Close()

	var catalog []models.Folder
	for rows.Next() {
		var c models.Folder

		if err := rows.Scan(&c.ID, &c.Title, &c.Path, &c.Type); err != nil {
			log.Printf("Erro ao escanear linha do catalogo: %v", err)
			continue
		}

		catalog = append(catalog, c)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Erro ao iterar pelas linhas: %v", err)
		return nil, err
	}

	return catalog, nil
}

func AddFolderInCatalog(folder models.Folder) error {
	sql := `INSERT INTO catalog (id, title, path, type)
	        VALUES (?, ?, ?, ?)`

	_, err := DB.Exec(sql, folder.ID, folder.Title, folder.Path, folder.Type)

	if err != nil {
		log.Printf("Erro ao executar inserção de pasta: %v", err)
		return err
	}

	log.Printf("Pasta inserida com sucesso: %s", folder.Title)
	return nil
}

func CreateUser(user, name, email, hashedPassword string, typeUser int) error {
	sql := `INSERT INTO users (user, name, email, password_hash, type) VALUES (?, ?, ?, ?, ?)`

	_, err := DB.Exec(sql, user, name, email, hashedPassword, typeUser)
	if err != nil {
		log.Printf("Erro ao criar usuário: %v", err)
		return err
	}

	return nil
}

func DeleteUser(user string) error {
	sql := `DELETE FROM users WHERE user = ?`

	result, err := DB.Exec(sql, user)
	if err != nil {
		log.Printf("Erro ao deletar usuário '%s': %v", user, err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		// return sql.ErrNoRows
	}

	log.Printf("Usuário deletado com sucesso: %s", user)

	return nil
}

func GetUserByUsername(username string) (*models.User, error) {
	sql := `SELECT id, user, email, password_hash, type FROM users WHERE user = ?`

	row := DB.QueryRow(sql, username)

	var u models.User
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.Type)

	if err != nil {
		log.Printf("Erro ao pegar usuário: %v", err)
		return nil, err
	}

	return &u, nil
}
