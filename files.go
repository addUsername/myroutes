package main

import (
	"fmt"
	"os"
)

func readRoutesJson(path string) string {

	fmt.Println("reading " + path)
	lines, err := os.ReadFile(path)

	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(lines)
}
func saveJson(path string, json string) bool {

	file, err := os.Create(path)
	if err != nil {
		return false
	}
	defer file.Close()
	file.WriteString(json)
	return true
}
