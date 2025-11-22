package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func main() {
	// get home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	// create Notes directory
	notesDir := filepath.Join(homeDir, "Notes")
	err = os.MkdirAll(notesDir, 0755)
	if err != nil {
		log.Fatal(err)
	}

	// get current time
	now := time.Now()

	// format note filename (YYYYMMDDSS)
	ext := ".md"
	filename := now.Format("20060102150405") + ext

	// create temp file of note
	tmpFile, err := os.CreateTemp("", filename)
	if err != nil {
		log.Fatal(err)
	}
	tmpPath := tmpFile.Name()

	// write yaml front matter
	initialContent := fmt.Sprintf("---\ncreated: %s\n---\n\n\n", now.Format(time.RFC3339))
	_, err = tmpFile.WriteString(initialContent)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// get initial modification time
	initialInfo, err := os.Stat(tmpPath)
	if err != nil {
		log.Fatal(err)
	}
	initialModTime := initialInfo.ModTime()

	// close temp file
	tmpFile.Close()

	// open temp file in editor
	editor := getEditor()

	var cmd *exec.Cmd
	if editor == "vim" {
		cmd = exec.Command(editor, "+", tmpPath)
	} else {
		cmd = exec.Command(editor, tmpPath)
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	// check if file was modified (saved in vim)
	finalInfo, err := os.Stat(tmpPath)
	if err != nil {
		log.Fatal(err)
	}

	if finalInfo.ModTime().After(initialModTime) {
		savedNotePath := filepath.Join(notesDir, filename)
		err = os.Rename(tmpPath, savedNotePath)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("✓ File saved to", savedNotePath)
	} else {
		os.Remove(tmpPath)
		fmt.Println("✗ File not saved, discarded")
	}
}

func getEditor() string {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		fmt.Fprintln(os.Stderr, "No $EDITOR set. Using vim as default.")
		fmt.Fprintln(os.Stderr, "To change: export EDITOR=nano (or your preferred editor)")
		return "vim"
	}
	return editor
}
