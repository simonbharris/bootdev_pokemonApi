package main

import (
	"fmt"
	"log/slog"
	"os"
)

func main() {
	initLogger()

	fmt.Println("Welcome to Someone's pokedex!")
	beginCli()
}

func initLogger() {
	programLevel := new(slog.LevelVar)
	programLevel.Set(slog.LevelDebug)
	//programLevel.Set(slog.LevelInfo)
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: programLevel})
	logger := slog.New(handler)
	slog.SetDefault(logger)
}
