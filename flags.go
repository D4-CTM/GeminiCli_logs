package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func getLogsPath() string {
	homeDir, err := os.UserHomeDir()

	if err != nil {
		log.Panicf("Error: %v", err)
	}

	return filepath.Join(homeDir, ".gemini/tmp")
}

var (
	defaultPath = getLogsPath()
	logExtension = ".txt"

	outputDir = flag.String("output-dir", "logs", "Directory to save output logs.")
	targetDir = flag.String("target-dir", defaultPath, "Directory where Gemini logs are located.")
	verbose = flag.Bool("verbose", false, "Enabling this mode will simply show more output on the terminal.")
	outputMarkdown = flag.Bool("output-markdown", false, "Output as markdown .md file.")
	listLogs = flag.Bool("list", false, "List the Gemini logs found in \"-target-dir\".")
)

func initFlags() {
	flag.StringVar(outputDir, "o", "logs", "Alias for --output-dir.")
	flag.StringVar(targetDir, "t", defaultPath, "Alias for --target-dir.")
	flag.BoolVar(outputMarkdown, "md", false, "Alias for --output-markdown.")
	flag.BoolVar(listLogs, "l", false, "Alias for --list-logs.")
	flag.BoolVar(verbose, "vb", false, "Alias for --verbose.")

	flag.Parse()

	if *verbose {
		fmt.Printf("Output dir: %s\n", *outputDir)
		fmt.Printf("Target dir: %s\n", *targetDir)
		fmt.Printf("Is verbose enabled?\n  > %t\n", *verbose)
		fmt.Printf("Is it gonna be a list?\n  > %t\n", *listLogs)
		fmt.Printf("Is the output markdown?\n  > %t\n", *outputMarkdown)
	}

	if *outputMarkdown {
		logExtension = ".md"
	}
}

