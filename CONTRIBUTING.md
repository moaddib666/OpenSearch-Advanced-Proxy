# Contributing

We welcome contributions to `OpenSearchAdvancedProxy` and are grateful for every pull request and issue submitted. This
document outlines the process for contributing to the project.

## Development Workflow

We use the GitHub Flow for the development process. This means:

1. **Fork the repository** and create your branch from `main`.
2. **Clone your fork**: `git clone https://github.com/your-username/your-repository`
3. **Create a branch**: `git checkout -b my-branch-name`
4. **Make your changes**: Add or edit files as necessary.
5. **Commit your changes**: Follow the [Conventional Commits](#conventional-commits) guidelines.
6. **Push to your fork** and [submit a pull request](https://github.com/your-username/your-repository/pulls).
7. **Wait for review** and make any necessary [revisions](#revision-process).
8. **Your changes get merged!** 🎉

## Conventional Commits

We use Conventional Commits to make the commit history easier to understand and navigate. Each commit message should be
structured as follows:

```plaintext
type(scope): description

BREAKING CHANGE: breaking change description if applicable
```

### Types

- **feat**: A new feature
- **fix**: A bug fix
- **docs**: Documentation only changes
- **style**: Changes that do not affect the meaning of the code
- **refactor**: A code change that neither fixes a bug nor adds a feature
- **perf**: A code change that improves performance
- **test**: Adding missing tests or correcting existing tests
- **chore**: Changes to the build process or auxiliary tools and libraries

### Scope

The scope should be the name of the module or aspect of the project affected by the change.

### Description

The description should be a concise explanation of the changes. Use the imperative, present tense: "change" not "
changed" nor "changes".

## Linting

We use `golangci-lint` to ensure code quality and consistency. Please run `golangci-lint run` before submitting a pull
request to check for any issues.

## Writing Tests

We aim to maintain high test coverage, so new code should be accompanied by corresponding tests. Ensure your code passes
all existing tests and any new tests you have written:

```bash
go test ./...
```

## Revision Process

After you submit your pull request, it may be reviewed by maintainers or other contributors. They may suggest changes,
improvements, or ask questions. Be responsive to feedback to get your contribution merged more quickly.

## Questions or Problems?

If you have any questions or encounter any problems,
please [submit an issue](https://github.com/moaddib666/OpenSearch-Advanced-Proxy/issues). We're here to help!
