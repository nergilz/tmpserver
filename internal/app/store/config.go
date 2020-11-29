package store

// DBConfig ...
type DBConfig struct {
	DatabaseURL string
}

// NewDBConfig ...
func NewDBConfig() *DBConfig {
	return &DBConfig{
		DatabaseURL: "host=localhost user=postgres dbname=restapi_tmp sslmode=disable",
	}
}
