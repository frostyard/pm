package snap

import (
	"context"
	"net/http"
	"strings"

	"github.com/frostyard/pm/internal/runner"
	"github.com/frostyard/pm/internal/types"
)

// Backend implements the snap backend.
type Backend struct {
	httpClient *http.Client
	runner     runner.Runner
	progress   types.ProgressReporter
}

// New creates a new snap backend.
func New(httpClient *http.Client, r runner.Runner, progress types.ProgressReporter) *Backend {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &Backend{
		httpClient: httpClient,
		runner:     r,
		progress:   progress,
	}
}

// Available checks if snapd is available by querying /v2/system-info.
func (b *Backend) Available(ctx context.Context) (bool, error) {
	// Try to reach the snapd API
	// Note: In production, this would use a unix socket transport
	// For now, we test if the http client is functional
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost/v2/system-info", nil)
	if err != nil {
		return false, &types.NotAvailableError{Backend: "snap", Reason: "failed to create request: " + err.Error()}
	}

	resp, err := b.httpClient.Do(req)
	if err != nil {
		return false, &types.NotAvailableError{Backend: "snap", Reason: "failed to reach snapd API: " + err.Error()}
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return true, nil
	}

	return false, &types.NotAvailableError{Backend: "snap", Reason: "snapd API returned non-2xx status"}
}

// Capabilities returns snap capabilities.
func (b *Backend) Capabilities(ctx context.Context) ([]types.Capability, error) {
	// Snap backend supports operations when runner is available
	hasRunner := b.runner != nil
	return []types.Capability{
		{Operation: types.OperationSearch, Supported: hasRunner, Notes: "via snap find CLI"},
		{Operation: types.OperationUpdateMetadata, Supported: hasRunner, Notes: "via snap refresh CLI"},
		{Operation: types.OperationUpgradePackages, Supported: hasRunner, Notes: "via snap refresh CLI"},
		{Operation: types.OperationInstall, Supported: hasRunner, Notes: "via snap install CLI"},
		{Operation: types.OperationUninstall, Supported: hasRunner, Notes: "via snap remove CLI"},
		{Operation: types.OperationListInstalled, Supported: hasRunner, Notes: "via snap list CLI"},
	}, nil
}

// Update implements Updater using `snap refresh --list`.
func (b *Backend) Update(ctx context.Context, opts types.UpdateOptions) (types.UpdateResult, error) {
	if b.runner == nil {
		return types.UpdateResult{}, types.ErrNotSupported
	}

	helper := types.NewProgressHelper(opts.Progress)
	helper.BeginAction("Update")
	defer helper.EndAction()

	helper.BeginTask("Checking for snap updates")
	stdout, _, err := runner.RunWithExternalError(
		ctx,
		b.runner,
		types.OperationUpdateMetadata,
		"snap",
		"snap",
		"refresh",
		"--list",
	)
	helper.EndTask()

	if err != nil {
		helper.Error("Update check failed: " + err.Error())
		return types.UpdateResult{}, err
	}

	// Check if there are updates available
	changed := len(strings.TrimSpace(stdout)) > 0 && !strings.Contains(stdout, "All snaps up to date")

	helper.Info("Update check completed")
	return types.UpdateResult{Changed: changed}, nil
}

// Upgrade implements Upgrader using `snap refresh`.
func (b *Backend) Upgrade(ctx context.Context, opts types.UpgradeOptions) (types.UpgradeResult, error) {
	if b.runner == nil {
		return types.UpgradeResult{}, types.ErrNotSupported
	}

	helper := types.NewProgressHelper(opts.Progress)
	helper.BeginAction("Upgrade")
	defer helper.EndAction()

	helper.BeginTask("Running snap refresh")
	stdout, _, err := runner.RunWithExternalError(
		ctx,
		b.runner,
		types.OperationUpgradePackages,
		"snap",
		"snap",
		"refresh",
	)
	helper.EndTask()

	if err != nil {
		helper.Error("Upgrade failed: " + err.Error())
		return types.UpgradeResult{}, err
	}

	// Parse upgraded snaps from output
	var packagesChanged []types.PackageRef
	changed := false

	// Look for lines indicating refreshes
	lines := strings.Split(stdout, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Snap shows refreshed snaps like "<snap-name> <version> from <publisher> refreshed"
		if strings.Contains(line, "refreshed") || strings.Contains(line, "installed") {
			changed = true
			// Extract snap name (first field)
			fields := strings.Fields(line)
			if len(fields) >= 1 {
				snapName := fields[0]
				packagesChanged = append(packagesChanged, types.PackageRef{
					Name: snapName,
					Kind: "snap",
				})
			}
		}
	}

	// Also check for "All snaps up to date" message
	if strings.Contains(stdout, "All snaps up to date") {
		changed = false
	}

	if changed {
		helper.Info("Upgrade completed: upgraded snaps")
	} else {
		helper.Info("Upgrade completed: no snaps needed upgrading")
	}

	return types.UpgradeResult{
		Changed:         changed,
		PackagesChanged: packagesChanged,
	}, nil
}

// Install implements Installer using `snap install`.
func (b *Backend) Install(ctx context.Context, pkgs []types.PackageRef, opts types.InstallOptions) (types.InstallResult, error) {
	if b.runner == nil {
		return types.InstallResult{}, types.ErrNotSupported
	}

	if len(pkgs) == 0 {
		return types.InstallResult{}, nil
	}

	helper := types.NewProgressHelper(opts.Progress)
	helper.BeginAction("Install")
	defer helper.EndAction()

	// Build package list
	pkgNames := make([]string, 0, len(pkgs)+1)
	pkgNames = append(pkgNames, "install")
	for _, pkg := range pkgs {
		pkgNames = append(pkgNames, pkg.Name)
	}

	helper.BeginTask("Running snap install")
	stdout, _, err := runner.RunWithExternalError(
		ctx,
		b.runner,
		types.OperationInstall,
		"snap",
		"snap",
		pkgNames...,
	)
	helper.EndTask()

	if err != nil {
		helper.Error("Install failed: " + err.Error())
		return types.InstallResult{}, err
	}

	// Check if packages were installed
	var installed []types.PackageRef
	changed := false

	// Look for installation confirmations in output
	lines := strings.Split(stdout, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "installed") {
			changed = true
			// Try to extract snap name from the line
			for _, pkg := range pkgs {
				if strings.Contains(line, pkg.Name) {
					installed = append(installed, pkg)
					break
				}
			}
		}
	}

	// If we couldn't parse specific packages but the command succeeded, mark all as installed
	if changed && len(installed) == 0 {
		installed = pkgs
	}

	if changed {
		helper.Info("Install completed: installed snaps")
	} else {
		helper.Info("Install completed: snaps already installed")
	}

	return types.InstallResult{
		Changed:           changed,
		PackagesInstalled: installed,
	}, nil
}

// Uninstall implements Uninstaller using `snap remove`.
func (b *Backend) Uninstall(ctx context.Context, pkgs []types.PackageRef, opts types.UninstallOptions) (types.UninstallResult, error) {
	if b.runner == nil {
		return types.UninstallResult{}, types.ErrNotSupported
	}

	if len(pkgs) == 0 {
		return types.UninstallResult{}, nil
	}

	helper := types.NewProgressHelper(opts.Progress)
	helper.BeginAction("Uninstall")
	defer helper.EndAction()

	// Build package list
	pkgNames := make([]string, 0, len(pkgs)+1)
	pkgNames = append(pkgNames, "remove")
	for _, pkg := range pkgs {
		pkgNames = append(pkgNames, pkg.Name)
	}

	helper.BeginTask("Running snap remove")
	stdout, _, err := runner.RunWithExternalError(
		ctx,
		b.runner,
		types.OperationUninstall,
		"snap",
		"snap",
		pkgNames...,
	)
	helper.EndTask()

	if err != nil {
		helper.Error("Uninstall failed: " + err.Error())
		return types.UninstallResult{}, err
	}

	// Check if packages were uninstalled
	var uninstalled []types.PackageRef
	changed := false

	// Look for removal confirmations in output
	lines := strings.Split(stdout, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "removed") {
			changed = true
			// Try to extract snap name from the line
			for _, pkg := range pkgs {
				if strings.Contains(line, pkg.Name) {
					uninstalled = append(uninstalled, pkg)
					break
				}
			}
		}
	}

	// If we couldn't parse specific packages but the command succeeded, mark all as uninstalled
	if changed && len(uninstalled) == 0 {
		uninstalled = pkgs
	}

	if changed {
		helper.Info("Uninstall completed: removed snaps")
	} else {
		helper.Info("Uninstall completed: snaps were not installed")
	}

	return types.UninstallResult{
		Changed:             changed,
		PackagesUninstalled: uninstalled,
	}, nil
}

// Search implements Searcher using `snap find`.
func (b *Backend) Search(ctx context.Context, query string, opts types.SearchOptions) ([]types.PackageRef, error) {
	if b.runner == nil {
		return nil, types.ErrNotSupported
	}

	if query == "" {
		return []types.PackageRef{}, nil
	}

	helper := types.NewProgressHelper(opts.Progress)
	helper.BeginAction("Search")
	defer helper.EndAction()

	helper.BeginTask("Running snap find")
	stdout, _, err := runner.RunWithExternalError(
		ctx,
		b.runner,
		types.OperationSearch,
		"snap",
		"snap",
		"find",
		query,
	)
	helper.EndTask()

	if err != nil {
		helper.Error("Search failed: " + err.Error())
		return nil, err
	}

	// Parse search results
	// Snap find output format:
	// Name       Version    Publisher    Notes  Summary
	// firefox    123.0      mozilla✓     -      Mozilla Firefox web browser
	var results []types.PackageRef
	lines := strings.Split(stdout, "\n")

	// Skip header line
	for i, line := range lines {
		if i == 0 {
			continue // Skip header
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Parse fields - split by whitespace
		fields := strings.Fields(line)
		if len(fields) >= 1 {
			snapName := fields[0]

			results = append(results, types.PackageRef{
				Name: snapName,
				Kind: "snap",
			})
		}
	}

	helper.Info("Search completed")
	return results, nil
}

// ListInstalled implements Lister using `snap list`.
func (b *Backend) ListInstalled(ctx context.Context, opts types.ListOptions) ([]types.InstalledPackage, error) {
	if b.runner == nil {
		return nil, types.ErrNotSupported
	}

	helper := types.NewProgressHelper(opts.Progress)
	helper.BeginAction("ListInstalled")
	defer helper.EndAction()

	helper.BeginTask("Running snap list")
	stdout, _, err := runner.RunWithExternalError(
		ctx,
		b.runner,
		types.OperationListInstalled,
		"snap",
		"snap",
		"list",
	)
	helper.EndTask()

	if err != nil {
		helper.Error("ListInstalled failed: " + err.Error())
		return nil, err
	}

	// Parse output: columns are Name, Version, Rev, Tracking, Publisher, Notes
	var packages []types.InstalledPackage
	lines := strings.Split(stdout, "\n")

	// Skip header line
	for i, line := range lines {
		if i == 0 {
			continue // Skip header
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Split by whitespace
		fields := strings.Fields(line)
		if len(fields) >= 2 {
			snapName := fields[0]
			version := fields[1]

			packages = append(packages, types.InstalledPackage{
				Ref: types.PackageRef{
					Name: snapName,
					Kind: "snap",
				},
				Version: version,
			})
		}
	}

	helper.Info("ListInstalled completed")
	return packages, nil
}
