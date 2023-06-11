package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"strings"
)

type Task struct {
	Name     string   `json:"name"`
	Command  string   `json:"command"`
	Requires []string `json:"requires"`
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/process_tasks", ProcessTasks).Methods("POST")
	fmt.Println("Server at 4000")
	log.Fatal(http.ListenAndServe(":4000", router))
}

func CheckErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func ProcessTasks(rw http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	recBody, err := io.ReadAll(req.Body)
	CheckErr(err)

	// Parse the request JSON
	var taskList []Task
	err = json.Unmarshal(recBody, &taskList)
	CheckErr(err)

	sortedTasks, err := SortTasks(taskList)
	CheckErr(err)

	bashScript := GenerateBashScript(sortedTasks)

	// Set the response content type
	rw.Header().Set("Content-Type", "text/plain")

	// Write the bash script as the response body
	_, err = rw.Write([]byte(bashScript))
	CheckErr(err)

	rw.WriteHeader(http.StatusOK)
}

func SortTasks(tasks []Task) ([]Task, error) {
	tasksMap := make(map[string]Task)
	for _, task := range tasks {
		tasksMap[task.Name] = task
	}
	var sortedTasks []Task
	executed := make(map[string]bool)

	for _, task := range tasks {
		err := ExecuteTasks(task, tasksMap, executed, &sortedTasks)
		CheckErr(err)
	}
	return sortedTasks, nil
}

func ExecuteTasks(task Task, taskMap map[string]Task, executed map[string]bool, sortedTasks *[]Task) error {
	if executed[task.Name] {
		return nil
	}

	for _, requiredTask := range task.Requires {
		taskName, success := taskMap[requiredTask]
		if success != true {
			return fmt.Errorf("required task does not exist: %s", requiredTask)
		}
		err := ExecuteTasks(taskName, taskMap, executed, sortedTasks)
		CheckErr(err)
	}

	executed[task.Name] = true
	*sortedTasks = append(*sortedTasks, task)
	return nil
}

func GenerateBashScript(tasks []Task) string {
	var sb strings.Builder
	sb.WriteString("#!/usr/bin/env bash\n")

	for _, task := range tasks {
		sb.WriteString(task.Command)
		sb.WriteString("\n")
	}
	return sb.String()
}
