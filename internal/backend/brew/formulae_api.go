package brew

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/frostyard/pm/internal/types"
)

const (
	formulaeAPIBase = "https://formulae.brew.sh/api"
)

// searchFormulae searches for formulae by name using the API.
// Returns a list of matching package references.
func (b *Backend) searchFormulae(ctx context.Context, query string) ([]types.PackageRef, error) {
	// The Formulae API provides /api/formula.json which lists all formulae
	// We fetch it and filter client-side
	url := formulaeAPIBase + "/formula.json"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, &types.ExternalFailureError{
			Operation: types.OperationSearch,
			Backend:   "brew",
			Err:       fmt.Errorf("failed to create request: %w", err),
		}
	}

	resp, err := b.httpClient.Do(req)
	if err != nil {
		return nil, &types.ExternalFailureError{
			Operation: types.OperationSearch,
			Backend:   "brew",
			Err:       fmt.Errorf("failed to fetch formula list: %w", err),
		}
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, &types.ExternalFailureError{
			Operation: types.OperationSearch,
			Backend:   "brew",
			Err:       fmt.Errorf("API returned status %d", resp.StatusCode),
		}
	}

	// The API returns an array of formula names
	var names []string
	if err := json.NewDecoder(resp.Body).Decode(&names); err != nil {
		return nil, &types.ExternalFailureError{
			Operation: types.OperationSearch,
			Backend:   "brew",
			Err:       fmt.Errorf("failed to parse response: %w", err),
		}
	}

	// Filter formulae by query (case-insensitive substring match)
	var results []types.PackageRef
	queryLower := strings.ToLower(query)
	for _, name := range names {
		if strings.Contains(strings.ToLower(name), queryLower) {
			results = append(results, types.PackageRef{
				Name: name,
				Kind: "formula",
			})
		}
	}

	return results, nil
}
