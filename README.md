# cpf-cli-go

A command-line interface (CLI) tool for validating and working with Brazilian CPF numbers.

## Features

- Validate CPF numbers
- Clean formatting
- Cross-platform support (built in Go)

## Installation

### Using Homebrew (macOS)

```bash
brew tap diegopeixoto/cpf-cli-go
brew install cpf-cli-go
```

### Manual Installation

1. Download the latest release for your platform from the [releases page](https://github.com/diegopeixoto/cpf-cli-go/releases):

   - Windows: `cpf-cli-go_Windows_x86_64.zip`
   - Linux: `cpf-cli-go_Linux_x86_64.tar.gz`
   - macOS: `cpf-cli-go_Darwin_x86_64.zip`

2. Extract the archive:

   ```bash
   # For Windows:
   # Extract the .zip file using Windows Explorer or:
   unzip cpf-cli-go_Windows_x86_64.zip

   # For Linux:
   tar xzf cpf-cli-go_Linux_x86_64.tar.gz

   # For macOS:
   unzip cpf-cli-go_Darwin_x86_64.zip
   ```

3. Add to PATH:

   - **Windows**:
     - Move `cpf.exe` to a directory like `C:\Program Files\cpf-cli-go\`
     - Add that directory to your PATH environment variable
   - **Linux/macOS**:
     ```bash
     sudo mv cpf /usr/local/bin/
     ```

4. Verify the installation:

   ```bash
   # Windows
   cpf.exe --version

   # Linux/macOS
   cpf --version
   ```

### Building from Source

Requirements:

- Go 1.23 or higher

```bash
# Clone the repository
git clone https://github.com/diegopeixoto/cpf-cli-go.git

# Navigate to the project directory
cd cpf-cli-go

# Build the project
go build

# (Optional) Install globally
go install
```

## Usage

```bash
# Validate a CPF
cpf validate 123.456.789-09

# Clean CPF formatting
cpf clean "123.456.789-09"

# Telemetry Management
cpf telemetry enable    # Enable telemetry
cpf telemetry disable   # Disable telemetry
cpf telemetry status    # Check telemetry status
```

## Telemetry

This tool includes optional telemetry to help us understand how it's being used and improve it. We use [PostHog](https://posthog.com/) for telemetry collection. The telemetry:

- Is **disabled by default**
- Only collects anonymous usage data:
  - Commands used
  - Success/error rates
  - OS and architecture
  - CLI version
- Never collects personal information or CPF numbers
- Can be enabled/disabled at any time using the `cpf telemetry` command
- Stores its configuration in `~/.cpf-cli/telemetry.json`

To manage telemetry:

```bash
cpf telemetry enable    # Enable telemetry
cpf telemetry disable   # Disable telemetry
cpf telemetry status    # Check current status
```

### Building with Telemetry

When building from source, you can configure the PostHog API key at build time:

```bash
# Build with PostHog API key
POSTHOG_API_KEY=your_api_key go build -ldflags="-X github.com/diegopeixoto/cpf-cli-go/pkg/telemetry.apiKey=$POSTHOG_API_KEY"

# Or using make (if you have a Makefile)
make build POSTHOG_API_KEY=your_api_key
```

The official releases are built with telemetry enabled and configured to send data to our PostHog instance. This helps us understand how the tool is being used and improve it. You can always disable telemetry after installation using `cpf telemetry disable`.

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Author

Diego Peixoto - [GitHub](https://github.com/diegopeixoto)
