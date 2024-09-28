package main

import (
	"bufio"
	"fmt"
	"os"
)

func beginCli() {
	scanner := bufio.NewScanner(os.Stdin)
	cliCommands := getCliCommands()
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		userInput := scanner.Text()
		cmd, found := cliCommands[userInput]
		if !found {
			fmt.Println("Invalid command. see 'help'")
			continue
		}

		err := cmd.callback()
		if err != nil {
			fmt.Printf("Error with command '%v': %v\n", cliCommands[userInput].name, err)
		}
	}
}
