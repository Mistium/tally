package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var allowAll = false

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		printHelp()
		return
	}

	recursive := false
	outputJSON := false
	allowAll = false
	allowed := make(map[string]bool)

	for _, arg := range args {
		switch arg {
		case "-c":
			recursive = false
		case "-d":
			recursive = true
		case "-a":
			recursive = true
			allowAll = true
		case "--json":
			outputJSON = true
		default:
			extensions := strings.Split(arg, ",")
			for _, e := range extensions {
				ext := strings.TrimSpace(e)
				if ext != "" {
					if !strings.HasPrefix(ext, ".") {
						ext = "." + ext
					}
					allowed[ext] = true
				}
			}
		}
	}

	stats := make(map[string]int)

	if recursive {
		countLinesRecursive(".", allowed, stats)
	} else {
		countLinesCurrentDir(".", allowed, stats)
	}

	if outputJSON {
		printJSON(stats)
	} else {
		for k := range allowed {
			if _, ok := stats[k]; !ok {
				stats[k] = 0
			}
		}
		printTable(stats)
	}
}

func countLinesCurrentDir(dir string, allowed map[string]bool, stats map[string]int) {
	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading directory: %v\n", err)
		return
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if !shouldCount(file.Name(), allowed) {
			continue
		}

		path := filepath.Join(dir, file.Name())
		lines, err := countLines(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file %s: %v\n", path, err)
			continue
		}
		ext := filepath.Ext(file.Name())
		stats[ext] += lines
	}
}

func countLinesRecursive(root string, allowed map[string]bool, stats map[string]int) {
	filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if d.IsDir() && path != root && (strings.HasPrefix(d.Name(), ".") || d.Name() == "node_modules") {
			return filepath.SkipDir
		}

		if d.IsDir() {
			return nil
		}

		if !shouldCount(d.Name(), allowed) {
			return nil
		}

		lines, err := countLines(path)
		if err != nil {
			return nil
		}

		ext := filepath.Ext(d.Name())
		stats[ext] += lines
		return nil
	})
}

func shouldCount(name string, allowed map[string]bool) bool {
	if len(allowed) == 0 || allowAll {
		return true
	}
	ext := filepath.Ext(name)
	return allowed[ext]
}

func countLines(filePath string) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := 0
	for scanner.Scan() {
		lines++
	}
	return lines, scanner.Err()
}

func printTable(stats map[string]int) {
	if len(stats) == 0 {
		fmt.Println("No matching files found.")
		return
	}

	fmt.Println("| Extension | Lines of Code | Percentage |")
	fmt.Println("| --------- | ------------- | ---------- |")

	total := 0
	for _, lines := range stats {
		total += lines
	}

	type ExtensionStat struct {
		ext   string
		lines int
	}

	var sortedStats []ExtensionStat
	for ext, lines := range stats {
		sortedStats = append(sortedStats, ExtensionStat{ext, lines})
	}

	for i := 0; i < len(sortedStats); i++ {
		for j := i + 1; j < len(sortedStats); j++ {
			if sortedStats[j].lines > sortedStats[i].lines {
				sortedStats[i], sortedStats[j] = sortedStats[j], sortedStats[i]
			}
		}
	}

	for _, stat := range sortedStats {
		extName := strings.TrimPrefix(stat.ext, ".")
		if extName == "" {
			extName = "(no ext)"
		}
		percentage := float64(stat.lines) / float64(total) * 100
		fmt.Printf("| %-9s | %13d | %9.1f%% |\n", extName, stat.lines, percentage)
	}
	fmt.Printf("\nTotal lines of code: %d\n", total)
}

func printJSON(stats map[string]int) {
	total := 0
	for _, v := range stats {
		total += v
	}
	output := map[string]any{
		"files": stats,
		"total": total,
	}
	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return
	}
	fmt.Println(string(data))
}

func printHelp() {
	fmt.Println("Usage: tally [-c|-d|-a] [filetypes] [--json]")
	fmt.Println("Count lines of code in the current directory or recursively in all subdirectories.")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -c           Count lines of code in the current directory")
	fmt.Println("  -d           Count lines of code recursively in all subdirectories")
	fmt.Println("  -a           Count lines of code recursively in all subdirectories, including all files")
	fmt.Println("  filetypes    Comma-separated file extensions (e.g. osl,go,js)")
	fmt.Println("  --json       Output results as JSON")
}
