package broker

type Config struct {
	KafkaURL     string `json:"kafka_url"`
	Topic        string `json:"topic"`
	StartTimeout int    `json:"start_timeout"`
	StopTimeout  int    `json:"stop_timeout"`
}
