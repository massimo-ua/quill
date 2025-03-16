# GitHub Document Store Provider

This package implements a document store provider that uses GitHub as a backend for storing documentation.

## Features

- Store documentation in a GitHub repository with proper directory structure
- Support for versioning via Git commit history
- Commit messages that include type, category, and timestamp information
- Automatic directory creation for structured documentation
- Proper error handling and context propagation

## Usage

### Configuration

Create a configuration with your GitHub credentials:

```go
config := &github.Config{
    Token:          "your-github-token",
    Owner:          "your-username-or-org",
    Repo:           "your-docs-repo",
    Branch:         "main",                 // Optional, defaults to "main"
    BasePath:       "docs",                 // Optional, base directory in repo
    CommitterName:  "Quill Bot",
    CommitterEmail: "bot@example.com",
}
```

### Creating a Provider

Use the factory to create a provider:

```go
provider, err := github.NewGitHubDocumentStoreProvider(config)
if err != nil {
    // Handle error
}
```

### Injecting the Provider

Inject the provider into the documentation service:

```go
// Assuming you have an AI agent provider
aiAgent := yourAiAgentProvider

// Create the documentation service
docService := services.NewDocumentationService(provider, aiAgent)
```

### Example: Storing Documentation

```go
ctx := context.Background()

// Create documentation for a status update
err := docService.CreateDocumentation(
    ctx,
    domain.MessageTypeStatus,
    domain.CategoryDevelopment,
    "Weekly update on backend development progress",
    []*domain.Reference{
        // Any references
    },
)
```

## Directory Structure

Documentation is organized in the GitHub repository as follows:

```
/docs
  /development
    /decision-20220101-120000.md
    /status-20220102-150000.md
  /operations
    /information-20220105-113000.md
  /product
    /idea-20220110-090000.md
```

The structure is based on the message category and type, with timestamps included in filenames for proper versioning.

## Commit Messages

Commit messages are automatically generated based on the document's metadata:

- For new documents: "Add [type] documentation ([category])"
- For updates: "Update [type] documentation ([category]) at [timestamp]"
- For deletions: "Delete documentation [path]"

This provides clear versioning history in the GitHub repository.