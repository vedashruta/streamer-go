package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	Conf Server
)

type Server struct {
	Port        string
	Concurrency int
}

func Init() (err error) {
	err = godotenv.Load(".env")
	if err != nil {
		return
	}
	Conf.Port = os.Getenv("PORT")
	Conf.Concurrency, err = strconv.Atoi(os.Getenv("CONCURRENCY"))
	if err != nil {
		return
	}
	return
}
