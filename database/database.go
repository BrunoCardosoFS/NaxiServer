package database

import (
	"database/sql"
	"log"
	"time"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB(dataSourceName string) {
	var err error
	DB, err = sql.Open("sqlite", dataSourceName)

	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	log.Println("Database successfully connected.")

	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(25)
	DB.SetConnMaxLifetime(5 * time.Minute)

	createTables()
}

func createTables() {
	catalogTable := `
		CREATE TABLE IF NOT EXISTS catalog (
			id    TEXT    UNIQUE NOT NULL,
			title TEXT    NOT NULL,
			path  TEXT    NOT NULL,
			type  INTEGER NOT NULL
		);
	`

	usersTable := `
		CREATE TABLE IF NOT EXISTS users (
			id            INTEGER PRIMARY KEY AUTOINCREMENT,
			user          TEXT    UNIQUE NOT NULL,
			name          TEXT    NOT NULL,
			email         TEXT    UNIQUE NOT NULL,
			password_hash TEXT    NOT NULL,
			type          INTEGER NOT NULL
		);
	`

	userPermissionsTable := `
		CREATE TABLE IF NOT EXISTS user_permissions (
			id         INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
			user       TEXT    REFERENCES users (user) ON DELETE CASCADE NOT NULL,
			permission TEXT    NOT NULL
		);
	`

	blocksTable := `
		CREATE TABLE IF NOT EXISTS blocks (
			id          INTEGER PRIMARY KEY,
			description TEXT
		);
	`

	musicalTable := `
		CREATE TABLE IF NOT EXISTS musical_items (
			id                    INTEGER PRIMARY KEY AUTOINCREMENT,
			block_id              INTEGER REFERENCES blocks (id) ON DELETE CASCADE ON UPDATE CASCADE NOT NULL,
			file_path             TEXT    NOT NULL,
			block_position        INTEGER NOT NULL,
			start_point           INTEGER,
			end_point             INTEGER,
			intro_point           INTEGER,
			mix_start_point       INTEGER,
			mix_end_point         INTEGER,
			pre_ramp_up_level     INTEGER,
			post_ramp_down_level  INTEGER,
			fade_in_end_point     INTEGER,
			ramp_up_start_point   INTEGER,
			ramp_up_end_point     INTEGER,
			ramp_down_start_point INTEGER,
			ramp_down_end_point   INTEGER,
			fade_out_start_point  INTEGER
		);
	`

	commercialTable := `
		CREATE TABLE IF NOT EXISTS commercial_items (
			id                    INTEGER PRIMARY KEY AUTOINCREMENT,
			block_id              INTEGER REFERENCES blocks (id) ON DELETE CASCADE ON UPDATE CASCADE NOT NULL,
			file_path             TEXT    NOT NULL,
			block_position        INTEGER NOT NULL,
			start_point           INTEGER,
			end_point             INTEGER,
			intro_point           INTEGER,
			mix_start_point       INTEGER,
			mix_end_point         INTEGER,
			pre_ramp_up_level     INTEGER,
			post_ramp_down_level  INTEGER,
			fade_in_end_point     INTEGER,
			ramp_up_start_point   INTEGER,
			ramp_up_end_point     INTEGER,
			ramp_down_start_point INTEGER,
			ramp_down_end_point   INTEGER,
			fade_out_start_point  INTEGER
		);
	`

	loggerTable := `
		CREATE TABLE IF NOT EXISTS logger (
			id             INTEGER PRIMARY KEY AUTOINCREMENT,
			block_id       INTEGER REFERENCES blocks (id) ON DELETE NO ACTION ON UPDATE NO ACTION NOT NULL,
			time           INTEGER NOT NULL,
			file_path      TEXT    NOT NULL,
			media_type     INTEGER NOT NULL,
			execution_type INTEGER NOT NULL
		);
	`

	execSQL(catalogTable)
	log.Println("Table 'catalog' verified/created.")

	execSQL(usersTable)
	log.Println("Table 'users' verified/created.")

	execSQL(userPermissionsTable)
	log.Println("Table 'user_permissions' verified/created.")

	execSQL(blocksTable)
	log.Println("Table 'blocks' verified/created.")

	execSQL(musicalTable)
	log.Println("Table 'musical_items' verified/created.")

	execSQL(commercialTable)
	log.Println("Table 'commercial_items' verified/created.")

	execSQL(loggerTable)
	log.Println("Table 'logger' verified/created.")
}

func execSQL(sqlStatement string) {
	statement, err := DB.Prepare(sqlStatement)

	if err != nil {
		log.Fatal("Error preparing table creation statement: ", err)
	}

	statement.Exec()
}
