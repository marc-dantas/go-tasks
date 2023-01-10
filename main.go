package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"bufio"
	"strings"
)

type Todo struct {
	Id    int
	Title string
	State bool
}

func gets(prompt string) string {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print(prompt)
	var line string
	if scanner.Scan() {
		line = scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return line
}

func geti(prompt string) int {
	v, e := strconv.Atoi(gets(prompt))
	if e != nil {
		log.Fatal("Invalid input.")
	}
	return v
}

func show_options(opts []string) {
	for c, x := range opts {
		fmt.Println("=>  [", c+1, "]", x)
	}
}

func apply_file(path string, todos []Todo) {
	var jsons []string
	for _, todo := range todos {
		bs, e := json.Marshal(todo)
		if e != nil {
			log.Fatal(e)
		}
		jsons = append(jsons, string(bs))
	}
	writelines(path, jsons)
}

func writelines(filepath string, lines []string) {
	f, err := os.Create(filepath)
	if err != nil {
		log.Fatal(err)
	}
	for i, content := range lines {
		if i == len(lines)-1 {
			_, err = f.WriteString(content)
		} else {
			_, err = f.WriteString(content + "\n")
		}

		if err != nil {
			log.Fatal(err)
		}
		f.Sync()
	}
}

func readlines(filepath string) []string {
	buf, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Fatal(err)
	}
	result := string(buf)
	return strings.Split(result, "\n")
}

// I know this is slow, but wat can I do?
func remove_todo(seq *[]Todo, idx int) {
	var new_seq []Todo
	for i, x := range *seq {
		if i != idx {
			new_seq = append(new_seq, x)
		}
	}
	*seq = new_seq
}

func handle_option(opt int, todos *[]Todo) {
	if opt == 1 {
		id := len(*todos)
		title := gets("Title of the task: ")
		*todos = append(*todos, Todo{
			Id:    id,
			Title: title,
			State: false,
		})
	} else if opt == 2 {
		id := geti("Task ID: ")
		for i := range *todos {
			if (*todos)[i].Id == id {
				(*todos)[i].State = !(*todos)[i].State
			}
		}
	} else if opt == 3 {
		var idx int
		id := geti("Task ID: ")
		for i := range *todos {
			if i == id {
				break
			}
			idx++
		}
		remove_todo(todos, idx)
	}
	apply_file("data", *todos)
}

func apply_todos(filepath string, todos *[]Todo) {
	lines := readlines(filepath)
	var new_todos []Todo
	for _, line := range lines {
		var result Todo
		if line == "" {
			continue
		}
		e := json.Unmarshal([]byte(line), &result)
		if e != nil {
			log.Fatal(e)
		}
		new_todos = append(new_todos, result)
	}
	*todos = new_todos
}

func show_todos(todos []Todo) {
	if len(todos) == 0 {
		fmt.Println("\nYou don't have tasks!")
	} else {
		fmt.Println("\nYour tasks:")
		for _, x := range todos {
			if x.State {
				fmt.Printf("%d [X] %s\n", x.Id, x.Title)
			} else {
				fmt.Printf("%d [ ] %s\n", x.Id, x.Title)
			}
		}
	}
	fmt.Println()
}
func main() {
	fmt.Println("Welcome to GoTasks. A simple to do list written in Go")
	fmt.Println("Choose your option below. Type ^C (Control + c) to exit")
	var todos []Todo
	var opt int
	apply_todos("data", &todos)
	for {
		show_options([]string{
			"Add a new task",
			"Toggle task's state (checked or not)",
			"Remove task",
		})
		show_todos(todos)
		opt = geti("your option: ")
		handle_option(opt, &todos)
		fmt.Printf("\x1bc")
	}
}
