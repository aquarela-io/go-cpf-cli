package cpf

import (
	"testing"
)

func TestUnformatCPF(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"already unformatted", "12345678909", "12345678909"},
		{"formatted CPF", "123.456.789-09", "12345678909"},
		{"with spaces", "123 456 789 09", "12345678909"},
		{"with letters", "123abc456def789ghi09", "12345678909"},
		{"empty string", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UnformatCPF(tt.input); got != tt.expected {
				t.Errorf("UnformatCPF() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsRepeated(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"empty string", "", false},
		{"single char", "1", true},
		{"all zeros", "00000", true},
		{"all ones", "11111", true},
		{"mixed digits", "12345", false},
		{"repeated with suffix", "11111a", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsRepeated(tt.input); got != tt.expected {
				t.Errorf("IsRepeated() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestFormatCPF(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    string
		expectError bool
	}{
		{"valid unformatted", "12345678909", "123.456.789-09", false},
		{"valid formatted", "123.456.789-09", "123.456.789-09", false},
		{"too short", "1234567890", "", true},
		{"too long", "123456789099", "", true},
		{"with letters", "123abc45678", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FormatCPF(tt.input)
			if (err != nil) != tt.expectError {
				t.Errorf("FormatCPF() error = %v, expectError %v", err, tt.expectError)
				return
			}
			if got != tt.expected {
				t.Errorf("FormatCPF() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestValidateCPF(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		byLength  bool
		expected  bool
	}{
		{"valid CPF", "11144477735", false, true},
		{"valid formatted CPF", "111.444.777-35", false, true},
		{"invalid check digit", "11144477734", false, false},
		{"repeated digits", "11111111111", false, false},
		{"too short", "1234567890", false, false},
		{"too long", "123456789012", false, false},
		{"by length valid", "12345678901", true, true},
		{"by length invalid length", "1234567890", true, false},
		{"by length repeated", "11111111111", true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateCPF(tt.input, tt.byLength); got != tt.expected {
				t.Errorf("ValidateCPF() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestGenerateCPF(t *testing.T) {
	tests := []struct {
		name       string
		formatted  bool
		invalid    bool
		wantFormat bool
		wantValid  bool
	}{
		{"valid unformatted", false, false, false, true},
		{"valid formatted", true, false, true, true},
		{"invalid unformatted", false, true, false, false},
		{"invalid formatted", true, true, true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateCPF(tt.formatted, tt.invalid)
			if err != nil {
				t.Errorf("GenerateCPF() error = %v", err)
				return
			}

			// Check formatting
			if tt.wantFormat {
				if len(got) != 14 || got[3] != '.' || got[7] != '.' || got[11] != '-' {
					t.Errorf("GenerateCPF() formatting incorrect = %v", got)
				}
			} else {
				if len(got) != 11 {
					t.Errorf("GenerateCPF() length incorrect = %v", got)
				}
			}

			// Check validity
			isValid := ValidateCPF(got, false)
			if isValid != tt.wantValid {
				t.Errorf("GenerateCPF() validity = %v, want %v", isValid, tt.wantValid)
			}
		})
	}
} 