package models

type Settings struct {
	DbPath string `json:"dbpath"`
	Port   string `json:"port"`
}

type Folder struct {
	Title string `json:"title"`
	Path  string `json:"path"`
	Type  int    `json:"type"`
	ID    string `json:"id"`
}

type CdnFileEntry struct {
	Name     string `json:"name"`
	FileType string `json:"type"`
	IsDir    bool   `json:"isDir"`
}

type User struct {
	ID           uint   `json:"id"`
	Username     string `json:"user"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
	Type         uint   `json:"type"`
}
