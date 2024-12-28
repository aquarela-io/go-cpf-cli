package main

import (
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
		{"invalid CPF", "123.111.111-11", false, false},
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