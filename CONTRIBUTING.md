# Contributing to k8s-yaml-splitter

Thank you for your interest in contributing! This document provides guidelines for contributing to the project.

## Development Setup

1. **Prerequisites**
   - Go 1.25.7 or later
   - Make
   - Git

2. **Clone and build**
   ```sh
   git clone https://github.com/ohauer/k8s-yaml-splitter.git
   cd k8s-yaml-splitter
   make build
   ```

3. **Run tests**
   ```sh
   make test
   ```

## Making Changes

1. **Fork the repository** on GitHub
2. **Create a feature branch** from `main`
3. **Make your changes** with clear, focused commits
4. **Add tests** for new functionality
5. **Run the test suite** to ensure nothing breaks
6. **Update documentation** if needed

## Code Style

- Follow standard Go formatting (`gofmt`)
- Write clear, descriptive commit messages
- Keep functions focused and well-documented
- Maintain backward compatibility when possible

## Testing

- All new features must include tests
- Run `make test` before submitting
- Test with real Kubernetes manifests when possible

## Submitting Changes

1. **Push your branch** to your fork
2. **Create a Pull Request** with:
   - Clear description of changes
   - Reference to any related issues
   - Test results

## Reporting Issues

- Use GitHub Issues for bug reports and feature requests
- Include reproduction steps for bugs
- Provide example YAML files when relevant

## Release Process

Releases are handled by maintainers using the automated release script.
