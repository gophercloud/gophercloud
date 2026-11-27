package documentation_test

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"testing"
	"time"
)

// TestFAQFragments compiles the partial examples with their documented inputs.
// It also runs the JSON decoding example without making any network requests.
func TestFAQFragments(t *testing.T) {
	contents, err := os.ReadFile(filepath.Join("..", "..", "docs", "FAQ.md"))
	if err != nil {
		t.Fatal(err)
	}
	blocks := regexp.MustCompile("(?s)```go\\n(.*?)```").FindAllSubmatch(contents, -1)
	if len(blocks) != 6 {
		t.Fatalf("expected 6 FAQ code blocks, got %d", len(blocks))
	}
	sources := map[string]string{
		"request": fmt.Sprintf(`package main
import (
 "context"
 "github.com/gophercloud/gophercloud/v2"
 "github.com/gophercloud/gophercloud/v2/openstack/compute/v2/servers"
)
%s
%s
func exampleCustomRequest() {
 var computeClient *gophercloud.ServiceClient
 %s
 _, _ = server, err
}
`, blocks[1][1], blocks[2][1], blocks[3][1]),
		"response": fmt.Sprintf(`package main
import (
 "context"
 "github.com/gophercloud/gophercloud/v2"
 "github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/volumes"
)
func exampleCustomResponse() {
 var client *gophercloud.ServiceClient
 volID := "volume-id"
 %s
 _ = err
}
`, blocks[4][1]),
		"unmarshal": fmt.Sprintf(`package main
import (
 "encoding/json"
 "testing"
 "time"
 "github.com/gophercloud/gophercloud/v2"
)
%s
func TestUnmarshal(t *testing.T) {
 var volume MyVolume
 if err := json.Unmarshal([]byte("{\"id\":\"volume-id\",\"created_at\":\"2024-01-02T03:04:05.000\"}"), &volume); err != nil { t.Fatal(err) }
 expected := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)
 if volume.ID != "volume-id" || !volume.TimeCreated.Equal(expected) {
 t.Fatalf("unexpected decoded volume: %%+v", volume)
 }
}
`, blocks[5][1]),
	}
	for name, source := range sources {
		t.Run(name, func(t *testing.T) {
			filename := filepath.Join(t.TempDir(), "example_test.go")
			if err := os.WriteFile(filename, []byte(source), 0600); err != nil {
				t.Fatal(err)
			}
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
			defer cancel()
			cmd := exec.CommandContext(ctx, "go", "test", filename)
			if output, err := cmd.CombinedOutput(); err != nil {
				t.Fatalf("FAQ example failed: %v\n%s", err, output)
			}
		})
	}
}
