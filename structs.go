package main

import "fmt"

type funcUtil interface {
	ToString() string;
};

type LogMetadata struct {
	Path    string
	LogName string
}

func (lm *LogMetadata) ToString() string {
	return fmt.Sprintf("Log name: %s\nPath to file: %s\n", lm.LogName, lm.Path)
}

type Part struct {
	Text             string            `json:"text,omitempty"`
	FunctionResponse *FunctionResponse `json:"functionResponse,omitempty"`
	FunctionCall     *FunctionCall     `json:"functionCall,omitempty"`
}

type FunctionCall struct {
	Name string `json:"name,omitempty"`
	Args any    `json:"args,omitempty"`
}

// Function corrected by chatGPT
func (fc *FunctionCall) GetField(field string) (string, bool) {
	if fc.Args == nil {
		return "", false
	}
	if m, ok := fc.Args.(map[string]any); ok {
		if v, ok := m[field]; ok {
			if s, ok := v.(string); ok {
				return s, true
			}
		}
	}
	return "", false
}

type FunctionResponse struct {
	Name     string `json:"name"`
	Response any    `json:"response"`
}

// Function corrected by ChatGPT
func (fr *FunctionResponse) GetField(field string) (string, bool) {
	if fr.Response == nil {
		return "", false
	}
	if m, ok := fr.Response.(map[string]any); ok {
		if v, ok := m[field]; ok {
			if s, ok := v.(string); ok {
				return s, true
			}
		}
	}
	return "", false
}

type Message struct {
	Role  string `json:"role"`
	Parts []Part `json:"parts"`
}
