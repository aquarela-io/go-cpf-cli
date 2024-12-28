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

### Linux and Manual Installation

1. Download the latest release for your platform (Linux/macOS) from the [releases page](https://github.com/diegopeixoto/cpf-cli-go/releases)
2. Extract the archive:
   ```bash
   # For Linux:
   tar xzf cpf-cli-go_Linux_x86_64.tar.gz
   # For macOS:
   unzip cpf-cli-go_Darwin_x86_64.zip
   ```
3. Move the binary to your PATH:
   ```bash
   sudo mv cpf /usr/local/bin/
   ```
4. Verify the installation:
   ```bash
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
```

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
