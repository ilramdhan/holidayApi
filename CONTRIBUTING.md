# Contributing to Holiday API Indonesia

Thank you for your interest in contributing to Holiday API Indonesia! We welcome contributions from the community.

## 🚀 Getting Started

1. **Fork the repository** on GitHub
2. **Clone your fork** locally
3. **Create a branch** for your feature or bug fix
4. **Make your changes**
5. **Submit a pull request**

## 📋 Development Setup

### Prerequisites
- Go 1.23 or higher
- Git
- Docker (optional)

### Local Development

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/holidayapi.git
cd holidayapi

# Install dependencies
go mod tidy

# Run the application
go run cmd/server/main.go

# Run tests
go test ./...
```

## 📝 Code Style

- Follow standard Go formatting (`go fmt`)
- Use meaningful variable and function names
- Add comments for exported functions
- Keep functions focused and small
- Write tests for new features

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test ./internal/services/...
```

## 📦 Submitting Changes

1. **Create a feature branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes**
   - Write clear, concise commit messages
   - Follow existing code style
   - Add tests for new functionality

3. **Commit your changes**
   ```bash
   git add .
   git commit -m "feat: add new feature"
   ```

4. **Push to your fork**
   ```bash
   git push origin feature/your-feature-name
   ```

5. **Create a Pull Request**
   - Describe what changes you made and why
   - Reference any related issues
   - Ensure all tests pass

## 🏷️ Commit Message Convention

We follow conventional commits:

- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation changes
- `style:` - Code style changes (formatting, etc.)
- `refactor:` - Code refactoring
- `test:` - Test changes
- `chore:` - Build process or auxiliary tool changes

Examples:
```
feat: add support for filtering by holiday type
fix: resolve timezone issue in today endpoint
docs: update API documentation
```

## 🐛 Reporting Bugs

When reporting bugs, please include:

- **Description**: Clear description of the bug
- **Steps to Reproduce**: How to reproduce the issue
- **Expected Behavior**: What you expected to happen
- **Actual Behavior**: What actually happened
- **Environment**: OS, Go version, etc.
- **Screenshots**: If applicable

## 💡 Feature Requests

We welcome feature requests! Please:

- Check if the feature has already been requested
- Provide clear use case and rationale
- Describe the expected behavior

## 📄 License

By contributing, you agree that your contributions will be licensed under the MIT License.

## 🤝 Code of Conduct

- Be respectful and inclusive
- Welcome newcomers
- Focus on constructive feedback
- Respect different viewpoints

## 📞 Questions?

- Open an issue on GitHub
- Start a discussion
- Email: ilramdhan@gmail.com

Thank you for contributing! 🎉
