package models

type DatabaseConfig struct {
	Driver     string   `json:"driver"`
	Connection string   `json:"connection"`
	Post       []string `json:"post"`
}

// StaticConfig are global configs from the start up file
type StaticConfig struct {
	ImportDir string         `json:"import"`
	Database  DatabaseConfig `json:"database"`
}
