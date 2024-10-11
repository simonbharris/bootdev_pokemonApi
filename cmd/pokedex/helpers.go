package main

import (
	"fmt"
	"strings"
)

func printMainValue(key, value string) {
	fmt.Printf("%v: %v\n", key, value)
}

// if value is 0-length, don't print a `:`
func printListValues(key, value string) {
	if len(value) == 0 {
		fmt.Printf("  - %v\n", strings.ToLower(key))
	} else {
		fmt.Printf("  -%v: %v\n", strings.ToLower(key), value)
	}
}
