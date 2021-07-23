package main

import (
	"log"

	"github.com/IfuryI/ChatAPI/internal/server"
	constants "github.com/IfuryI/ChatAPI/pkg/const"
)

func main() {
	app := server.NewApp()

	if err := app.Run(constants.Port); err != nil {
		log.Fatal(err)
	}
}
