package brew

import (
	"context"
	"testing"

	"github.com/frostyard/pm/internal/types"
)

// Integration test for Search with fixture data
func TestBackend_Search_Integration(t *testing.T) {
	t.Run("Empty query returns empty result", func(t *testing.T) {
		b := New(nil, nil, nil)
		ctx := context.Background()

		results, err := b.Search(ctx, "", types.SearchOptions{})
		if err != nil {
			t.Fatalf("Expected no error for empty query, got %v", err)
		}
		if len(results) != 0 {
			t.Errorf("Expected empty results for empty query, got %d results", len(results))
		}
	})

	// Note: Real API tests would require network access or a test server
	// For deterministic tests, we'd need to make the API URL injectable
	t.Run("Real API search", func(t *testing.T) {
		// Skip unless we're running integration tests
		if testing.Short() {
			t.Skip("Skipping integration test in short mode")
		}

		b := New(nil, nil, nil)
		ctx := context.Background()

		// This will hit the real API
		results, err := b.Search(ctx, "git", types.SearchOptions{})
		if err != nil {
			// It's okay if the API is unreachable in tests
			if !types.IsExternalFailure(err) && !types.IsNotAvailable(err) {
				t.Errorf("Expected ExternalFailure or NotAvailable, got %v", err)
			}
			t.Skipf("API unreachable: %v", err)
		}

		// Verify we got some results
		if len(results) == 0 {
			t.Error("Expected some results for 'git' query")
		}

		// Verify result structure
		for _, pkg := range results {
			if pkg.Name == "" {
				t.Error("Expected non-empty package name")
			}
			if pkg.Kind != "formula" {
				t.Errorf("Expected Kind='formula', got %q", pkg.Kind)
			}
		}
	})
}
