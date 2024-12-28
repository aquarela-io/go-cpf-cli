package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// unformatCPF removes all non-digit characters from the input string.
func unformatCPF(cpfStr string) string {
	re := regexp.MustCompile(`\D`)
	return re.ReplaceAllString(cpfStr, "")
}

// isRepeated checks if the string is composed entirely of the same character
// (e.g. "00000000000 or 11111111111G").
func isRepeated(s string) bool {
	if len(s) == 0 {
		return false
	}
	firstChar := s[0]
	for i := 1; i < len(s); i++ {
		if s[i] != firstChar {
			return false
		}
	}
	return true
}

// formatCPF formats an 11-digit CPF string as ###.###.###-##.
func formatCPF(cpfStr string) (string, error) {
	digits := unformatCPF(cpfStr)
	if len(digits) != 11 {
		return "", fmt.Errorf("invalid CPF number (must have 11 digits)")
	}
	return fmt.Sprintf("%s.%s.%s-%s",
		digits[0:3],
		digits[3:6],
		digits[6:9],
		digits[9:11],
	), nil
}

// calc implements the internal sum with the CPF weighting logic.
func calc(nums []int) int {
	total := 0
	for i, num := range nums {
		// multiply `num` by (9 - (i % 10)) and accumulate
		total += num * (9 - (i % 10))
	}
	return total
}

// getCD computes the 2 check digits (DV) from the first 9 digits of a CPF.
func getCD(digits9 []int) ([2]int, error) {
	if len(digits9) != 9 {
		return [2]int{}, fmt.Errorf("invalid digits length: expected 9, got %d", len(digits9))
	}

	// reverse the 9 digits
	reversed := make([]int, 9)
	for i := 0; i < 9; i++ {
		reversed[i] = digits9[9-1-i]
	}

	cd1 := calc(reversed) % 11 % 10

	// for the second digit, pretend there's a '0' at the front + incorporate cd1
	secondInput := append([]int{0}, reversed...)
	cd2 := (calc(secondInput) + cd1*9) % 11 % 10

	return [2]int{cd1, cd2}, nil
}

// validateCPF checks if the provided CPF string is valid. If `byLength` is true,
// we only validate the length (11 digits) and that it isn't repeated.
func validateCPF(cpfStr string, byLength bool) bool {
	unformatted := unformatCPF(cpfStr)
	if len(unformatted) != 11 {
		return false
	}
	if isRepeated(unformatted) {
		return false
	}

	// If we only validate by length, skip check digits
	if byLength {
		return true
	}

	number9 := unformatted[:9]
	dv2 := unformatted[9:11]

	// convert the first 9 digits to int
	var digits9 []int
	for _, ch := range number9 {
		d, err := strconv.Atoi(string(ch))
		if err != nil {
			return false
		}
		digits9 = append(digits9, d)
	}

	cd, err := getCD(digits9)
	if err != nil {
		return false
	}
	trueDV := fmt.Sprintf("%d%d", cd[0], cd[1])

	return dv2 == trueDV
}

// generateCPF creates a random CPF number. If `invalid` is true, the check digits
// are randomized (most likely not matching the real digits). If `formatted` is true,
// it returns a string in ###.###.###-## form.
func generateCPF(formatted, invalid bool) (string, error) {
	digits9 := make([]int, 9)
	for i := 0; i < 9; i++ {
		digits9[i] = rng.Intn(10)
	}

	var dv [2]int
	if invalid {
		dv[0] = rng.Intn(10)
		dv[1] = rng.Intn(10)
	} else {
		correctDV, err := getCD(digits9)
		if err != nil {
			return "", fmt.Errorf("failed to generate check digits: %w", err)
		}
		dv = correctDV
	}

	// build the full 11-digit CPF
	allDigits := append(digits9, dv[0], dv[1])
	cpfStr := strings.Builder{}
	for _, d := range allDigits {
		cpfStr.WriteString(strconv.Itoa(d))
	}

	if formatted {
		result, err := formatCPF(cpfStr.String())
		if err != nil {
			return "", fmt.Errorf("failed to format CPF: %w", err)
		}
		return result, nil
	}
	return cpfStr.String(), nil
}


func printVersion() {
	fmt.Printf("CPF Tool version %s (%s) built on %s\n", version, commit, date)
	fmt.Println("Developed by Diego Peixoto for aquarela.io")
	fmt.Printf("Copyleft © 2024-%d\n", time.Now().Year())
}

func printHelp() {
	help := `CPF Tool
Developed by Diego Peixoto for aquarela.io
Copyleft © 2024-%d

Usage:
  cpf <command> [options]

Commands:
  validate, -v          Validate CPF(s). Use --file to validate from file.
  format, -f <cpf>      Format a given CPF to ###.###.###-##.
  generate, -g          Generate random CPF(s).
  version, -V           Show version information.
  help, -h, --help     Show this help message.

Options for "generate":
  --invalid          Generate invalid CPF(s).
  --unformatted     Generate unformatted CPF(s).
  --count=N         Generate N CPFs (default: 1).
  --separator=X     Separator between multiple CPFs (default: newline).
  --json            Output in JSON format.

File processing:
  --file=FILE       Process CPFs from a file (one per line).
  --output=FILE     Write output to a file instead of stdout.

Examples:
  cpf -v 123.456.789-09              Validate a single CPF
  cpf validate --file=cpfs.txt       Validate CPFs from file
  cpf -f 12345678909                 Format a CPF
  cpf -g                             Generate a CPF
  cpf -g --invalid --json            Generate invalid CPF in JSON format
  cpf format --file=cpfs.txt --output=formatted.json`

	fmt.Printf(help, time.Now().Year())
	fmt.Println()
}


func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		printHelp()
		os.Exit(0)
	}

	command := strings.ToLower(args[0])

	switch command {
	case "version":
		printVersion()
		return
	case "help", "--help", "-h":
		printHelp()
		return
	case "validate", "-v":
		var results []CPFResult
		var err error

		// Check if we're processing a file
		hasFile := false
		for i := 1; i < len(args); i++ {
			if strings.HasPrefix(args[i], "--file=") {
				hasFile = true
				filename := strings.TrimPrefix(args[i], "--file=")
				results, err = processFile(filename, validateProcessor)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error: %v\n", err)
					os.Exit(1)
				}
				break
			}
		}

		if !hasFile {
			if len(args) < 2 {
				fmt.Fprintln(os.Stderr, "Error: Missing CPF to validate.")
				printHelp()
				os.Exit(1)
			}
			// Single CPF validation
			cpfToValidate := args[1]
			results = []CPFResult{validateProcessor(cpfToValidate)}
		}

		// Check if we should output to a file
		outputFile := ""
		for i := 1; i < len(args); i++ {
			if strings.HasPrefix(args[i], "--output=") {
				outputFile = strings.TrimPrefix(args[i], "--output=")
				break
			}
		}

		// Output results
		output, err := json.MarshalIndent(results, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if outputFile != "" {
			if err := os.WriteFile(outputFile, output, 0644); err != nil {
				fmt.Fprintf(os.Stderr, "Error writing to file: %v\n", err)
				os.Exit(1)
			}
		} else {
			fmt.Println(string(output))
		}

	case "format", "-f":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Error: Missing CPF to format.")
			printHelp()
			os.Exit(1)
		}
		cpfToFormat := args[1]
		formatted, err := formatCPF(cpfToFormat)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(formatted)

	case "generate", "-g":
		invalid := false
		unformatted := false
		count := 1
		separator := "\n"
		useJSON := false
		outputFile := ""

		for i := 1; i < len(args); i++ {
			arg := args[i]
			switch {
			case arg == "--invalid":
				invalid = true
			case arg == "--unformatted":
				unformatted = true
			case arg == "--json":
				useJSON = true
			case strings.HasPrefix(arg, "--count="):
				countStr := strings.TrimPrefix(arg, "--count=")
				n, err := strconv.Atoi(countStr)
				if err != nil || n <= 0 {
					fmt.Fprintf(os.Stderr, "Error: Invalid count value '%s'. Must be a positive number.\n", countStr)
					os.Exit(1)
				}
				count = n
			case strings.HasPrefix(arg, "--separator="):
				separator = strings.TrimPrefix(arg, "--separator=")
			case strings.HasPrefix(arg, "--output="):
				outputFile = strings.TrimPrefix(arg, "--output=")
			default:
				fmt.Fprintf(os.Stderr, "Error: Unknown option '%s'\n", arg)
				printHelp()
				os.Exit(1)
			}
		}

		if useJSON {
			results, err := generateCPFsJSON(count, !unformatted, invalid)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error generating CPFs: %v\n", err)
				os.Exit(1)
			}

			output, err := json.MarshalIndent(results, "", "  ")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}

			if outputFile != "" {
				if err := os.WriteFile(outputFile, output, 0644); err != nil {
					fmt.Fprintf(os.Stderr, "Error writing to file: %v\n", err)
					os.Exit(1)
				}
			} else {
				fmt.Println(string(output))
			}
		} else {
			// Generate multiple CPFs
			cpfs := make([]string, 0, count)
			for i := 0; i < count; i++ {
				cpf, err := generateCPF(!unformatted, invalid)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error generating CPF: %v\n", err)
					os.Exit(1)
				}
				cpfs = append(cpfs, cpf)
			}
			fmt.Print(strings.Join(cpfs, separator))
			if separator == "\n" {
				fmt.Println()
			}
		}

	default:
		fmt.Fprintf(os.Stderr, "Error: Unknown command '%s'\n", command)
		printHelp()
		os.Exit(1)
	}
}

// Add these new types for JSON output
type CPFResult struct {
	CPF     string `json:"cpf"`
	Valid   bool   `json:"valid,omitempty"`
	Error   string `json:"error,omitempty"`
	Original string `json:"original,omitempty"`
}

// Add this new function to process files
func processFile(filename string, processFunc func(string) CPFResult) ([]CPFResult, error) {
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

// Add validation processor
func validateProcessor(cpf string) CPFResult {
	return CPFResult{
		CPF:      cpf,
		Valid:    validateCPF(cpf, false),
		Original: cpf,
	}
}

// Add format processor
func formatProcessor(cpf string) CPFResult {
	formatted, err := formatCPF(cpf)
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

// Update generateCPF to support JSON output
func generateCPFsJSON(count int, formatted, invalid bool) ([]CPFResult, error) {
	results := make([]CPFResult, 0, count)
	for i := 0; i < count; i++ {
		cpf, err := generateCPF(formatted, invalid)
		if err != nil {
			return nil, err
		}
		results = append(results, CPFResult{CPF: cpf})
	}
	return results, nil
}
