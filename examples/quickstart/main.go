package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/kazz187/claude-code-sdk-go/pkg/claudecode"
)

func main() {
	// Create a context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle interrupts gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\nInterrupted, shutting down...")
		cancel()
	}()

	// Example 1: Simple query
	fmt.Println("=== Example 1: Simple Query ===")
	if err := simpleQuery(ctx); err != nil {
		log.Printf("Error in simple query: %v\n", err)
	}

	fmt.Println("\n=== Example 2: Query with Options ===")
	if err := queryWithOptions(ctx); err != nil {
		log.Printf("Error in query with options: %v\n", err)
	}

	fmt.Println("\n=== Example 3: Using Tools ===")
	if err := queryWithTools(ctx); err != nil {
		log.Printf("Error in query with tools: %v\n", err)
	}
}

// simpleQuery demonstrates a basic query to Claude
func simpleQuery(ctx context.Context) error {
	// Query Claude
	messages, err := claudecode.Query(ctx, "What is 2 + 2?", nil)
	if err != nil {
		return fmt.Errorf("failed to query: %w", err)
	}

	// Process messages
	for msg := range messages {
		if err := processMessage(msg); err != nil {
			return err
		}
	}

	return nil
}

// queryWithOptions demonstrates using options
func queryWithOptions(ctx context.Context) error {
	// Create options
	options := claudecode.NewClaudeCodeOptions()
	systemPrompt := "You are a helpful assistant that responds concisely."
	options.SystemPrompt = &systemPrompt
	maxTurns := 1
	options.MaxTurns = &maxTurns

	// Query with options
	messages, err := claudecode.Query(ctx, "Tell me a short joke", options)
	if err != nil {
		return fmt.Errorf("failed to query: %w", err)
	}

	// Process messages
	for msg := range messages {
		if err := processMessage(msg); err != nil {
			return err
		}
	}

	return nil
}

// queryWithTools demonstrates using tools
func queryWithTools(ctx context.Context) error {
	// Create options with tools enabled
	options := claudecode.NewClaudeCodeOptions()
	options.AllowedTools = []string{"Read", "Write", "Bash"}
	options.PermissionMode = &[]claudecode.PermissionMode{claudecode.PermissionModeAcceptEdits}[0]

	// Create a temporary directory for the example
	tempDir, err := os.MkdirTemp("", "claude-sdk-example")
	if err != nil {
		return fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tempDir) // Clean up

	options.Cwd = &tempDir

	// Query that uses tools
	prompt := "Create a file called hello.txt with the content 'Hello from Claude SDK!'"
	messages, err := claudecode.Query(ctx, prompt, options)
	if err != nil {
		return fmt.Errorf("failed to query: %w", err)
	}

	// Process messages
	for msg := range messages {
		if err := processMessage(msg); err != nil {
			return err
		}
	}

	// Verify the file was created
	content, err := os.ReadFile(tempDir + "/hello.txt")
	if err == nil {
		fmt.Printf("\nFile created successfully! Content: %s\n", string(content))
	}

	return nil
}

// processMessage handles different message types
func processMessage(msg claudecode.Message) error {
	switch m := msg.(type) {
	case claudecode.UserMessage:
		fmt.Printf("[User]: %s\n", m.Content)

	case claudecode.AssistantMessage:
		fmt.Print("[Assistant]: ")
		for _, block := range m.Content {
			switch b := block.(type) {
			case claudecode.TextBlock:
				fmt.Print(b.Text)
			case claudecode.ToolUseBlock:
				fmt.Printf("\n[Tool Use - %s]: %v\n", b.Name, b.Input)
			case claudecode.ToolResultBlock:
				fmt.Printf("[Tool Result]: %v\n", b.Content)
			}
		}
		fmt.Println()

	case claudecode.SystemMessage:
		fmt.Printf("[System - %s]: %v\n", m.Subtype, m.Data)

	case claudecode.ResultMessage:
		fmt.Printf("\n[Result]: Session ID: %s, Duration: %dms, Cost: $%.4f\n",
			m.SessionID, m.DurationMs, *m.TotalCostUSD)
		if m.IsError {
			return fmt.Errorf("query ended with error")
		}

	default:
		fmt.Printf("[Unknown Message Type]: %T\n", msg)
	}

	return nil
}
