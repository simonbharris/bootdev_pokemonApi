package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func beginCli() {
	scanner := bufio.NewScanner(os.Stdin)
	cliCommands := getCliCommands()
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		userInput := scanner.Text()
		words := strings.Split(userInput, " ")
		commandWord := words[0]
		cmd, found := cliCommands[commandWord]
		if !found {
			fmt.Println("Invalid command. see 'help'")
			continue
		}

		err := cmd.callback(words[1:]...)
		if err != nil {
			fmt.Printf("Error with command '%v': %v\n", cliCommands[commandWord].name, err)
		}
	}
}
