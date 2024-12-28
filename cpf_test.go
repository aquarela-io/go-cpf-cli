package main

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func TestValidateCPF(t *testing.T) {
	tests := []struct {
		name     string
		cpf      string
		byLength bool
		want     bool
	}{
		{"valid formatted CPF", "529.982.247-25", false, true},
		{"valid unformatted CPF", "52998224725", false, true},
		{"invalid CPF", "113.111.111-11", false, false},
		{"invalid length", "123", false, false},
		{"valid length only", "12345678901", true, true},
		{"invalid length check", "123456", true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validateCPF(tt.cpf, tt.byLength)
			if got != tt.want {
				t.Errorf("validateCPF(%q, %v) = %v, want %v", tt.cpf, tt.byLength, got, tt.want)
			}
		})
	}
}

func TestFormatCPF(t *testing.T) {
	tests := []struct {
		name    string
		cpf     string
		want    string
		wantErr bool
	}{
		{"valid unformatted", "52998224725", "529.982.247-25", false},
		{"already formatted", "529.982.247-25", "529.982.247-25", false},
		{"invalid length", "123", "", true},
		{"with letters", "123abc45678", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := formatCPF(tt.cpf)
			if (err != nil) != tt.wantErr {
				t.Errorf("formatCPF(%q) error = %v, wantErr %v", tt.cpf, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("formatCPF(%q) = %v, want %v", tt.cpf, got, tt.want)
			}
		})
	}
}

func TestGenerateCPF(t *testing.T) {
	t.Run("formatted valid", func(t *testing.T) {
		cpf, err := generateCPF(true, false)
		if err != nil {
			t.Fatalf("generateCPF(true, false) error = %v", err)
		}
		if !strings.Contains(cpf, ".") || !strings.Contains(cpf, "-") {
			t.Errorf("generateCPF(true, false) = %v, want formatted CPF", cpf)
		}
		if !validateCPF(cpf, false) {
			t.Errorf("generateCPF(true, false) = %v, generated invalid CPF", cpf)
		}
	})

	t.Run("unformatted valid", func(t *testing.T) {
		cpf, err := generateCPF(false, false)
		if err != nil {
			t.Fatalf("generateCPF(false, false) error = %v", err)
		}
		if strings.Contains(cpf, ".") || strings.Contains(cpf, "-") {
			t.Errorf("generateCPF(false, false) = %v, want unformatted CPF", cpf)
		}
		if !validateCPF(cpf, false) {
			t.Errorf("generateCPF(false, false) = %v, generated invalid CPF", cpf)
		}
	})
}

func TestProcessFile(t *testing.T) {
	// Create a temporary file with test CPFs
	content := []byte("529.982.247-25\n111.111.111-11\n123.456.789-09\n")
	tmpfile, err := os.CreateTemp("", "cpf_test_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	// Test validation processing
	results, err := processFile(tmpfile.Name(), validateProcessor)
	if err != nil {
		t.Fatalf("processFile() error = %v", err)
	}

	expected := []CPFResult{
		{CPF: "529.982.247-25", Valid: true, Original: "529.982.247-25"},
		{CPF: "111.111.111-11", Valid: false, Original: "111.111.111-11"},
		{CPF: "123.456.789-09", Valid: true, Original: "123.456.789-09"},
	}

	if len(results) != len(expected) {
		t.Fatalf("processFile() returned %d results, want %d", len(results), len(expected))
	}

	for i, result := range results {
		if result.Valid != expected[i].Valid {
			t.Errorf("result[%d].Valid = %v, want %v", i, result.Valid, expected[i].Valid)
		}
	}
}

func TestGenerateMultipleCPFs(t *testing.T) {
	t.Run("generate multiple with newline separator", func(t *testing.T) {
		count := 3
		results, err := generateCPFsJSON(count, true, false)
		if err != nil {
			t.Fatalf("generateCPFsJSON() error = %v", err)
		}

		if len(results) != count {
			t.Errorf("generateCPFsJSON() returned %d CPFs, want %d", len(results), count)
		}

		for i, result := range results {
			if !validateCPF(result.CPF, false) {
				t.Errorf("CPF[%d] = %v is invalid", i, result.CPF)
			}
		}
	})

	t.Run("generate invalid CPFs", func(t *testing.T) {
		count := 3
		results, err := generateCPFsJSON(count, true, true)
		if err != nil {
			t.Fatalf("generateCPFsJSON() error = %v", err)
		}

		if len(results) != count {
			t.Errorf("generateCPFsJSON() returned %d CPFs, want %d", len(results), count)
		}

		validCount := 0
		for _, result := range results {
			if validateCPF(result.CPF, false) {
				validCount++
			}
		}

		// Most generated CPFs should be invalid
		if validCount > count/2 {
			t.Errorf("Too many valid CPFs generated: %d out of %d", validCount, count)
		}
	})
}

func TestFileOutput(t *testing.T) {
	// Test JSON output to file
	tmpfile, err := os.CreateTemp("", "cpf_output_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	// Generate some CPFs
	results, err := generateCPFsJSON(3, true, false)
	if err != nil {
		t.Fatalf("generateCPFsJSON() error = %v", err)
	}

	// Write to file
	output, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	if err := os.WriteFile(tmpfile.Name(), output, 0644); err != nil {
		t.Fatalf("Failed to write to file: %v", err)
	}

	// Read and verify
	content, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	var readResults []CPFResult
	if err := json.Unmarshal(content, &readResults); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if len(readResults) != len(results) {
		t.Errorf("Read %d results from file, want %d", len(readResults), len(results))
	}

	for i, result := range readResults {
		if result.CPF != results[i].CPF {
			t.Errorf("Result[%d] = %v, want %v", i, result.CPF, results[i].CPF)
		}
	}
} 