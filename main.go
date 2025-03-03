package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
)

func main() {
	// Check if the application has permission to modify the /etc/hosts file
	if !checkPermissions() {
		fmt.Println("Error: The application does not have permission to modify the /etc/hosts file.")
		return
	}

	// Set up routes and start the server
	setupRoutes()
}

func checkPermissions() bool {
	hostsFilePath := getHostsFilePath()
	file, err := os.OpenFile(hostsFilePath, os.O_WRONLY, 0644)
	if err != nil {
		return false
	}
	defer file.Close()
	return true
}

func getHostsFilePath() string {
	if runtime.GOOS == "windows" {
		return "C:\\Windows\\System32\\drivers\\etc\\hosts"
	}
	return "/etc/hosts"
}

func setupRoutes() {
	// Call the function to set up routes and start the server
	// This function will be implemented in routes.go
}
