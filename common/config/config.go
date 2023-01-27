package config

import "time"

type Config struct {
	Server Server `fig:"server" validate:"required"`
}

type Server struct {
	Port    int           `fig:"port" default:"8000"`
	Timeout time.Duration `fig:"timeout" default:"15s"`
}
