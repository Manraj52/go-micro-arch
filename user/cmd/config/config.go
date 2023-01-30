package config

import "time"

type Config struct {
	Server   Server `fig:"server" validate:"required"`
	Mongo    Mongo  `fig:"mongo" validate:"required"`
	Location string `fig:"location" validate:"required"`
}

type Server struct {
	Port      int           `fig:"port" default:"8000"`
	Timeout   time.Duration `fig:"timeout" default:"15s"`
	SecretKey string        `fig:"secretkey" validate:"required"`
}

type Mongo struct {
	Host     string        `fig:"host"`
	Port     int           `fig:"port"`
	Username string        `fig:"username" default:"user"`
	Password string        `fig:"password" default:"password"`
	Timeout  time.Duration `fig:"timeout" default:"10s"`
	Database string        `fig:"database" default:"MicroArch"`
}
