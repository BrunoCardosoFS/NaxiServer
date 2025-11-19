package database

import (
	"log"

	"github.com/BrunoCardosoFS/NaxiServer/models"
)

func GetCatalog() ([]models.Folder, error) {
	rows, err := DB.Query("SELECT id, title, path, type FROM catalog ORDER BY type ASC, title ASC")

	if err != nil {
		log.Printf("Error querying catalog: %v", err)
		return nil, err
	}
	defer rows.Close()

	var catalog []models.Folder
	for rows.Next() {
		var c models.Folder

		if err := rows.Scan(&c.ID, &c.Title, &c.Path, &c.Type); err != nil {
			log.Printf("Error scanning catalog line: %v", err)
			continue
		}

		catalog = append(catalog, c)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating through rows: %v", err)
		return nil, err
	}

	return catalog, nil
}

func AddFolderInCatalog(folder models.Folder) error {
	sql := `INSERT INTO catalog (id, title, path, type)
	        VALUES (?, ?, ?, ?)`

	_, err := DB.Exec(sql, folder.ID, folder.Title, folder.Path, folder.Type)

	if err != nil {
		log.Printf("Error while inserting folder: %v", err)
		return err
	}

	log.Printf("Folder inserted successfully: %s", folder.Title)
	return nil
}
