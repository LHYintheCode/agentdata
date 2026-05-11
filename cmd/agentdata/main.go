package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/LHYintheCode/agentdata/internal/exporter"
	"github.com/LHYintheCode/agentdata/internal/model"
	"github.com/LHYintheCode/agentdata/internal/search"
	"github.com/LHYintheCode/agentdata/internal/source"
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
	sourceName := flags.String("source", "jsonl", "input source: jsonl, codex, or claude")
	if err := flags.Parse(args); err != nil {
		return 2
	}

	sessions, err := loadSessions(*sourceName, *path)
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
	sourceName := flags.String("source", "jsonl", "input source: jsonl, codex, or claude")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	query := strings.Join(flags.Args(), " ")
	if strings.TrimSpace(query) == "" {
		fmt.Fprintln(stderr, "search query is required")
		return 2
	}

	sessions, err := loadSessions(*sourceName, *path)
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
	sourceName := flags.String("source", "jsonl", "input source: jsonl, codex, or claude")
	format := flags.String("format", "jsonl", "export format: jsonl or markdown")
	if err := flags.Parse(args); err != nil {
		return 2
	}

	sessions, err := loadSessions(*sourceName, *path)
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

func loadSessions(sourceName string, path string) ([]model.Session, error) {
	resolvedPath, err := resolveSourcePath(sourceName, path)
	if err != nil {
		return nil, err
	}

	files, err := sourceFiles(resolvedPath, sourceName)
	if err != nil {
		return nil, err
	}

	sessions := make([]model.Session, 0)
	for _, file := range files {
		opened, err := os.Open(file)
		if err != nil {
			return nil, err
		}
		parsed, parseErr := parseSessions(sourceName, opened, file)
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

func resolveSourcePath(sourceName string, path string) (string, error) {
	if strings.TrimSpace(path) != "" {
		return path, nil
	}
	if sourceName == "codex" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, ".codex", "sessions"), nil
	}
	if sourceName == "claude" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, ".claude", "projects"), nil
	}
	return "", errors.New("--path is required")
}

func parseSessions(sourceName string, reader io.Reader, file string) ([]model.Session, error) {
	switch sourceName {
	case "jsonl":
		return source.ParseJSONLSessions(reader)
	case "codex":
		session, err := source.ParseCodexRollout(reader, file)
		if err != nil {
			return nil, err
		}
		if len(session.Messages) == 0 {
			return nil, nil
		}
		return []model.Session{session}, nil
	case "claude":
		session, err := source.ParseClaudeTranscript(reader, file)
		if err != nil {
			return nil, err
		}
		if len(session.Messages) == 0 {
			return nil, nil
		}
		return []model.Session{session}, nil
	default:
		return nil, fmt.Errorf("unsupported source: %s", sourceName)
	}
}

func sourceFiles(path string, sourceName string) ([]string, error) {
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
		if matchesSourceFile(current, sourceName) {
			files = append(files, current)
		}
		return nil
	})
	return files, err
}

func matchesSourceFile(path string, sourceName string) bool {
	switch sourceName {
	case "codex":
		return strings.HasPrefix(filepath.Base(path), "rollout-") && strings.EqualFold(filepath.Ext(path), ".jsonl")
	case "claude":
		return source.IsClaudeTranscriptPath(path)
	default:
		return strings.EqualFold(filepath.Ext(path), ".jsonl")
	}
}

func countMessages(sessions []model.Session) int {
	total := 0
	for _, session := range sessions {
		total += len(session.Messages)
	}
	return total
}
