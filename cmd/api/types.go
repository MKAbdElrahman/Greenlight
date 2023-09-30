package main

import (
	"log"

	"greenlight.mkabdelrahman.net/internal/data"
)

type application struct {
	config config
	logger *log.Logger
	models data.Models
}

type envelope map[string]interface{}
