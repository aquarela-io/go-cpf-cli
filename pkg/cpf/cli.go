package cpf

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// CPFResult represents the result of a CPF operation
type CPFResult struct {
	CPF      string `json:"cpf"`
	Valid    bool   `json:"valid,omitempty"`
	Error    string `json:"error,omitempty"`
	Original string `json:"original,omitempty"`
}

// ProcessFile processes CPFs from a file using the provided processor function
func ProcessFile(filename string, processFunc func(string) CPFResult) ([]CPFResult, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var results []CPFResult
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		results = append(results, processFunc(line))
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return results, nil
}

// ValidateProcessor creates a CPFResult for validation
func ValidateProcessor(cpf string) CPFResult {
	return CPFResult{
		CPF:      cpf,
		Valid:    ValidateCPF(cpf, false),
		Original: cpf,
	}
}

// FormatProcessor creates a CPFResult for formatting
func FormatProcessor(cpf string) CPFResult {
	formatted, err := FormatCPF(cpf)
	if err != nil {
		return CPFResult{
			CPF:      cpf,
			Error:    err.Error(),
			Original: cpf,
		}
	}
	return CPFResult{
		CPF:      formatted,
		Original: cpf,
	}
}

// GenerateCPFsJSON generates multiple CPFs in JSON format
func GenerateCPFsJSON(count int, formatted, invalid bool) ([]CPFResult, error) {
	results := make([]CPFResult, 0, count)
	for i := 0; i < count; i++ {
		cpf, err := GenerateCPF(formatted, invalid)
		if err != nil {
			return nil, err
		}
		results = append(results, CPFResult{CPF: cpf})
	}
	return results, nil
}

// WriteJSONOutput writes JSON results to a file or stdout
func WriteJSONOutput(results []CPFResult, outputFile string) error {
	output, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %w", err)
	}

	if outputFile != "" {
		if err := os.WriteFile(outputFile, output, 0644); err != nil {
			return fmt.Errorf("error writing to file: %w", err)
		}
		return nil
	}

	fmt.Println(string(output))
	return nil
} 