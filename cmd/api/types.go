package main

import "log"

type application struct {
	config config
	logger *log.Logger
}

type envelope map[string]interface{}
