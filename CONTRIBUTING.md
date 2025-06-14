# Contributing to Quiver

Thank you for your interest in contributing to Quiver! This document provides guidelines and information for contributors.

## Table of Contents

- [Contributing to Quiver](#contributing-to-quiver)
	- [Table of Contents](#table-of-contents)
	- [Code of Conduct](#code-of-conduct)
	- [Getting Started](#getting-started)
		- [Prerequisites](#prerequisites)
		- [Development Setup](#development-setup)
	- [Project Structure](#project-structure)
	- [Contributing Guidelines](#contributing-guidelines)
		- [Code Style](#code-style)
		- [Testing](#testing)
		- [Commit Messages](#commit-messages)
	- [Pull Request Process](#pull-request-process)
		- [PR Requirements](#pr-requirements)
	- [Issue Reporting](#issue-reporting)
	- [Development Workflow](#development-workflow)
		- [Adding New Features](#adding-new-features)
		- [Package Development](#package-development)
		- [Testing Guidelines](#testing-guidelines)
		- [Documentation](#documentation)
	- [Getting Help](#getting-help)
	- [Recognition](#recognition)

## Code of Conduct

By participating in this project, you agree to abide by our Code of Conduct. Please be respectful and professional in all interactions.

## Getting Started

### Prerequisites

- Go 1.24.2 or later
- Git
- Basic understanding of Go programming language
- Familiarity with REST APIs and HTTP servers

### Development Setup

1. **Fork the repository**
   ```bash
   git clone https://github.com/your-username/quiver.git
   cd quiver
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Build the project**
   ```bash
   go build -o bin/quiver ./cmd/quiver
   ```

4. **Run tests**
   ```bash
   go test ./...
   ```

5. **Run the application**
   ```bash
   ./bin/quiver
   ```

## Project Structure

```
quiver/
├── cmd/                    # Application entry points
│   └── quiver/            # Main application
├── internal/              # Private application code
│   ├── config/            # Configuration management
│   ├── logger/            # Logging utilities
│   ├── server/            # HTTP server implementation
│   └── ui/               # User interface components
├── pkg/                   # Public library code
│   └── packages/          # Package management
├── api/                   # API definitions and documentation
├── docs/                  # Documentation
├── scripts/               # Build and deployment scripts
├── test/                  # Test files and test data
└── examples/              # Example configurations and packages
```

## Contributing Guidelines

### Code Style

- Follow Go best practices and conventions
- Use `gofmt` to format your code
- Use `golint` to check for style issues
- Write clear, documented code with meaningful variable names
- Add comments for complex logic

### Testing

- Write unit tests for new functionality
- Ensure all tests pass before submitting a PR
- Maintain or improve code coverage
- Use table-driven tests where appropriate

### Commit Messages

Use conventional commit format:

```
type(scope): description

[optional body]

[optional footer]
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes
- `refactor`: Code refactoring
- `test`: Test additions or modifications
- `chore`: Maintenance tasks

Examples:
```
feat(packages): add package validation
fix(server): resolve memory leak in HTTP handler
docs(api): update REST API documentation
```

## Pull Request Process

1. **Create a feature branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes**
   - Follow the coding guidelines
   - Add tests for new functionality
   - Update documentation as needed

3. **Test your changes**
   ```bash
   go test ./...
   go vet ./...
   ```

4. **Commit your changes**
   ```bash
   git add .
   git commit -m "feat: add your feature description"
   ```

5. **Push to your fork**
   ```bash
   git push origin feature/your-feature-name
   ```

6. **Open a Pull Request**
   - Provide a clear description of the changes
   - Reference any related issues
   - Include screenshots if applicable

### PR Requirements

- [ ] Code follows project style guidelines
- [ ] Tests are added for new functionality
- [ ] All tests pass
- [ ] Documentation is updated
- [ ] Commit messages follow conventional format
- [ ] PR description clearly explains the changes

## Issue Reporting

When reporting issues, please include:

1. **Bug Reports**
   - Go version
   - Operating system
   - Steps to reproduce
   - Expected vs actual behavior
   - Error messages or logs

2. **Feature Requests**
   - Clear description of the feature
   - Use case and benefits
   - Possible implementation approach

## Development Workflow

### Adding New Features

1. **Planning**
   - Discuss major features in an issue first
   - Consider backward compatibility
   - Plan for testing and documentation

2. **Implementation**
   - Create feature branch
   - Implement in small, logical commits
   - Add comprehensive tests
   - Update documentation

3. **Review**
   - Self-review your code
   - Ensure all tests pass
   - Submit pull request

### Package Development

For creating new game server packages:

1. **Package Structure**
   ```
   packages/your-game/
   ├── manifest.yaml      # Package manifest
   ├── README.md         # Package documentation
   ├── scripts/          # Setup and management scripts
   └── config/           # Default configurations
   ```

2. **Manifest Format**
   ```yaml
   version: "0.1"
   name: "Your Game Server"
   description: "Description of your game server"
   author: "Your Name"
   type: "your-game"
   commands:
     start: "./scripts/start.sh"
     stop: "./scripts/stop.sh"
     build: "./scripts/build.sh"
   config:
     port: 27015
     max_players: 32
   ```

### Testing Guidelines

- Unit tests for individual functions
- Integration tests for component interactions
- End-to-end tests for complete workflows
- Use mock objects for external dependencies
- Test error conditions and edge cases

### Documentation

- Update API documentation for new endpoints
- Add examples for new features
- Keep README up to date
- Document configuration options
- Provide migration guides for breaking changes

## Getting Help

- Join our community discussions
- Check existing issues and documentation
- Ask questions in issue comments
- Follow us on social media for updates

## Recognition

Contributors will be recognized in:
- CONTRIBUTORS.md file
- Release notes for significant contributions
- Project documentation

Thank you for contributing to Quiver! 🏹 