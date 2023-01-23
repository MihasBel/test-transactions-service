package app

// Config exported variable to contain config values
var Config Configuration

// Configuration exported type for config
type Configuration struct {
	Host                  string `json:"host"`
	Port                  int    `json:"port"`
	DBName                string `json:"dbname"`
	User                  string `json:"user"`
	Password              string `json:"password"`
	TimeZone              string `json:"timezone"`
	SSLMode               string `json:"sslmode"`
	DebugSQL              bool   `json:"debugsql"`
	StartTimeout          int    `json:"StartTimeout"`
	StopTimeout           int    `json:"StopTimeout"`
	Address               string `json:"Address"`
	ConsoleLoggingEnabled bool   `json:"ConsoleLoggingEnabled"`
	FileLoggingEnabled    bool   `json:"FileLoggingEnabled"`
	LogDirectory          string `json:"LogDirectory"`
	LogFilename           string `json:"LogFilename"`
	LogMaxBackups         int    `json:"LogMaxBackups"`
	LogMaxSize            int    `json:"LogMaxSize"`
	LogMaxAge             int    `json:"LogMaxAge"`
	KafkaURL              string `json:"KafkaURL"`
	PartitionsCount       int    `json:"PartitionsCount"`
}
