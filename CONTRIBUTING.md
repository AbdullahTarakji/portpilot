# Contributing to portpilot

Thanks for wanting to contribute! 🎉

## Quick Start

```bash
# Clone the repo
git clone https://github.com/AbdullahTarakji/portpilot.git
cd portpilot

# Install dependencies
go mod download

# Build
go build -o portpilot ./cmd/portpilot

# Run tests
go test ./...

# Run linter
golangci-lint run
```

## Development Workflow

1. Fork the repo
2. Create a feature branch from `develop`: `git checkout -b feature/my-feature develop`
3. Make your changes
4. Write/update tests
5. Run `go test ./...` and `golangci-lint run`
6. Commit with conventional commits: `feat:`, `fix:`, `docs:`, `chore:`, `test:`, `refactor:`
7. Push and open a PR against `develop`

## Commit Messages

We use [Conventional Commits](https://www.conventionalcommits.org/):

- `feat: add port filtering` — new feature
- `fix: handle permission error on Linux` — bug fix
- `docs: update README install section` — docs only
- `test: add scanner unit tests` — tests only
- `refactor: extract port parser` — code restructuring
- `chore: update dependencies` — maintenance

## Code Style

- Run `gofmt` (enforced by CI)
- Follow [Effective Go](https://go.dev/doc/effective_go)
- Add comments for exported functions
- Keep functions focused and small

## Reporting Bugs

Use the [bug report template](https://github.com/AbdullahTarakji/portpilot/issues/new?template=bug_report.md).

## Suggesting Features

Use the [feature request template](https://github.com/AbdullahTarakji/portpilot/issues/new?template=feature_request.md).

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
