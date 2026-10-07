package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/gearmobile/todo-manager/internal/storage"
	"github.com/gearmobile/todo-manager/internal/todo"
)

const tasksFile = "tasks.json"

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "add":
		cmdAdd(os.Args[2:])
	case "list":
		cmdList(os.Args[2:])
	case "complete":
		cmdComplete(os.Args[2:])
	case "delete":
		cmdDelete(os.Args[2:])
	case "export":
		cmdExport(os.Args[2:])
	case "load":
		cmdLoad(os.Args[2:])
	default:
		fmt.Printf("неизвестная команда %s\n", os.Args[1])
		usage()
		os.Exit(1)
	}
}

func cmdAdd(args []string) {
	fs := flag.NewFlagSet("add", flag.ExitOnError)
	desc := fs.String("desc", "", "task description")
	fs.Parse(args)

	if *desc == "" {
		fmt.Println("error: --desc is required")
		os.Exit(1)
	}

	tasks, err := storage.LoadJSON(tasksFile)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	mgr := todo.NewManager(tasks)
	task := mgr.Add(*desc)

	if err := storage.SaveJSON(tasksFile, mgr.Tasks()); err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	fmt.Printf("добавлена задача #%d: %s\n", task.ID, task.Description)
}

func cmdList(args []string) {
	fs := flag.NewFlagSet("list", flag.ExitOnError)
	filter := fs.String("filter", "all", "filter: all, done, pending")
	fs.Parse(args)

	tasks, err := storage.LoadJSON(tasksFile)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	mgr := todo.NewManager(tasks)
	list := mgr.List(*filter)

	if len(list) == 0 {
		fmt.Println("no tasks")
		return
	}
	for _, t := range list {
		status := "pending"
		if t.Done {
			status = "done"
		}
		fmt.Printf("[%d] %s — %s\n", t.ID, status, t.Description)
	}
}

func cmdComplete(args []string) {
	fs := flag.NewFlagSet("complete", flag.ExitOnError)
	id := fs.Int("id", 0, "task id")
	fs.Parse(args)

	if *id <= 0 {
		fmt.Println("error: --id must be a positive number")
		os.Exit(1)
	}

	tasks, err := storage.LoadJSON(tasksFile)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	mgr := todo.NewManager(tasks)
	if err := mgr.Complete(*id); err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	if err := storage.SaveJSON(tasksFile, mgr.Tasks()); err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	fmt.Printf("completed task #%d\n", *id)
}

func cmdDelete(args []string) {
	fs := flag.NewFlagSet("delete", flag.ExitOnError)
	id := fs.Int("id", 0, "task id")
	fs.Parse(args)

	if *id <= 0 {
		fmt.Println("error: --id must be a positive number")
		os.Exit(1)
	}

	tasks, err := storage.LoadJSON(tasksFile)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	mgr := todo.NewManager(tasks)
	if err := mgr.Delete(*id); err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	if err := storage.SaveJSON(tasksFile, mgr.Tasks()); err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	fmt.Printf("deleted task #%d\n", *id)
}

func cmdExport(args []string) {
	fs := flag.NewFlagSet("export", flag.ExitOnError)
	format := fs.String("format", "json", "output format: csv or json")
	out := fs.String("out", "", "output file path")
	fs.Parse(args)

	if *out == "" {
		fmt.Println("error: --out is required")
		os.Exit(1)
	}

	tasks, err := storage.LoadJSON(tasksFile)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	switch *format {
	case "csv":
		err = storage.SaveCSV(*out, tasks)
	case "json":
		err = storage.SaveJSON(*out, tasks)
	default:
		fmt.Printf("error: unknown format %q (use csv or json)\n", *format)
		os.Exit(1)
	}
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	fmt.Printf("exported %d tasks to %s\n", len(tasks), *out)
}

func cmdLoad(args []string) {
	fs := flag.NewFlagSet("load", flag.ExitOnError)
	file := fs.String("file", "", "file to import (csv or json)")
	fs.Parse(args)

	if *file == "" {
		fmt.Println("error: --file is required")
		os.Exit(1)
	}

	var tasks []todo.Task
	var err error
	if strings.HasSuffix(*file, ".csv") {
		tasks, err = storage.LoadCSV(*file)
	} else {
		tasks, err = storage.LoadJSON(*file)
	}
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	if err := storage.SaveJSON(tasksFile, tasks); err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	fmt.Printf("loaded %d tasks from %s\n", len(tasks), *file)
}

func usage() {
	fmt.Println("usage: todo <command> [flags]")
	fmt.Println("commands: add, list, complete, delete, export, load")
}
