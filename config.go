package main

import (
	"log"
	"os"
)

type Config struct {
	DB struct {
		Username string
		Password string
		Host     string
		Port     string
		Database string
	}
}

func ParseConfig() Config {
	var config Config

	if val, exists := os.LookupEnv("DB_USERNAME"); exists {
		config.DB.Username = val
	} else {
		log.Fatal("DB_USERNAME is missing")
	}

	if val, exists := os.LookupEnv("DB_PASSWORD"); exists {
		config.DB.Password = val
	} else {
		log.Fatal("DB_PASSWORD is missing")
	}

	if val, exists := os.LookupEnv("DB_HOST"); exists {
		config.DB.Host = val
	} else {
		log.Fatal("DB_HOST is missing")
	}

	if val, exists := os.LookupEnv("DB_PORT"); exists {
		config.DB.Port = val
	} else {
		log.Fatal("DB_PORT is missing")
	}

	if val, exists := os.LookupEnv("DB_DATABASE"); exists {
		config.DB.Database = val
	} else {
		log.Fatal("DB_DATABASE is missing")
	}

	return config
}
