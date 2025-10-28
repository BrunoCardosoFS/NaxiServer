package models

type Folder struct {
	Title string `json:"title"`
	Path  string `json:"path"`
	Type  int    `json:"type"`
	ID    string `json:"id"`
}
