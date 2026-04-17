# gh-commit

A GitHub CLI extension that automatically generates Conventional Commits messages using GitHub Models API and commits them. It uses free [GitHub Models](https://docs.github.com/en/github-models) for inference, so you don't need to do any token setup - your existing GitHub CLI token will work!

## Installation

```bash
gh extension install ziishaned/gh-commit
gh commit
```

## Prerequisites

- [GitHub CLI](https://cli.github.com/) installed and authenticated
- Git repository with changes to commit

## Usage

### Basic Usage

Automatically stage changes (if needed), generate a commit message, and commit:

```bash
gh commit
```

### Advanced Options

```bash
# Use a different model
gh commit --model openai/gpt-4o-mini

# Dry run to see what would be committed
gh commit --dry-run

# Combine options
gh commit --model xai/grok-3-mini --dry-run
```

## How It Works

1. **Check for staged changes**: If no changes are staged, automatically stages all changes with `git add .`
2. **Analyze diff**: Captures the full diff of staged changes
3. **Generate message**: Sends diff to GitHub Models API with Conventional Commits prompt
4. **Commit**: Creates commit with generated message
5. **Error handling**: If anything fails, falls back to interactive mode

## Conventional Commits

The extension generates commit messages following the Conventional Commits specification:

```
type(scope): description
```

**Types:** `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`, `perf`, `ci`, `build`, `revert`

**Examples:**
- `feat(auth): add JWT token refresh mechanism`
- `fix(api): resolve memory leak in request handler`
- `docs(readme): update installation instructions`
- `refactor(database): extract query builder to separate module`

## Available Models

The extension supports any model available through GitHub Models API. Default is `openai/gpt-4o`.

Popular options:
- `openai/gpt-4o` (default)
- `openai/gpt-4o-mini` (faster, cheaper)
- `xai/grok-3-mini` (good reasoning)

## Contributing

Contributions are welcome! The extension is built with Go and follows the structure of [gh-standup](https://github.com/sgoedecke/gh-standup).

In particular, tweaking the [prompt](internal/llm/commit.prompt.yml) can significantly improve commit message quality.

## License

MIT
