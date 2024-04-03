package templates

import (
	"encoding/json"
	"os"
	"path/filepath"

	gitignore "github.com/sabhiram/go-gitignore"
)

type Node struct {
	Name     string  `json:"name"`
	IsDir    bool    `json:"isDir"`
	Content  []byte  `json:"content,omitempty"`
	Children []*Node `json:"children,omitempty"`
}

func WalkDir(path string, ignorePatterns *gitignore.GitIgnore, parent *Node) (*Node, error) {
	rootNode, err := addNode(path, parent)
	if err != nil {
		return nil, err
	}

	if !rootNode.IsDir {
		return rootNode, nil
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		fullPath := filepath.Join(path, entry.Name())

		if ignorePatterns != nil && ignorePatterns.MatchesPath(fullPath) {
			continue
		}

		if entry.IsDir() {
			_, err := WalkDir(fullPath, ignorePatterns, rootNode)
			if err != nil {
				return nil, err
			}
		} else {
			_, err := addNode(fullPath, rootNode)
			if err != nil {
				return nil, err
			}
		}
	}

	return rootNode, nil
}

func addNode(path string, parent *Node) (*Node, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	node := &Node{
		Name:  info.Name(),
		IsDir: info.IsDir(),
	}

	if !node.IsDir {
		content, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}

		node.Content = content
	} else {
		node.Children = []*Node{}
	}

	if parent != nil {
		parent.Children = append(parent.Children, node)
	}

	return node, nil
}

func SaveTree(root *Node) error {
	file, err := os.Create("template.json")
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	return encoder.Encode(root)
}

func LoadTree() (*Node, error) {
	file, err := os.Open("template.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var root Node
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&root); err != nil {
		return nil, err
	}

	return &root, nil
}

func WalkTree(node *Node, path string) error {
	currentPath := filepath.Join(path, node.Name)
	if node.IsDir {
		err := os.Mkdir(currentPath, 0755)
		if err != nil {
			return err
		}

		for _, child := range node.Children {
			err := WalkTree(child, currentPath)
			if err != nil {
				return err
			}
		}
	} else {
		file, err := os.Create(currentPath)
		if err != nil {
			return err
		}
		defer file.Close()

		if _, err := file.WriteString(string(node.Content)); err != nil {
			return err
		}
	}

	return nil
}
