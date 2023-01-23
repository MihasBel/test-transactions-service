package kafka

type Config struct {
	StartTimeout int `json:"start_timeout"`
	StopTimeout  int `json:"stop_timeout"`
}
