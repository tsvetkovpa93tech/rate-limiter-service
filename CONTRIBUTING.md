# Contributing to Rate Limiter Service

Thank you for your interest in contributing! This document provides guidelines and instructions for contributing to the project.

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/yourusername/rate-limiter-service.git`
3. Create a branch: `git checkout -b feature/your-feature-name`
4. Make your changes
5. Test your changes: `go test ./...`
6. Commit your changes: `git commit -am 'Add some feature'`
7. Push to your fork: `git push origin feature/your-feature-name`
8. Submit a pull request

## Code Style

- Follow Go conventions and best practices
- Use `gofmt` to format code
- Run `golint` and fix any issues
- Write clear, self-documenting code
- Add comments for exported functions and types

## Testing

- Write tests for all new features
- Ensure all tests pass: `go test ./...`
- Aim for high test coverage
- Include integration tests for new algorithms or storage backends

## Documentation

- Update README.md if adding new features
- Update API documentation in `api/openapi.yaml`
- Add code comments for complex logic
- Update examples if API changes

## Pull Request Process

1. Ensure your code follows the project's style guidelines
2. Update documentation as needed
3. Add tests for new functionality
4. Ensure all tests pass
5. Update CHANGELOG.md (if applicable)
6. Submit PR with clear description of changes

## Reporting Issues

When reporting issues, please include:
- Description of the issue
- Steps to reproduce
- Expected behavior
- Actual behavior
- Environment (OS, Go version, etc.)
- Relevant logs or error messages

## Feature Requests

For feature requests, please:
- Describe the feature and use case
- Explain why it would be useful
- Consider implementation approach
- Discuss in an issue before implementing

Thank you for contributing!

