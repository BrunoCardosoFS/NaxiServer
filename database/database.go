package database

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB(dataSourceName string) {
	var err error
	DB, err = sql.Open("sqlite3", dataSourceName)

	if err != nil {
		log.Fatalf("Erro ao abrir banco de dados: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}

	log.Println("Banco de dados conectado com sucesso.")

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

	blocksTable := `
		CREATE TABLE IF NOT EXISTS blocks (
			id          INTEGER PRIMARY KEY,
			description TEXT
		);
	`

	musicalTable := `
		CREATE TABLE IF NOT EXISTS musical_items (
			id        INTEGER PRIMARY KEY AUTOINCREMENT,
			block_id  INTEGER REFERENCES blocks (id) ON DELETE CASCADE ON UPDATE CASCADE NOT NULL,
			file_path TEXT    NOT NULL,
			position  INTEGER NOT NULL
		);
	`

	commercialTable := `
		CREATE TABLE IF NOT EXISTS commercial_items (
			id        INTEGER PRIMARY KEY AUTOINCREMENT,
			block_id  INTEGER REFERENCES blocks (id) ON DELETE CASCADE ON UPDATE CASCADE NOT NULL,
			file_path TEXT    NOT NULL,
			position  INTEGER NOT NULL
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
	log.Println("Tabela 'catalog' verificada/criada.")

	execSQL(usersTable)
	log.Println("Tabela 'users' verificada/criada.")

	execSQL(blocksTable)
	log.Println("Tabela 'blocks' verificada/criada.")

	execSQL(musicalTable)
	log.Println("Tabela 'musical_items' verificada/criada.")

	execSQL(commercialTable)
	log.Println("Tabela 'commercial_items' verificada/criada.")

	execSQL(loggerTable)
	log.Println("Tabela 'logger' verificada/criada.")
}

func execSQL(sqlStatement string) {
	statement, err := DB.Prepare(sqlStatement)

	if err != nil {
		log.Fatal("Erro ao preparar statement de criação de tabela: ", err)
	}

	statement.Exec()
}
