package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/diegopeixoto/cpf-cli-go/pkg/cpf"
	"github.com/diegopeixoto/cpf-cli-go/pkg/telemetry"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

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
  version, -V          Show version information.
  help, -h, --help     Show this help message.
  telemetry            Configure telemetry settings.

Telemetry Commands:
  telemetry enable              Enable telemetry
  telemetry disable             Disable telemetry
  telemetry status              Show telemetry status

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
  cpf format --file=cpfs.txt --output=formatted.json
  cpf telemetry enable               Enable telemetry
  cpf telemetry disable              Disable telemetry
  cpf telemetry status               Show telemetry status`

	fmt.Printf(help, time.Now().Year())
	fmt.Println()
}

func main() {
	// Initialize telemetry
	if err := telemetry.Initialize(version); err != nil {
		// Silently continue if telemetry initialization fails
		_ = err
	}
	// Ensure we close the telemetry client
	defer telemetry.Close()

	args := os.Args[1:]

	if len(args) == 0 {
		printHelp()
		os.Exit(0)
	}

	command := strings.ToLower(args[0])

	// Handle telemetry commands first
	if command == "telemetry" {
		if len(args) < 2 {
			fmt.Println("Usage: cpf telemetry [enable|disable|status]")
			os.Exit(1)
		}

		switch strings.ToLower(args[1]) {
		case "enable":
			if err := telemetry.SetEnabled(true); err != nil {
				fmt.Fprintf(os.Stderr, "Error enabling telemetry: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("Telemetry enabled")
			return
		case "disable":
			if err := telemetry.SetEnabled(false); err != nil {
				fmt.Fprintf(os.Stderr, "Error disabling telemetry: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("Telemetry disabled")
			return
		case "status":
			if telemetry.IsEnabled() {
				fmt.Println("Telemetry is enabled")
			} else {
				fmt.Println("Telemetry is disabled")
			}
			return
		default:
			fmt.Println("Usage: cpf telemetry [enable|disable|status]")
			os.Exit(1)
		}
	}

	// Track command execution
	defer func() {
		metadata := make(map[string]string)
		if len(args) > 1 {
			metadata["args"] = strings.Join(args[1:], " ")
		}
		telemetry.Track(command, true, nil, metadata)
	}()

	switch command {
	case "version":
		printVersion()
		return
	case "help", "--help", "-h":
		printHelp()
		return
	case "validate", "-v":
		var results []cpf.CPFResult
		var err error

		// Check if we're processing a file
		hasFile := false
		for i := 1; i < len(args); i++ {
			if strings.HasPrefix(args[i], "--file=") {
				hasFile = true
				filename := strings.TrimPrefix(args[i], "--file=")
				results, err = cpf.ProcessFile(filename, cpf.ValidateProcessor)
				if err != nil {
					telemetry.Track(command, false, err, nil)
					fmt.Fprintf(os.Stderr, "Error: %v\n", err)
					os.Exit(1)
				}
				break
			}
		}

		if !hasFile {
			if len(args) < 2 {
				err := fmt.Errorf("missing CPF to validate")
				telemetry.Track(command, false, err, nil)
				fmt.Fprintln(os.Stderr, "Error: Missing CPF to validate.")
				printHelp()
				os.Exit(1)
			}
			// Single CPF validation
			cpfToValidate := args[1]
			results = []cpf.CPFResult{cpf.ValidateProcessor(cpfToValidate)}
		}

		// Check if we should output to a file
		outputFile := ""
		for i := 1; i < len(args); i++ {
			if strings.HasPrefix(args[i], "--output=") {
				outputFile = strings.TrimPrefix(args[i], "--output=")
				break
			}
		}

		if err := cpf.WriteJSONOutput(results, outputFile); err != nil {
			telemetry.Track(command, false, err, nil)
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

	case "format", "-f":
		if len(args) < 2 {
			err := fmt.Errorf("missing CPF to format")
			telemetry.Track(command, false, err, nil)
			fmt.Fprintln(os.Stderr, "Error: Missing CPF to format.")
			printHelp()
			os.Exit(1)
		}
		cpfToFormat := args[1]
		formatted, err := cpf.FormatCPF(cpfToFormat)
		if err != nil {
			telemetry.Track(command, false, err, nil)
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
					telemetry.Track(command, false, err, nil)
					fmt.Fprintf(os.Stderr, "Error: Invalid count value '%s'. Must be a positive number.\n", countStr)
					os.Exit(1)
				}
				count = n
			case strings.HasPrefix(arg, "--separator="):
				separator = strings.TrimPrefix(arg, "--separator=")
			case strings.HasPrefix(arg, "--output="):
				outputFile = strings.TrimPrefix(arg, "--output=")
			default:
				err := fmt.Errorf("unknown option '%s'", arg)
				telemetry.Track(command, false, err, nil)
				fmt.Fprintf(os.Stderr, "Error: Unknown option '%s'\n", arg)
				printHelp()
				os.Exit(1)
			}
		}

		if useJSON {
			results, err := cpf.GenerateCPFsJSON(count, !unformatted, invalid)
			if err != nil {
				telemetry.Track(command, false, err, nil)
				fmt.Fprintf(os.Stderr, "Error generating CPFs: %v\n", err)
				os.Exit(1)
			}

			if err := cpf.WriteJSONOutput(results, outputFile); err != nil {
				telemetry.Track(command, false, err, nil)
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
		} else {
			// Generate multiple CPFs
			cpfs := make([]string, 0, count)
			for i := 0; i < count; i++ {
				generatedCPF, err := cpf.GenerateCPF(!unformatted, invalid)
				if err != nil {
					telemetry.Track(command, false, err, nil)
					fmt.Fprintf(os.Stderr, "Error generating CPF: %v\n", err)
					os.Exit(1)
				}
				cpfs = append(cpfs, generatedCPF)
			}
			fmt.Print(strings.Join(cpfs, separator))
			if separator == "\n" {
				fmt.Println()
			}
		}

	default:
		err := fmt.Errorf("unknown command '%s'", command)
		telemetry.Track(command, false, err, nil)
		fmt.Fprintf(os.Stderr, "Error: Unknown command '%s'\n", command)
		printHelp()
		os.Exit(1)
	}
} 