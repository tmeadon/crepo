package main

import (
	"fmt"
	"os"
)

var helpTemplate = `
Usage: cr <search_term>

Configuration can be found at %s
`

func printHelp() {
	configPath, err := getConfigFilePath()

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error getting config path:", err)
	}

	fmt.Println(fmt.Sprintf(helpTemplate, configPath))
}

func main() {
	config, err := loadConfig()

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error loading config:", err)
		os.Exit(1)
	}

	switch len(os.Args) {
	case 1:
		selectRepo("", config)
	case 2:
		if os.Args[1] == "-h" || os.Args[1] == "--help" {
			printHelp()
			return
		}

		if os.Args[1] == "-p" || os.Args[1] == "--print" {
			printShellFunction()
			return
		}

		selectRepo(os.Args[1], config)
	default:
		fmt.Fprintln(os.Stderr, "Multiple search terms not supported")
		printHelp()
		os.Exit(1)
	}
}
