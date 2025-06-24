# Claude Code SDK for Go

Go SDK for Claude Code. See the [Claude Code SDK documentation](https://docs.anthropic.com/en/docs/claude-code/sdk) for more information.

## Installation

```bash
go get github.com/kazz187/claude-code-sdk-go/pkg/claudecode
```

**Prerequisites:**
- Go 1.21+
- Node.js
- Claude Code: `npm install -g @anthropic-ai/claude-code`

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/kazz187/claude-code-sdk-go/pkg/claudecode"
)

func main() {
    ctx := context.Background()
    
    // Query Claude
    messages, err := claudecode.Query(ctx, "What is 2 + 2?", nil)
    if err != nil {
        log.Fatal(err)
    }
    
    // Process messages
    for msg := range messages {
        if assistantMsg, ok := msg.(claudecode.AssistantMessage); ok {
            for _, block := range assistantMsg.Content {
                if textBlock, ok := block.(claudecode.TextBlock); ok {
                    fmt.Println(textBlock.Text)
                }
            }
        }
    }
}
```

## Usage

### Basic Query

```go
import (
    "context"
    "github.com/kazz187/claude-code-sdk-go/pkg/claudecode"
)

// Simple query
messages, err := claudecode.Query(ctx, "Hello Claude", nil)
if err != nil {
    return err
}

for msg := range messages {
    switch m := msg.(type) {
    case claudecode.AssistantMessage:
        for _, block := range m.Content {
            if textBlock, ok := block.(claudecode.TextBlock); ok {
                fmt.Println(textBlock.Text)
            }
        }
    }
}

// With options
options := claudecode.NewClaudeCodeOptions()
systemPrompt := "You are a helpful assistant"
options.SystemPrompt = &systemPrompt
maxTurns := 1
options.MaxTurns = &maxTurns

messages, err = claudecode.Query(ctx, "Tell me a joke", options)
```

### Using Tools

```go
options := claudecode.NewClaudeCodeOptions()
options.AllowedTools = []string{"Read", "Write", "Bash"}
permMode := claudecode.PermissionModeAcceptEdits
options.PermissionMode = &permMode

messages, err := claudecode.Query(ctx, "Create a hello.py file", options)
if err != nil {
    return err
}

// Process tool use and results
for msg := range messages {
    // Handle messages
}
```

### Working Directory

```go
options := claudecode.NewClaudeCodeOptions()
workDir := "/path/to/project"
options.Cwd = &workDir
```

## API Reference

### `Query(ctx, prompt, options)`

Main function for querying Claude.

**Parameters:**
- `ctx` (context.Context): Context for cancellation and timeouts
- `prompt` (string): The prompt to send to Claude
- `options` (*ClaudeCodeOptions): Optional configuration

**Returns:** 
- `<-chan Message`: Channel of response messages
- `error`: Error if connection fails

### Types

See [pkg/claudecode/types.go](pkg/claudecode/types.go) for complete type definitions:
- `ClaudeCodeOptions` - Configuration options
- `AssistantMessage`, `UserMessage`, `SystemMessage`, `ResultMessage` - Message types
- `TextBlock`, `ToolUseBlock`, `ToolResultBlock` - Content blocks

## Error Handling

```go
import "github.com/kazz187/claude-code-sdk-go/pkg/claudecode"

messages, err := claudecode.Query(ctx, "Hello", nil)
if err != nil {
    switch e := err.(type) {
    case *claudecode.CLINotFoundError:
        fmt.Println("Please install Claude Code")
    case *claudecode.ProcessError:
        fmt.Printf("Process failed with exit code: %v\n", e.ExitCode)
    case *claudecode.CLIJSONDecodeError:
        fmt.Printf("Failed to parse response: %v\n", e)
    default:
        fmt.Printf("Unknown error: %v\n", err)
    }
}
```

See [pkg/claudecode/errors.go](pkg/claudecode/errors.go) for all error types.

## Available Tools

See the [Claude Code documentation](https://docs.anthropic.com/en/docs/claude-code/security#tools-available-to-claude) for a complete list of available tools.

## Examples

See [examples/quickstart/main.go](examples/quickstart/main.go) for a complete working example.

## Development

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Format code
go fmt ./...

# Build
go build ./pkg/claudecode
```

## License

MIT