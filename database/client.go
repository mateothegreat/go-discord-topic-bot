package database

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mateothegreat/go-discord-topic-bot/prisma/db"
)

var DatabaseClient *db.PrismaClient

func Connect() {
	log.Printf("Asdf")
	DatabaseClient = db.NewClient(db.WithDatasourceURL(os.Getenv("DATABASE_URI")))

	if err := DatabaseClient.Prisma.Connect(); err != nil {
		panic(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-c
		Disconnect()
		os.Exit(0)
	}()
}

func Disconnect() {
	log.Println("disconnecting from database")
	if err := DatabaseClient.Prisma.Disconnect(); err != nil {
		panic(fmt.Errorf("could not disconnect: %w", err))
	}
}
