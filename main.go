package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println("Hello, World!")

}

func cleanInput(text string) []string {
	lowerString := strings.ToLower(text)
	sliceString := strings.Fields(lowerString)
	return sliceString

}
