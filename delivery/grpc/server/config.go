package server

type Config struct {
	Address      string `json:"address"`
	StartTimeout int    `json:"start_timeout"`
	StopTimeout  int    `json:"stop_timeout"`
}
