package commands

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/okzmo/nixus/cmd/templates"
	gitignore "github.com/sabhiram/go-gitignore"
)

func Execute(commandName string) {
	switch commandName {
	case "save":
		var path string
		if len(os.Args) < 3 {
			path = "./"
		} else {
			path = os.Args[2]
		}

		saveTemplate(path)
	case "create":
		if len(os.Args) < 3 {
			fmt.Println("Please give a name to your project. e.g. nixus create <name>.")
			os.Exit(2)
		}

		projectName := os.Args[2]

		createNewProject(projectName)
	default:
		fmt.Println("Unknown command. Usage: nixus <command>")
		os.Exit(2)
	}
}

func saveTemplate(path string) {
	fmt.Println("saving...")

	ignorePatterns, err := gitignore.CompileIgnoreFile(filepath.Join(path, ".gitignore"))
	if err != nil {
		fmt.Println("Failed to compile .gitignore patterns:", err)
		return
	}

	node, err := templates.WalkDir(path, ignorePatterns, nil)
	if err != nil {
		log.Fatal(err)
	}

	if err := templates.SaveTree(node); err != nil {
		log.Fatal(err)
	}
}

func createNewProject(projectName string) {
	node, err := templates.LoadTree()
	if err != nil {
		log.Fatal(err)
	}

	node.Name = projectName
	err = templates.WalkTree(node, "./")
	if err != nil {
		log.Fatal(err)
	}
}
