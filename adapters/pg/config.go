package pg

type Config struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	DBName   string `json:"dbname"`
	User     string `json:"user"`
	Password string `json:"password"`
	TimeZone string `json:"timezone"`
	SSLMode  string `json:"sslmode"`
	DebugSQL bool   `json:"debugsql"`
}
