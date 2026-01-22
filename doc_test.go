package pm_test

import (
	"context"
	"fmt"
	"log"

	"github.com/frostyard/pm"
)

// Example_basicUsage demonstrates basic package manager usage.
func Example_basicUsage() {
	// Create a backend (e.g., Homebrew)
	backend := pm.NewBrew()

	// Check if backend is available
	available, err := backend.Available(context.Background())
	if err != nil {
		log.Printf("Error checking availability: %v", err)
		return
	}

	if !available {
		log.Println("Backend not available")
		return
	}

	fmt.Println("Backend is available")
	// Output: Backend is available
}

// Example_searchPackages demonstrates searching for packages.
func Example_searchPackages() {
	backend := pm.NewBrew()

	// Search for packages
	packages, err := backend.(pm.Searcher).Search(context.Background(), "git", pm.SearchOptions{})
	if pm.IsNotSupported(err) {
		fmt.Println("Search not supported")
		return
	}
	if err != nil {
		log.Printf("Search failed: %v", err)
		return
	}

	fmt.Printf("Found %d packages\n", len(packages))
	// Output varies based on backend availability
}

// Example_capabilities demonstrates checking backend capabilities.
func Example_capabilities() {
	backend := pm.NewBrew()

	caps, err := backend.Capabilities(context.Background())
	if err != nil {
		log.Printf("Error getting capabilities: %v", err)
		return
	}

	fmt.Printf("Backend supports %d operations\n", len(caps))
	for _, cap := range caps {
		if cap.Supported {
			fmt.Printf("- %s: supported\n", cap.Operation)
		}
	}
	// Output varies based on backend implementation
}

// Example_updateVsUpgrade demonstrates the difference between Update and Upgrade.
func Example_updateVsUpgrade() {
	backend := pm.NewBrew()

	// Update: refresh metadata only (never modifies packages)
	updateResult, err := backend.(pm.Updater).Update(context.Background(), pm.UpdateOptions{})
	if pm.IsNotSupported(err) {
		fmt.Println("Update operation not supported")
	} else if err == nil {
		fmt.Printf("Metadata updated: %v\n", updateResult.Changed)
	}

	// Upgrade: may modify installed packages
	upgradeResult, err := backend.(pm.Upgrader).Upgrade(context.Background(), pm.UpgradeOptions{})
	if pm.IsNotSupported(err) {
		fmt.Println("Upgrade operation not supported")
	} else if err == nil {
		fmt.Printf("Packages upgraded: %d\n", len(upgradeResult.PackagesChanged))
	}

	// Output varies based on backend support
}

// Example_errorHandling demonstrates error detection and handling.
func Example_errorHandling() {
	backend := pm.NewSnap()

	_, err := backend.(pm.Installer).Install(
		context.Background(),
		[]pm.PackageRef{{Name: "example"}},
		pm.InstallOptions{},
	)

	// Check error types
	if pm.IsNotSupported(err) {
		fmt.Println("Operation not supported by this backend")
	} else if pm.IsNotAvailable(err) {
		fmt.Println("Backend not available")
	} else if pm.IsExternalFailure(err) {
		fmt.Println("External command/API failed")
	} else if err != nil {
		fmt.Printf("Other error: %v\n", err)
	}

	// Output: External command/API failed
}

// Example_progressReporting demonstrates progress reporting during operations.
func Example_progressReporting() {
	backend := pm.NewBrew()

	// Create a progress reporter
	reporter := &SimpleProgressReporter{}

	// Search with progress reporting
	_, err := backend.(pm.Searcher).Search(
		context.Background(),
		"nodejs",
		pm.SearchOptions{Progress: reporter},
	)

	if err != nil && !pm.IsNotSupported(err) {
		log.Printf("Search failed: %v", err)
	}

	fmt.Printf("Received %d progress events\n", reporter.EventCount)
	// Output varies based on backend behavior
}

// SimpleProgressReporter is a basic progress reporter for examples.
type SimpleProgressReporter struct {
	EventCount int
}

func (r *SimpleProgressReporter) OnAction(action pm.ProgressAction) {
	r.EventCount++
	if action.EndedAt.IsZero() {
		fmt.Printf("Action started: %s\n", action.Name)
	}
}

func (r *SimpleProgressReporter) OnTask(task pm.ProgressTask) {
	r.EventCount++
}

func (r *SimpleProgressReporter) OnStep(step pm.ProgressStep) {
	r.EventCount++
}

func (r *SimpleProgressReporter) OnMessage(msg pm.ProgressMessage) {
	r.EventCount++
	if msg.Severity == pm.SeverityWarning {
		fmt.Printf("Warning: %s\n", msg.Text)
	}
}

// Example_multipleBackends demonstrates working with multiple backends.
func Example_multipleBackends() {
	backends := map[string]pm.Manager{
		"brew":    pm.NewBrew(),
		"flatpak": pm.NewFlatpak(),
		"snap":    pm.NewSnap(),
	}

	ctx := context.Background()

	// Check which backends are available
	fmt.Println("Checking backend availability:")
	for name, backend := range backends {
		available, err := backend.Available(ctx)
		if err != nil {
			fmt.Printf("- %s: error (%v)\n", name, err)
			continue
		}
		if available {
			fmt.Printf("- %s: available\n", name)
		} else {
			fmt.Printf("- %s: not available\n", name)
		}
	}

	// Output varies based on system configuration
}

// Example_typeAssertion demonstrates safe type assertion for optional interfaces.
func Example_typeAssertion() {
	backend := pm.NewBrew()

	// Check if backend supports search
	if searcher, ok := backend.(pm.Searcher); ok {
		fmt.Println("Backend supports Search")
		_, _ = searcher.Search(context.Background(), "test", pm.SearchOptions{})
	}

	// Check if backend supports install
	if installer, ok := backend.(pm.Installer); ok {
		fmt.Println("Backend supports Install")
		_, _ = installer.Install(context.Background(), []pm.PackageRef{}, pm.InstallOptions{})
	}

	// Check if backend supports upgrade
	if upgrader, ok := backend.(pm.Upgrader); ok {
		fmt.Println("Backend supports Upgrade")
		_, _ = upgrader.Upgrade(context.Background(), pm.UpgradeOptions{})
	}

	// Output varies based on backend implementation
}

// Example_packageReferences demonstrates working with package references.
func Example_packageReferences() {
	// Create package references
	packages := []pm.PackageRef{
		{Name: "git", Kind: "formula"},
		{Name: "nodejs", Kind: "formula"},
		{Name: "python@3.11", Kind: "formula"},
	}

	fmt.Printf("Preparing to install %d packages:\n", len(packages))
	for _, pkg := range packages {
		fmt.Printf("- %s (%s)\n", pkg.Name, pkg.Kind)
	}

	// Output:
	// Preparing to install 3 packages:
	// - git (formula)
	// - nodejs (formula)
	// - python@3.11 (formula)
}

// Example_operationOptions demonstrates configuring operation options.
func Example_operationOptions() {
	backend := pm.NewBrew()

	// Configure install options (with nil progress for simplicity)
	installOpts := pm.InstallOptions{
		Progress: nil, // Progress is optional
	}

	// Configure search options
	searchOpts := pm.SearchOptions{
		Progress: nil, // Progress is optional
	}

	// Use options in operations
	_, _ = backend.(pm.Installer).Install(
		context.Background(),
		[]pm.PackageRef{{Name: "example"}},
		installOpts,
	)

	_, _ = backend.(pm.Searcher).Search(
		context.Background(),
		"test",
		searchOpts,
	)

	fmt.Println("Operations configured with options")
	// Output: Operations configured with options
}

// Example_contextCancellation demonstrates context-based cancellation.
func Example_contextCancellation() {
	backend := pm.NewBrew()
	ctx, cancel := context.WithCancel(context.Background())

	// Cancel immediately for demonstration
	cancel()

	// Operation should respect context cancellation
	_, err := backend.(pm.Searcher).Search(ctx, "test", pm.SearchOptions{})
	if err == context.Canceled {
		fmt.Println("Operation was cancelled")
	} else if pm.IsNotSupported(err) {
		fmt.Println("Operation not supported")
	}

	// Output varies based on backend implementation
}
