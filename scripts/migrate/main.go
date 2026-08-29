package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/lukman-ss/software-engineering-lab/pkg/database"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: migrate <up|down>")
		os.Exit(1)
	}

	cmd := os.Args[1]

	ctx := context.Background()
	cfg := database.FromEnv()
	db, err := database.Connect(ctx, cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := runMigrations(ctx, db, cmd); err != nil {
		fmt.Fprintf(os.Stderr, "Migration failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Migration successful.")
}

func runMigrations(ctx context.Context, db *sql.DB, direction string) error {
	dir := "migrations"
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	var targetFiles []string
	suffix := "." + direction + ".sql"

	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), suffix) {
			targetFiles = append(targetFiles, f.Name())
		}
	}

	sort.Strings(targetFiles)
	if direction == "down" {
		// Reverse for down
		for i, j := 0, len(targetFiles)-1; i < j; i, j = i+1, j-1 {
			targetFiles[i], targetFiles[j] = targetFiles[j], targetFiles[i]
		}
	}

	for _, name := range targetFiles {
		path := filepath.Join(dir, name)
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %s: %w", name, err)
		}

		fmt.Printf("Executing %s...\n", name)
		if _, err := db.ExecContext(ctx, string(content)); err != nil {
			return fmt.Errorf("execute %s: %w", name, err)
		}
	}

	return nil
}
