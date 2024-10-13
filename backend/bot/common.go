package bot

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// readFile reads a file and creates an array of lines from file, delimiter is newline
func readFile(filePath string) []string {
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("[ERROR] Error encountered while reading script %s: %s", filePath, err.Error())
		return []string{}
	}
	// replace all to get rid of carriage return in windows os
	return strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n")
}

func cleanStringOfSpaces(input string) string {
	input = strings.ReplaceAll(input, "\t", "")
	input = strings.ReplaceAll(input, "\n", "")
	input = strings.TrimSpace(input)
	return input
}

func validateStrForClickEvent(input string) bool {
	if _, err := strconv.Atoi(input); err != nil {
		fmt.Printf("[ERROR] Invalid number: %s, error: %s", input, err.Error())
		return false
	}
	return true
}
