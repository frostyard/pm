package flatpak

import (
	"context"
	"strings"

	"github.com/frostyard/pm/internal/runner"
	"github.com/frostyard/pm/internal/types"
)

// Backend implements the flatpak backend.
type Backend struct {
	runner   runner.Runner
	progress types.ProgressReporter
}

// New creates a new flatpak backend.
func New(r runner.Runner, progress types.ProgressReporter) *Backend {
	return &Backend{
		runner:   r,
		progress: progress,
	}
}

// Available checks if flatpak is available by running `flatpak --version`.
func (b *Backend) Available(ctx context.Context) (bool, error) {
	if b.runner == nil {
		return false, &types.NotAvailableError{Backend: "flatpak", Reason: "no runner configured"}
	}

	stdout, stderr, err := b.runner.Run(ctx, "flatpak", "--version")
	if err != nil {
		return false, &types.NotAvailableError{Backend: "flatpak", Reason: "flatpak --version failed: " + stderr + ": " + err.Error()}
	}

	// If we got a successful execution and output, flatpak is available
	if len(stdout) > 0 {
		return true, nil
	}

	return false, &types.NotAvailableError{Backend: "flatpak", Reason: "flatpak --version returned no output"}
}

// Capabilities returns flatpak capabilities.
func (b *Backend) Capabilities(ctx context.Context) ([]types.Capability, error) {
	// Flatpak backend supports operations when runner is available
	hasRunner := b.runner != nil
	return []types.Capability{
		{Operation: types.OperationSearch, Supported: hasRunner, Notes: "via flatpak search CLI"},
		{Operation: types.OperationUpdateMetadata, Supported: hasRunner, Notes: "via flatpak update CLI"},
		{Operation: types.OperationUpgradePackages, Supported: hasRunner, Notes: "via flatpak update CLI"},
		{Operation: types.OperationInstall, Supported: hasRunner, Notes: "via flatpak install CLI"},
		{Operation: types.OperationUninstall, Supported: hasRunner, Notes: "via flatpak uninstall CLI"},
		{Operation: types.OperationListInstalled, Supported: hasRunner, Notes: "via flatpak list CLI"},
	}, nil
}

// Update implements Updater using `flatpak update --appstream`.
func (b *Backend) Update(ctx context.Context, opts types.UpdateOptions) (types.UpdateResult, error) {
	if b.runner == nil {
		return types.UpdateResult{}, types.ErrNotSupported
	}

	helper := types.NewProgressHelper(b.progress, opts.Progress)
	helper.BeginAction("Update")
	defer helper.EndAction()

	helper.BeginTask("Running flatpak update --appstream")
	stdout, _, err := runner.RunWithExternalError(
		ctx,
		b.runner,
		types.OperationUpdateMetadata,
		"flatpak",
		"flatpak",
		"update",
		"--appstream",
	)
	helper.EndTask()

	if err != nil {
		helper.Error("Update failed: " + err.Error())
		return types.UpdateResult{}, err
	}

	// Check if there were updates by looking at output
	changed := strings.Contains(stdout, "Updating") || strings.Contains(stdout, "Updated")

	helper.Info("Update completed")
	return types.UpdateResult{Changed: changed}, nil
}

// Upgrade implements Upgrader using `flatpak update`.
func (b *Backend) Upgrade(ctx context.Context, opts types.UpgradeOptions) (types.UpgradeResult, error) {
	if b.runner == nil {
		return types.UpgradeResult{}, types.ErrNotSupported
	}

	helper := types.NewProgressHelper(b.progress, opts.Progress)
	helper.BeginAction("Upgrade")
	defer helper.EndAction()

	helper.BeginTask("Running flatpak update")
	stdout, _, err := runner.RunWithExternalError(
		ctx,
		b.runner,
		types.OperationUpgradePackages,
		"flatpak",
		"flatpak",
		"update",
		"-y",
	)
	helper.EndTask()

	if err != nil {
		helper.Error("Upgrade failed: " + err.Error())
		return types.UpgradeResult{}, err
	}

	// Parse upgraded packages from output
	var packagesChanged []types.PackageRef
	changed := false

	// Look for lines indicating updates
	lines := strings.Split(stdout, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Flatpak shows updating apps like "Updating <app-id>"
		if strings.HasPrefix(line, "Updating") {
			changed = true
			// Extract app ID
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				appID := parts[1]
				packagesChanged = append(packagesChanged, types.PackageRef{
					Name: appID,
					Kind: "app",
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

// Install implements Installer using `flatpak install`.
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

	// Build package list - flatpak install requires app IDs
	pkgNames := make([]string, 0, len(pkgs)+2)
	pkgNames = append(pkgNames, "install", "-y")
	for _, pkg := range pkgs {
		pkgNames = append(pkgNames, pkg.Name)
	}

	helper.BeginTask("Running flatpak install")
	stdout, _, err := runner.RunWithExternalError(
		ctx,
		b.runner,
		types.OperationInstall,
		"flatpak",
		"flatpak",
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
		if strings.Contains(line, "Installing") || strings.Contains(line, "installed") {
			changed = true
			// Try to extract app ID from the line
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
		helper.Info("Install completed: installed packages")
	} else {
		helper.Info("Install completed: packages already installed")
	}

	return types.InstallResult{
		Changed:           changed,
		PackagesInstalled: installed,
	}, nil
}

// Uninstall implements Uninstaller using `flatpak uninstall`.
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
	pkgNames := make([]string, 0, len(pkgs)+2)
	pkgNames = append(pkgNames, "uninstall", "-y")
	for _, pkg := range pkgs {
		pkgNames = append(pkgNames, pkg.Name)
	}

	helper.BeginTask("Running flatpak uninstall")
	stdout, _, err := runner.RunWithExternalError(
		ctx,
		b.runner,
		types.OperationUninstall,
		"flatpak",
		"flatpak",
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

	// Look for uninstall confirmations in output
	lines := strings.Split(stdout, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "Uninstalling") || strings.Contains(line, "uninstalled") {
			changed = true
			// Try to extract app ID from the line
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
		helper.Info("Uninstall completed: uninstalled packages")
	} else {
		helper.Info("Uninstall completed: packages were not installed")
	}

	return types.UninstallResult{
		Changed:             changed,
		PackagesUninstalled: uninstalled,
	}, nil
}

// Search implements Searcher using `flatpak search`.
func (b *Backend) Search(ctx context.Context, query string, opts types.SearchOptions) ([]types.PackageRef, error) {
	if b.runner == nil {
		return nil, types.ErrNotSupported
	}

	if query == "" {
		return []types.PackageRef{}, nil
	}

	helper := types.NewProgressHelper(b.progress, opts.Progress)
	helper.BeginAction("Search")
	defer helper.EndAction()

	helper.BeginTask("Running flatpak search")
	stdout, _, err := runner.RunWithExternalError(
		ctx,
		b.runner,
		types.OperationSearch,
		"flatpak",
		"flatpak",
		"search",
		query,
	)
	helper.EndTask()

	if err != nil {
		helper.Error("Search failed: " + err.Error())
		return nil, err
	}

	// Parse search results
	// Flatpak search output format:
	// Name          Description                     Application ID          Version Branch Remotes
	// Firefox       Web Browser                     org.mozilla.firefox     ...     ...    flathub
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

		// Parse fields - split by whitespace but handle multiple spaces
		fields := strings.Fields(line)
		if len(fields) >= 3 {
			appID := fields[2]

			results = append(results, types.PackageRef{
				Name: appID,
				Kind: "app",
			})
		}
	}

	helper.Info("Search completed")
	return results, nil
}

// ListInstalled implements Lister using `flatpak list`.
func (b *Backend) ListInstalled(ctx context.Context, opts types.ListOptions) ([]types.InstalledPackage, error) {
	if b.runner == nil {
		return nil, types.ErrNotSupported
	}

	helper := types.NewProgressHelper(b.progress, opts.Progress)
	helper.BeginAction("ListInstalled")
	defer helper.EndAction()

	helper.BeginTask("Running flatpak list")
	stdout, _, err := runner.RunWithExternalError(
		ctx,
		b.runner,
		types.OperationListInstalled,
		"flatpak",
		"flatpak",
		"list",
		"--app",
		"--columns=name,application,version,installation",
	)
	helper.EndTask()

	if err != nil {
		helper.Error("ListInstalled failed: " + err.Error())
		return nil, err
	}

	// Parse output: columns are name, application ID, version, installation
	var packages []types.InstalledPackage
	lines := strings.Split(stdout, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Split by tab (flatpak uses tabs for column separation with --columns)
		fields := strings.Split(line, "\t")
		if len(fields) >= 4 {
			appID := strings.TrimSpace(fields[1])
			version := strings.TrimSpace(fields[2])
			installation := strings.TrimSpace(fields[3])

			packages = append(packages, types.InstalledPackage{
				Ref: types.PackageRef{
					Name:      appID,
					Kind:      "app",
					Namespace: installation, // "user" or "system"
				},
				Version: version,
			})
		} else if len(fields) >= 3 {
			// Fallback: if installation column is missing, still parse what we can
			appID := strings.TrimSpace(fields[1])
			version := strings.TrimSpace(fields[2])

			packages = append(packages, types.InstalledPackage{
				Ref: types.PackageRef{
					Name: appID,
					Kind: "app",
				},
				Version: version,
			})
		} else {
			// Fallback: split by whitespace if tabs not present
			fields = strings.Fields(line)
			if len(fields) >= 3 {
				appID := fields[1]
				version := fields[2]
				installation := ""
				if len(fields) >= 4 {
					installation = fields[3]
				}

				packages = append(packages, types.InstalledPackage{
					Ref: types.PackageRef{
						Name:      appID,
						Kind:      "app",
						Namespace: installation,
					},
					Version: version,
				})
			}
		}
	}

	helper.Info("ListInstalled completed")
	return packages, nil
}
