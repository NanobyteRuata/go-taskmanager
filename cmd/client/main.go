package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/NanobyteRuata/go-taskmanager/internal/models"
	"github.com/joho/godotenv"
)

var (
	hostUrl = "http://localhost:8080"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading it. Using default values.")
	}

	host := os.Getenv("HOST")
	if host != "" {
		hostUrl = host
	}

	command := os.Args[1]

	switch command {
	case "list":
		listTasks()
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Error: Task description required")
			printUsage()
			os.Exit(1)
		}
		title := os.Args[2]
		var dueDate string
		if len(os.Args) > 3 {
			dueDate = os.Args[3]
		}
		addTask(title, dueDate)
	case "complete":
		if len(os.Args) < 3 {
			fmt.Println("Error: Task ID required")
			printUsage()
			os.Exit(1)
		}
		completeTask(os.Args[2])
	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Error: Task ID required")
			printUsage()
			os.Exit(1)
		}
		deleteTask(os.Args[2])
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  client list                        - List all tasks")
	fmt.Println("  client add \"Task description\" [due-date] - Add a new task")
	fmt.Println("  client complete [task-id]          - Mark a task as completed")
	fmt.Println("  client delete [task-id]            - Delete a task")
}

func listTasks() {
	resp, err := http.Get(hostUrl + "/tasks")
	if err != nil {
		fmt.Printf("Error connecting to server: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Server returned error: %s\n", resp.Status)
		os.Exit(1)
	}

	var tasks []*models.Task
	if err := json.NewDecoder(resp.Body).Decode(&tasks); err != nil {
		fmt.Printf("Error parsing response: %v\n", err)
		os.Exit(1)
	}

	if len(tasks) == 0 {
		fmt.Println("No tasks found.")
		return
	}

	fmt.Println("Tasks:")
	for _, task := range tasks {
		status := " "
		if task.Completed {
			status = "✓"
		}
		dueStr := ""
		if !task.DueDate.IsZero() {
			dueStr = fmt.Sprintf(" (Due: %s)", task.DueDate.Format("2006-01-02"))
		}
		fmt.Printf("[%s] %s: %s%s\n", status, task.ID, task.Title, dueStr)
	}
}

func addTask(title, dueDate string) {
	payload := map[string]interface{}{
		"title": title,
	}

	if dueDate != "" {
		payload["due_date"] = dueDate
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		os.Exit(1)
	}

	resp, err := http.Post(
		hostUrl+"/tasks",
		"application/json",
		strings.NewReader(string(jsonData)),
	)
	if err != nil {
		fmt.Printf("Error connecting to server: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Server returned error: %s - %s\n", resp.Status, string(body))
		os.Exit(1)
	}

	var task *models.Task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		fmt.Printf("Error parsing response: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Task added with ID: %s\n", task.ID)
}

func completeTask(id string) {
	client := &http.Client{}
	req, err := http.NewRequest(
		http.MethodPatch,
		fmt.Sprintf("%s/tasks/%s/complete", hostUrl, id),
		nil,
	)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		os.Exit(1)
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error connecting to server: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Server returned error: %s - %s\n", resp.Status, string(body))
		os.Exit(1)
	}

	fmt.Printf("Task %s marked as completed\n", id)
}

func deleteTask(id string) {
	client := &http.Client{}
	req, err := http.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("%s/tasks/%s", hostUrl, id),
		nil,
	)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		os.Exit(1)
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error connecting to server: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Server returned error: %s - %s\n", resp.Status, string(body))
		os.Exit(1)
	}

	fmt.Printf("Task %s deleted\n", id)
}
