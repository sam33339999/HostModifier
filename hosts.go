package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	commentPrefix = "# Host Modifier Generation."
)

func readHostsFile() ([]string, error) {
	hostsFilePath := getHostsFilePath()
	file, err := os.Open(hostsFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func writeHostsFile(lines []string) error {
	hostsFilePath := getHostsFilePath()
	file, err := os.OpenFile(hostsFilePath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}
	return writer.Flush()
}

func addHost(ip, hostname string) error {
	lines, err := readHostsFile()
	if err != nil {
		return err
	}

	newEntry := fmt.Sprintf("%s %s   %s (%s)", ip, hostname, commentPrefix, time.Now().Format("2006/01/02 15:04:05"))
	lines = append(lines, newEntry)

	return writeHostsFile(lines)
}

func modifyHost(hostname, newIP string) error {
	lines, err := readHostsFile()
	if err != nil {
		return err
	}

	for i, line := range lines {
		if strings.Contains(line, hostname) {
			lines[i] = fmt.Sprintf("%s %s   %s (%s)", newIP, hostname, commentPrefix, time.Now().Format("2006/01/02 15:04:05"))
			break
		}
	}

	return writeHostsFile(lines)
}

func confirmHostIP(hostname string) (string, error) {
	lines, err := readHostsFile()
	if err != nil {
		return "", err
	}

	for _, line := range lines {
		if strings.Contains(line, hostname) {
			parts := strings.Fields(line)
			if len(parts) > 0 {
				return parts[0], nil
			}
		}
	}

	return "", fmt.Errorf("host not found")
}

func deleteHost(hostname string) error {
	lines, err := readHostsFile()
	if err != nil {
		return err
	}

	for i, line := range lines {
		if strings.Contains(line, hostname) {
			lines[i] = fmt.Sprintf("# DELETE AT(%s) %s", time.Now().Format("2006/01/02 15:04:05"), line)
			break
		}
	}

	return writeHostsFile(lines)
}
