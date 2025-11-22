// Package main implements a CLI writing tool that opens
// an editor for quick note-taking.
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const (
	programName  = "Rue"
	notesDirName = "Notes"
	noteExt      = ".txt"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("getting home directory: %w", err)
	}

	topLevelPath := filepath.Join(homeDir, programName)
	err = os.MkdirAll(topLevelPath, 0755)
	if err != nil {
		return fmt.Errorf("creating top-level directory: %w", err)
	}

	now := time.Now()
	filename := now.Format("20060102150405") + noteExt

	tempFile, err := os.CreateTemp("", filename)
	if err != nil {
		return fmt.Errorf("creating temp file: %w", err)
	}
	tempPath := tempFile.Name()
	defer tempFile.Close()

	tempInfo, err := os.Stat(tempPath)
	if err != nil {
		return fmt.Errorf("getting initial file info: %w", err)
	}
	initialModTime := tempInfo.ModTime()

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}

	cmd := exec.Command(editor, tempPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("opening editor: %w", err)
	}

	tempInfo, err = os.Stat(tempPath)
	if err != nil {
		return fmt.Errorf("getting final file info: %w", err)
	}
	finalModTime := tempInfo.ModTime()

	if finalModTime.After(initialModTime) {
		notesPath := filepath.Join(topLevelPath, notesDirName)
		err = os.MkdirAll(notesPath, 0755)
		if err != nil {
			return fmt.Errorf("creating notes directory: %w", err)
		}

		savedPath := filepath.Join(notesPath, filename)
		err = os.Rename(tempPath, savedPath)
		if err != nil {
			return fmt.Errorf("saving note: %w", err)
		}

		fmt.Printf("Note saved at %s\n", savedPath)
	} else {
		err = os.Remove(tempPath)
		if err != nil {
			return fmt.Errorf("removing temp file: %w", err)
		}
		fmt.Println("Note discarded")
	}

	return nil
}
