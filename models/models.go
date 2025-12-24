package models

type Settings struct {
	DbPath string `json:"dbpath"`
	Port   string `json:"port"`
	JwtKey string `json:"token"`
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
	Name         string `json:"name"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
	Type         uint   `json:"type"`
}

type LoginResponse struct {
	Token       string   `json:"token"`
	Username    string   `json:"user"`
	Name        string   `json:"name"`
	Email       string   `json:"email"`
	Permissions []string `json:"permissions"`
}

type UserCredentialsRegister struct {
	Username    string   `json:"user"`
	Name        string   `json:"name"`
	Email       string   `json:"email"`
	Password    string   `json:"password"`
	Type        int      `json:"type"`
	Permissions []string `json:"permissions"`
}
