package config

type Config struct {
	Server Server `yaml:"Server"`
}

type Server struct {
	Port       string `yaml:"port"`
	WorkersNum int    `yaml:"workers_num"`
}
