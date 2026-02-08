# Contributing to Drift Detector

Thank you for your interest in contributing to Drift Detector! 🎉

## Code of Conduct

By participating in this project, you agree to be respectful and professional in all interactions.

## How to Contribute

### Reporting Bugs

1. Check if the bug has already been reported in [Issues](https://github.com/MeowTux/drift-detector/issues)
2. If not, create a new issue with:
   - Clear title and description
   - Steps to reproduce
   - Expected vs actual behavior
   - Your environment (OS, Go version, cloud provider)
   - Relevant logs or error messages

### Suggesting Features

1. Check existing [Issues](https://github.com/MeowTux/drift-detector/issues) and [Discussions](https://github.com/MeowTux/drift-detector/discussions)
2. Create a new issue with:
   - Clear description of the feature
   - Use case and benefits
   - Potential implementation approach (optional)

### Contributing Code

#### Prerequisites

- Go 1.21 or higher
- Git
- Basic understanding of Terraform and cloud providers

#### Development Setup

1. Fork the repository
2. Clone your fork:
   ```bash
   git clone https://github.com/YOUR_USERNAME/drift-detector.git
   cd drift-detector
   ```

3. Add upstream remote:
   ```bash
   git remote add upstream https://github.com/MeowTux/drift-detector.git
   ```

4. Install dependencies:
   ```bash
   go mod download
   ```

5. Create a branch:
   ```bash
   git checkout -b feature/your-feature-name
   ```

#### Making Changes

1. Write clean, documented code
2. Follow Go best practices and conventions
3. Add tests for new functionality
4. Update documentation as needed
5. Run tests locally:
   ```bash
   make test
   ```

6. Format code:
   ```bash
   make fmt
   ```

#### Commit Guidelines

Use conventional commits:

- `feat:` New feature
- `fix:` Bug fix
- `docs:` Documentation changes
- `test:` Test additions or changes
- `refactor:` Code refactoring
- `chore:` Maintenance tasks

Example:
```
feat: add Azure VM drift detection
fix: correct S3 encryption check
docs: update installation instructions
```

#### Submitting Pull Requests

1. Push your changes:
   ```bash
   git push origin feature/your-feature-name
   ```

2. Create a Pull Request on GitHub
3. Fill out the PR template with:
   - Description of changes
   - Related issues
   - Testing performed
   - Screenshots (if applicable)

4. Wait for review and address feedback

## Project Structure

```
drift-detector/
├── cmd/              # CLI commands
├── internal/         # Internal packages
│   ├── detectors/    # Cloud provider detectors
│   ├── notifiers/    # Notification implementations
│   ├── terraform/    # Terraform state handling
│   └── drift/        # Drift analysis logic
├── pkg/              # Public packages
├── examples/         # Example configurations
└── docs/             # Documentation
```

## Adding a New Cloud Provider

1. Create detector in `internal/detectors/`
2. Implement `Detector` interface
3. Add configuration options
4. Write tests
5. Update documentation

## Testing

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Run specific test
go test -v ./internal/detectors -run TestAWSDetector
```

## Documentation

- Keep README.md up to date
- Add inline code comments
- Update examples when changing features
- Document breaking changes

## Questions?

- Open a [Discussion](https://github.com/MeowTux/drift-detector/discussions)
- Check existing [Issues](https://github.com/MeowTux/drift-detector/issues)
- Email: meowtux@example.com

Thank you for contributing! 🚀
