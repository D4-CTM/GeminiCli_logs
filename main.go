package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/pmezard/go-difflib/difflib"
)

const dirCreatePermissions = 0755

// Reads the target *-target-dir* searching for files that matches the pattern of: *checkpoint-*.json*
func readChatLogsDir(logsPath string) []LogMetadata {
	if *verbose {
		fmt.Printf("Getting content on %s\n", logsPath)
	}

	logsDir, err := os.ReadDir(logsPath)

	if err != nil {
		log.Panicf("Error: %v", err)
	}

	return getChatLogs(logsDir, logsPath)
}

func getChatLogs(logsDir []os.DirEntry, path string) []LogMetadata {
	var logsMetadata []LogMetadata
	for _, logDir := range	logsDir {
		chatLogPath := filepath.Join(path, logDir.Name())

		chatLogs, err := os.ReadDir(chatLogPath)
		if err != nil {
			log.Printf("An error ocurred while attempting to read: %s\nError: %\n", chatLogPath, err)
			continue
		}

		for _, logs := range chatLogs {
			logName := logs.Name()
			_chatLogPath := filepath.Join(chatLogPath, logName)

			if logs.IsDir() {
				if *verbose {
					fmt.Printf("%s is a directory!\n", _chatLogPath)
				}

				logsMetadata = append(logsMetadata, readChatLogsDir(_chatLogPath)...)	
				continue
			}

			chopedName, hasPrefix := strings.CutPrefix(logName, "checkpoint-")
			if hasPrefix {
				if *verbose {
					fmt.Printf("%s is a valid log!\n", _chatLogPath)
				}
				
				chatLogName := strings.TrimSuffix(chopedName, ".json")

				logsMetadata = append(logsMetadata, LogMetadata{
					Path:    _chatLogPath,
					LogName: chatLogName,
				})
			}

		}
	}
	return logsMetadata
}

func createLog(metadata LogMetadata) {
	fmt.Printf("-- Reading %s --\n", metadata.LogName)
	file, err := os.Open(metadata.Path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	token, err := decoder.Token()
	if err != nil || token != json.Delim('[') {
		log.Fatal("Expected start of array")
	}

	content := ""
	contentBuffer := ""

	if *outputMarkdown {
		content += fmt.Sprintf("## `%s` Log\n", metadata.LogName)
	}

	for decoder.More() {
		var entry Message
		if err := decoder.Decode(&entry); err != nil {
			log.Fatal(err)
		}

		part := entry.Parts[0]
		if len(part.Text) > 0 {
			if *outputMarkdown {
				content += fmt.Sprintf("* %s:\n\t- %s\n\n", entry.Role, part.Text)
			} else {
				content += fmt.Sprintf("%s:\n> %s\n\n", entry.Role, part.Text)
			}	
			continue
		}

		if part.FunctionCall != nil && part.FunctionCall.Name == "write_file" {
			var _filePath string
			var _content string
			var functionCall = part.FunctionCall

			if field, ok := functionCall.GetField("file_path"); ok {
				_filePath = field
			} else {
				continue
			}

			if field, ok := functionCall.GetField("content"); ok {
				_content = field
			} else {
				continue
			}

			header := fmt.Sprintf("\n-> File written: %s\n", _filePath)
			if *outputMarkdown {
				extension := filepath.Ext(filepath.Base(_filePath))

				contentBuffer = content + fmt.Sprintf("%s``` %s\n%s\n```\n\n", header, extension, _content)
			} else {
				contentBuffer = content + fmt.Sprintf("%s%s\n\n", header, _content)
			}
			continue
		}

		if part.FunctionCall != nil && part.FunctionCall.Name == "replace" {
			var stringA string
			var stringB string
			var _filePath string
			var functionCall = part.FunctionCall

			if field, ok := functionCall.GetField("old_string"); ok {
				stringA = field
			} else {
				continue
			}

			if field, ok := functionCall.GetField("new_string"); ok {
				stringB = field
			} else {
				continue
			}

			if field, ok := functionCall.GetField("file_path"); ok {
				_filePath = field
			} else {
				continue
			}

			diff := difflib.UnifiedDiff{
				A:       difflib.SplitLines(stringA),
				B:       difflib.SplitLines(stringB),
				Context: 3,
			}
			_content, _ := difflib.GetUnifiedDiffString(diff)

			header := fmt.Sprintf("\n-> File modified: %s\n", _filePath)
			if *outputMarkdown {
				contentBuffer = content + fmt.Sprintf("%s``` diff\n%s\n```\n\n", header, _content)
			} else {
				contentBuffer = content + fmt.Sprintf("%s%s\n\n", header, _content)
			}
			continue
		}

		if part.FunctionResponse != nil && part.FunctionResponse.Name == "write_file" {
			if _, ok := part.FunctionResponse.GetField("error"); ok {
				if *outputMarkdown {
					contentBuffer = content + "\n-> **The user canceled the write file operation**\n\n"
				} else {
					contentBuffer = content + "\n-> The user canceled the write file operation\n\n"
				}
			}

			content = contentBuffer
		}

		if part.FunctionResponse != nil && part.FunctionResponse.Name == "replace" {
			if _, ok := part.FunctionResponse.GetField("error"); ok {
				if *outputMarkdown {
					contentBuffer = content + "\n-> **The user canceled the replacement**\n\n"
				} else {
					contentBuffer = content + "\n-> The user canceled the replacement\n\n"
				}
			}

			content = contentBuffer
		}
	}

	_, err = decoder.Token()

	logFile, _ := os.Create(filepath.Join(*outputDir, fmt.Sprintf("%s%s", metadata.LogName, logExtension)))

	logFile.WriteString(content)
}

func CreateOutputDir() {
	_, err := os.Stat(*outputDir)

	if os.IsNotExist(err) {
		if *verbose {
			fmt.Printf("creating \"%s\"!\n", *outputDir)
		}
		err = os.MkdirAll(*outputDir, dirCreatePermissions)
		if err != nil {
			log.Panicln(err)
		}
		return
	}

	if *verbose {
		fmt.Printf("\"%s\" already exists!\n", *outputDir)
	}
}

func main() {
	initFlags()

	logsPath := *targetDir
	if len(logsPath) == 0 {
		log.Panicf("The target directory should be specified!\n--target-dir: %s", logsPath)
	}

	logs := readChatLogsDir(logsPath)
	
	if *listLogs {
		if *targetDir != defaultPath {
			fmt.Println("\033[31mWhen listing the -output-dir is ignored.\033[0m")
		}
		
		for _, _log := range logs {
			fmt.Println(_log.ToString() + "---")
		}
		return 
	}

	CreateOutputDir()
	for _, _log := range logs {
		createLog(_log)
	}
}
