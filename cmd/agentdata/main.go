package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"agentdata/internal/exporter"
	"agentdata/internal/model"
	"agentdata/internal/search"
	"agentdata/internal/source"
)

const version = "agentdata 0.1.0"

func main() {
	os.Exit(Run(os.Args[1:], os.Stdout, os.Stderr))
}

func Run(args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "usage: agentdata <version|scan|search|export>")
		return 2
	}

	switch args[0] {
	case "version":
		fmt.Fprintln(stdout, version)
		return 0
	case "scan":
		return runScan(args[1:], stdout, stderr)
	case "search":
		return runSearch(args[1:], stdout, stderr)
	case "export":
		return runExport(args[1:], stdout, stderr)
	default:
		fmt.Fprintf(stderr, "unknown command: %s\n", args[0])
		return 2
	}
}

func runScan(args []string, stdout io.Writer, stderr io.Writer) int {
	flags := flag.NewFlagSet("scan", flag.ContinueOnError)
	flags.SetOutput(stderr)
	path := flags.String("path", "", "file or directory containing JSONL records")
	if err := flags.Parse(args); err != nil {
		return 2
	}

	sessions, err := loadSessions(*path)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	fmt.Fprintf(stdout, "sessions=%d messages=%d\n", len(sessions), countMessages(sessions))
	return 0
}

func runSearch(args []string, stdout io.Writer, stderr io.Writer) int {
	flags := flag.NewFlagSet("search", flag.ContinueOnError)
	flags.SetOutput(stderr)
	path := flags.String("path", "", "file or directory containing JSONL records")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	query := strings.Join(flags.Args(), " ")
	if strings.TrimSpace(query) == "" {
		fmt.Fprintln(stderr, "search query is required")
		return 2
	}

	sessions, err := loadSessions(*path)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	for _, result := range search.Messages(sessions, query) {
		fmt.Fprintf(stdout, "%s\t%s\t%s\n", result.SessionID, result.Message.Role, result.Message.Text)
	}
	return 0
}

func runExport(args []string, stdout io.Writer, stderr io.Writer) int {
	flags := flag.NewFlagSet("export", flag.ContinueOnError)
	flags.SetOutput(stderr)
	path := flags.String("path", "", "file or directory containing JSONL records")
	format := flags.String("format", "jsonl", "export format: jsonl or markdown")
	if err := flags.Parse(args); err != nil {
		return 2
	}

	sessions, err := loadSessions(*path)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	switch *format {
	case "jsonl":
		err = exporter.JSONL(stdout, sessions)
	case "markdown":
		err = exporter.Markdown(stdout, sessions)
	default:
		fmt.Fprintf(stderr, "unsupported export format: %s\n", *format)
		return 2
	}
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	return 0
}

func loadSessions(path string) ([]model.Session, error) {
	if strings.TrimSpace(path) == "" {
		return nil, errors.New("--path is required")
	}

	files, err := jsonlFiles(path)
	if err != nil {
		return nil, err
	}

	sessions := make([]model.Session, 0)
	for _, file := range files {
		opened, err := os.Open(file)
		if err != nil {
			return nil, err
		}
		parsed, parseErr := source.ParseJSONLSessions(opened)
		closeErr := opened.Close()
		if parseErr != nil {
			return nil, fmt.Errorf("%s: %w", file, parseErr)
		}
		if closeErr != nil {
			return nil, closeErr
		}
		sessions = append(sessions, parsed...)
	}
	return sessions, nil
}

func jsonlFiles(path string) ([]string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return []string{path}, nil
	}

	files := make([]string, 0)
	err = filepath.WalkDir(path, func(current string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}
		if strings.EqualFold(filepath.Ext(current), ".jsonl") {
			files = append(files, current)
		}
		return nil
	})
	return files, err
}

func countMessages(sessions []model.Session) int {
	total := 0
	for _, session := range sessions {
		total += len(session.Messages)
	}
	return total
}
