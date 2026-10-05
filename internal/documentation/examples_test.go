package documentation_test

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"
)

// TestMarkdownExamples compiles complete programs without running their main
// functions, so no cloud configuration or live OpenStack service is needed.
func TestMarkdownExamples(t *testing.T) {
	blocks := regexp.MustCompile("(?s)```go\\n(.*?)```")
	for _, path := range []string{"README.md", "docs/FAQ.md"} {
		t.Run(path, func(t *testing.T) {
			contents, err := os.ReadFile(filepath.Join("..", "..", path))
			if err != nil {
				t.Fatal(err)
			}
			count := 0
			for _, match := range blocks.FindAllSubmatch(contents, -1) {
				source := string(match[1])
				if !strings.HasPrefix(source, "package main\n") {
					continue
				}
				count++
				t.Run(fmt.Sprintf("program%d", count), func(t *testing.T) {
					filename := filepath.Join(t.TempDir(), "main.go")
					if err := os.WriteFile(filename, []byte(source), 0600); err != nil {
						t.Fatal(err)
					}
					ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
					defer cancel()
					cmd := exec.CommandContext(ctx, "go", "test", "-run", "^$", filename)
					if output, err := cmd.CombinedOutput(); err != nil {
						t.Fatalf("example does not compile: %v\n%s", err, output)
					}
				})
			}
			if count == 0 {
				t.Fatal("no complete Go examples found")
			}
		})
	}
}
