package brew

import (
	"context"
	"net/http"
	"strings"

	"github.com/frostyard/pm/internal/runner"
	"github.com/frostyard/pm/internal/types"
)

// Backend implements the brew backend.
type Backend struct {
	httpClient *http.Client
	runner     runner.Runner
	progress   types.ProgressReporter
}

// New creates a new brew backend.
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

// Available checks if brew is available by testing the Formulae API endpoint.
func (b *Backend) Available(ctx context.Context) (bool, error) {
	// Try a lightweight HEAD request to the formulae API
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, "https://formulae.brew.sh/api/formula.json", nil)
	if err != nil {
		return false, &types.NotAvailableError{Backend: "brew", Reason: "failed to create request: " + err.Error()}
	}

	resp, err := b.httpClient.Do(req)
	if err != nil {
		return false, &types.NotAvailableError{Backend: "brew", Reason: "failed to reach formulae API: " + err.Error()}
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return true, nil
	}

	return false, &types.NotAvailableError{Backend: "brew", Reason: "formulae API returned non-2xx status"}
}

// Capabilities returns brew capabilities.
func (b *Backend) Capabilities(ctx context.Context) ([]types.Capability, error) {
	// Brew backend supports operations when runner is available
	hasRunner := b.runner != nil
	return []types.Capability{
		{Operation: types.OperationSearch, Supported: true, Notes: "via Formulae API"},
		{Operation: types.OperationUpdateMetadata, Supported: hasRunner, Notes: "via brew update CLI"},
		{Operation: types.OperationUpgradePackages, Supported: hasRunner, Notes: "via brew upgrade CLI"},
		{Operation: types.OperationInstall, Supported: hasRunner, Notes: "via brew install CLI"},
		{Operation: types.OperationUninstall, Supported: hasRunner, Notes: "via brew uninstall CLI"},
		{Operation: types.OperationListInstalled, Supported: hasRunner, Notes: "via brew list CLI"},
	}, nil
}

// Update implements Updater using `brew update`.
func (b *Backend) Update(ctx context.Context, opts types.UpdateOptions) (types.UpdateResult, error) {
	if b.runner == nil {
		return types.UpdateResult{}, types.ErrNotSupported
	}

	helper := types.NewProgressHelper(b.progress, opts.Progress)
	helper.BeginAction("Update")
	defer helper.EndAction()

	helper.BeginTask("Running brew update")
	stdout, _, err := runner.RunWithExternalError(
		ctx,
		b.runner,
		types.OperationUpdateMetadata,
		"brew",
		"brew",
		"update",
	)
	helper.EndTask()

	if err != nil {
		helper.Error("Update failed: " + err.Error())
		return types.UpdateResult{}, err
	}

	// Check if there were updates by looking at output
	changed := strings.Contains(stdout, "Updated") || strings.Contains(stdout, "Homebrew updated")

	helper.Info("Update completed")
	return types.UpdateResult{Changed: changed}, nil
}

// Upgrade implements Upgrader using `brew upgrade`.
func (b *Backend) Upgrade(ctx context.Context, opts types.UpgradeOptions) (types.UpgradeResult, error) {
	if b.runner == nil {
		return types.UpgradeResult{}, types.ErrNotSupported
	}

	helper := types.NewProgressHelper(b.progress, opts.Progress)
	helper.BeginAction("Upgrade")
	defer helper.EndAction()

	helper.BeginTask("Running brew upgrade")
	stdout, _, err := runner.RunWithExternalError(
		ctx,
		b.runner,
		types.OperationUpgradePackages,
		"brew",
		"brew",
		"upgrade",
	)
	helper.EndTask()

	if err != nil {
		helper.Error("Upgrade failed: " + err.Error())
		return types.UpgradeResult{}, err
	}

	// Parse upgraded packages from output
	var packagesChanged []types.PackageRef
	changed := false

	// Look for lines like "==> Upgrading <package>"
	lines := strings.Split(stdout, "\n")
	for _, line := range lines {
		if strings.Contains(line, "==> Upgrading") {
			changed = true
			// Extract package name after "Upgrading "
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				pkgName := parts[2]
				packagesChanged = append(packagesChanged, types.PackageRef{
					Name: pkgName,
					Kind: "formula",
				})
			}
		}
	}

	if changed {
		helper.Info("Upgrade completed: upgraded packages")
	} else {
		helper.Info("Upgrade completed: no packages needed upgrading")
	}

	return types.UpgradeResult{
		Changed:         changed,
		PackagesChanged: packagesChanged,
	}, nil
}

// Install implements Installer using `brew install`.
func (b *Backend) Install(ctx context.Context, pkgs []types.PackageRef, opts types.InstallOptions) (types.InstallResult, error) {
	if b.runner == nil {
		return types.InstallResult{}, types.ErrNotSupported
	}

	if len(pkgs) == 0 {
		return types.InstallResult{}, nil
	}

	helper := types.NewProgressHelper(b.progress, opts.Progress)
	helper.BeginAction("Install")
	defer helper.EndAction()

	// Build package list
	pkgNames := make([]string, 0, len(pkgs)+1)
	pkgNames = append(pkgNames, "install")
	for _, pkg := range pkgs {
		pkgNames = append(pkgNames, pkg.Name)
	}

	helper.BeginTask("Running brew install")
	stdout, _, err := runner.RunWithExternalError(
		ctx,
		b.runner,
		types.OperationInstall,
		"brew",
		"brew",
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
		if strings.Contains(line, "==> Installing") || strings.Contains(line, "==> Downloading") {
			changed = true
		}
	}

	// Assume all requested packages were installed
	if changed {
		installed = pkgs
		helper.Info("Install completed: installed packages")
	} else {
		helper.Info("Install completed: packages already installed")
	}

	return types.InstallResult{
		Changed:           changed,
		PackagesInstalled: installed,
	}, nil
}

// Uninstall implements Uninstaller using `brew uninstall`.
func (b *Backend) Uninstall(ctx context.Context, pkgs []types.PackageRef, opts types.UninstallOptions) (types.UninstallResult, error) {
	if b.runner == nil {
		return types.UninstallResult{}, types.ErrNotSupported
	}

	if len(pkgs) == 0 {
		return types.UninstallResult{}, nil
	}

	helper := types.NewProgressHelper(b.progress, opts.Progress)
	helper.BeginAction("Uninstall")
	defer helper.EndAction()

	// Build package list
	pkgNames := make([]string, 0, len(pkgs)+1)
	pkgNames = append(pkgNames, "uninstall")
	for _, pkg := range pkgs {
		pkgNames = append(pkgNames, pkg.Name)
	}

	helper.BeginTask("Running brew uninstall")
	stdout, _, err := runner.RunWithExternalError(
		ctx,
		b.runner,
		types.OperationUninstall,
		"brew",
		"brew",
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

	// Look for uninstallation confirmations
	lines := strings.Split(stdout, "\n")
	for _, line := range lines {
		if strings.Contains(line, "==> Uninstalling") {
			changed = true
		}
	}

	// Assume all requested packages were uninstalled
	if changed {
		uninstalled = pkgs
		helper.Info("Uninstall completed: uninstalled packages")
	} else {
		helper.Info("Uninstall completed: packages not found")
	}

	return types.UninstallResult{
		Changed:             changed,
		PackagesUninstalled: uninstalled,
	}, nil
}

// Search implements Searcher using the Formulae API.
func (b *Backend) Search(ctx context.Context, query string, opts types.SearchOptions) ([]types.PackageRef, error) {
	helper := types.NewProgressHelper(b.progress, opts.Progress)
	helper.BeginAction("Search")
	defer helper.EndAction()

	if query == "" {
		helper.Info("Empty search query")
		return []types.PackageRef{}, nil
	}

	helper.BeginTask("Fetch formulae")
	results, err := b.searchFormulae(ctx, query)
	helper.EndTask()

	if err != nil {
		helper.Error("Search failed: " + err.Error())
		return nil, err
	}

	helper.Info("Search completed")
	return results, nil
}

// ListInstalled implements Lister using `brew list`.
func (b *Backend) ListInstalled(ctx context.Context, opts types.ListOptions) ([]types.InstalledPackage, error) {
	if b.runner == nil {
		return nil, types.ErrNotSupported
	}

	helper := types.NewProgressHelper(b.progress, opts.Progress)
	helper.BeginAction("ListInstalled")
	defer helper.EndAction()

	helper.BeginTask("Running brew list")
	stdout, _, err := runner.RunWithExternalError(
		ctx,
		b.runner,
		types.OperationListInstalled,
		"brew",
		"brew",
		"list",
		"--versions",
	)
	helper.EndTask()

	if err != nil {
		helper.Error("ListInstalled failed: " + err.Error())
		return nil, err
	}

	// Parse output: each line is "package version"
	var installed []types.InstalledPackage
	lines := strings.Split(stdout, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) >= 1 {
			pkg := types.InstalledPackage{
				Ref: types.PackageRef{
					Name: parts[0],
					Kind: "formula",
				},
			}
			if len(parts) >= 2 {
				pkg.Version = parts[1]
			}
			installed = append(installed, pkg)
		}
	}

	helper.Info("ListInstalled completed")
	return installed, nil
}
