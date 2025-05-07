package main

import (
	"kanbanify-api/db"
	"kanbanify-api/handler"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetOutput(os.Stdout)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.InfoLevel)
	logrus.Info("Starting Kanbanify-api...")

	err := godotenv.Load()
	if err != nil {
		logrus.Fatal("Error loading .env file")
	}

	conn, err := db.Connect()
	if err != nil {
		logrus.Fatal("Error connecting to database << ", err)
		os.Exit(1)
	}
	defer conn.Close()
	logrus.Info("Connected to database")

	handler.Handler()
}
