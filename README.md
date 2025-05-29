# Git Report Generator

A command-line tool for generating professional PDF reports of Git commits for specified time periods and authors.

## Features

- 📊 Generate PDF reports of Git commits
- 🎯 Filter commits by author, date range, and branch
- 🎨 Configurable header templates
- 📝 Professional Polish document format
- ⚙️ Customizable PDF styling
- 🔧 Easy-to-use CLI interface

## Installation

### Prerequisites

- Go 1.21 or higher
- Git repository access

### Build from Source

```bash
# Clone the repository
git clone <repository-url>
cd git-report-generator

# Download dependencies
go mod tidy

# Build the application
go build -o git-report-generator

# (Optional) Install globally
go install
```

## Usage

### Basic Usage

```bash
# Generate report for current repository
./git-report-generator --from 2024-01-01 --to 2024-01-31

# Generate report for specific repository
./git-report-generator --repo /path/to/repo --from 2024-01-01 --to 2024-01-31

# Generate report for specific author
./git-report-generator --from 2024-01-01 --to 2024-01-31 --author john@example.com

# Generate report for specific branch
./git-report-generator --from 2024-01-01 --to 2024-01-31 --branch feature/new-feature
```

### Command Line Options

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--repo` | `-r` | Path to Git repository | `.` (current directory) |
| `--from` | `-f` | Start date (YYYY-MM-DD) | **Required** |
| `--to` | `-t` | End date (YYYY-MM-DD) | **Required** |
| `--output` | `-o` | Output PDF file path | `report_YYYY-MM-DD.pdf` |
| `--author` | `-a` | Author email filter | Git config user.email |
| `--branch` | `-b` | Branch name | Current branch |
| `--config` | `-c` | Configuration file path | Default config |

### Examples

```bash
# Generate report for March 2024
./git-report-generator --from 2024-03-01 --to 2024-03-31

# Generate report with custom output file
./git-report-generator --from 2024-01-01 --to 2024-01-31 --output reports/january-2024.pdf

# Generate report for specific author and branch
./git-report-generator \
  --from 2024-01-01 \
  --to 2024-01-31 \
  --author kamil.kalandyk@codewise.com \
  --branch main \
  --output kamil-january-report.pdf
```

## Configuration

The tool supports configuration files for customizing the report header and PDF styling.

### Default Configuration

The default configuration includes:

```json
{
  "header": {
    "template": "Kraków, {{current_date}}\nProtokół odbioru prac programistycznych\n\nWykonawca: {{executor_name}} ({{executor_email}})\nOdbiorca: {{recipient_name}}\n\nRepozytorium: {{repository_name}}\n- Branch {{branch_name}}\n- Commits:",
    "executor_name": "Kamil Kalandyk",
    "executor_email": "kamil.kalandyk@codewise.com,kamil.kalandyk@cm.tech,kamil@mind-future.com",
    "recipient_name": "Commerce Media Tech Sp. z o. o.",
    "location": "Kraków"
  },
  "pdf": {
    "margin_top": 20,
    "margin_bottom": 20,
    "margin_left": 20,
    "margin_right": 20,
    "font_family": "Arial",
    "font_size": 10,
    "header_color": [0, 0, 0],
    "content_color": [50, 50, 50]
  }
}
```

### Custom Configuration

Create a configuration file and use it with the `--config` flag:

```bash
# Create custom config
cat > my-config.json << EOF
{
  "header": {
    "template": "{{current_date}}\nDevelopment Report\n\nDeveloper: {{executor_name}}\nProject: {{repository_name}}\nBranch: {{branch_name}}\n\nCommits:",
    "executor_name": "Your Name",
    "executor_email": "your.email@company.com",
    "recipient_name": "Your Company Ltd."
  }
}
EOF

# Use custom config
./git-report-generator --config my-config.json --from 2024-01-01 --to 2024-01-31
```

### Template Placeholders

Available placeholders for the header template:

- `{{current_date}}` - Current date (YYYY-MM-DD)
- `{{executor_name}}` - Developer/executor name
- `{{executor_email}}` - Developer/executor email
- `{{recipient_name}}` - Recipient organization name
- `{{repository_name}}` - Git repository name
- `{{branch_name}}` - Git branch name

## Report Format

The generated PDF report includes:

1. **Header Section**: Configurable header with date, executor info, recipient, and repository details
2. **Commits Section**: List of commits with:
   - Date (YYYY-MM-DD format)
   - Short SHA (8 characters)
   - Commit message
   - Commit description (if available)

### Sample Output Format

```
Kraków, 2024-01-15
Protokół odbioru prac programistycznych

Wykonawca: Kamil Kalandyk (kamil.kalandyk@codewise.com,kamil.kalandyk@cm.tech,kamil@mind-future.com)
Odbiorca: Commerce Media Tech Sp. z o. o.

Repozytorium: my-project
- Branch main
- Commits:

2024-01-15 a1b2c3d4 Add user authentication feature
    Implemented JWT-based authentication with refresh tokens

2024-01-14 e5f6g7h8 Fix database connection issue
    Resolved connection pool timeout problems

2024-01-13 i9j0k1l2 Update documentation
```

## Development

### Project Structure

```
git-report-generator/
├── cmd/                    # CLI commands
│   └── root.go            # Root command implementation
├── internal/              # Internal packages
│   ├── config/           # Configuration management
│   │   └── config.go
│   ├── git/              # Git operations
│   │   └── service.go
│   └── generator/        # PDF generation
│       └── pdf.go
├── main.go               # Application entry point
├── go.mod                # Go module definition
├── go.sum                # Go module checksums
└── README.md             # This file
```

### Dependencies

- [cobra](https://github.com/spf13/cobra) - CLI framework
- [go-git](https://github.com/go-git/go-git) - Git operations in Go
- [gofpdf](https://github.com/jung-kurt/gofpdf) - PDF generation

### Building

```bash
# Download dependencies
go mod tidy

# Run tests
go test ./...

# Build for current platform
go build -o git-report-generator

# Build for multiple platforms
GOOS=linux GOARCH=amd64 go build -o git-report-generator-linux
GOOS=windows GOARCH=amd64 go build -o git-report-generator.exe
GOOS=darwin GOARCH=amd64 go build -o git-report-generator-macos
```

## Troubleshooting

### Common Issues

1. **"Not a git repository"**
   - Ensure you're running the command in a Git repository or specify the correct path with `--repo`

2. **"User email not configured"**
   - Configure Git user email: `git config user.email "your.email@example.com"`
   - Or specify author explicitly: `--author your.email@example.com`

3. **"No commits found"**
   - Check the date range and author email
   - Verify the branch name is correct
   - Ensure commits exist in the specified time period

4. **Permission denied writing PDF**
   - Check write permissions for the output directory
   - Specify a different output path with `--output`

### Debug Mode

For verbose output, you can modify the source to add debug logging or use Git commands to verify data:

```bash
# Check commits in date range
git log --since="2024-01-01" --until="2024-01-31" --author="your.email@example.com" --oneline

# Check current branch
git branch --show-current

# Check user configuration
git config user.email
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Support

For support, please open an issue on the GitHub repository or contact the maintainer. 