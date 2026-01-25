package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type TestEvent struct {
	Time    string
	Action  string
	Package string
	Test    string
	Elapsed float64
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: notifier <path-to-json-file>")
		os.Exit(1)
	}

	filePath := os.Args[1]
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var failedTests []string
	var passedCount int
	var failedCount int

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var event TestEvent

		if err := json.Unmarshal(scanner.Bytes(), &event); err != nil {
			continue
		}

		if event.Test != "" {
			if event.Action == "fail" {
				failedTests = append(failedTests, event.Test)
				failedCount++
			} else if event.Action == "pass" {
				passedCount++
			}
		}
	}

	sendNotification(passedCount, failedCount, failedTests)
}

func sendNotification(passed, failed int, failures []string) {
	title := "Go Test Results"
	body := ""
	urgency := "normal"
	icon := "terminal"

	if failed > 0 {
		title = fmt.Sprintf("❌ FAILED: %d", failed)
		urgency = "critical"
		limit := 3
		if len(failures) < limit {
			limit = len(failures)
		}
		body = fmt.Sprintf("Passed: %d\nFailed:\n- %s", passed, strings.Join(failures[:limit], "\n- "))
		if len(failures) > 3 {
			body += "\n...and more"
		}
	} else {
		title = "✅ SUCCESS"
		body = fmt.Sprintf("All %d tests passed!", passed)
	}

	cmd := exec.Command("notify-send", title, body, "-u", urgency, "-i", icon, "-t", "5000")
	if err := cmd.Run(); err != nil {
		fmt.Println("Error sending notification:", err)
	}
}
