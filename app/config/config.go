package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type appConfig struct {
	Properties map[string]string
	Port       string
}

var AppConfig *appConfig

func LoadConfig() {
	AppConfig = new(appConfig)
	AppConfig.Port = "1234"
	AppConfig.Properties = loadWebsites()
}

func loadWebsites() map[string]string {
	websites := make(map[string]string)

	file, err := os.Open("websites.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return websites
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		parts := strings.Split(line, "|")

		if len(parts) >= 2 {
			first := strings.TrimSpace(parts[0])
			second := strings.TrimSpace(parts[1])
			websites[first] = second
			fmt.Printf("website: %s, property ID: %s\n", first, second)
		} else {
			fmt.Println("Line doesn't have enough parts:", line)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}
	return websites
}
