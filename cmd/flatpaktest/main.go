package main

import (
	"context"
	"fmt"
	"os"

	"github.com/frostyard/pm"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// Create a simple progress reporter
	progress := &progressReporter{}

	// Create Flatpak backend
	backend := pm.NewFlatpak(pm.WithProgress(progress))

	ctx := context.Background()

	// Check if backend is available
	available, err := backend.Available(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error checking availability: %v\n", err)
		if pm.IsNotAvailable(err) {
			fmt.Fprintf(os.Stderr, "Flatpak is not available on this system\n")
			os.Exit(1)
		}
	}
	if !available {
		fmt.Fprintf(os.Stderr, "Flatpak is not available\n")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "search":
		if len(os.Args) < 3 {
			fmt.Fprintf(os.Stderr, "Usage: %s search <query>\n", os.Args[0])
			os.Exit(1)
		}
		handleSearch(ctx, backend, os.Args[2])

	case "list":
		handleList(ctx, backend)

	case "install":
		if len(os.Args) < 3 {
			fmt.Fprintf(os.Stderr, "Usage: %s install <package>...\n", os.Args[0])
			os.Exit(1)
		}
		handleInstall(ctx, backend, os.Args[2:])

	case "uninstall":
		if len(os.Args) < 3 {
			fmt.Fprintf(os.Stderr, "Usage: %s uninstall <package>...\n", os.Args[0])
			os.Exit(1)
		}
		handleUninstall(ctx, backend, os.Args[2:])

	case "update":
		handleUpdate(ctx, backend)

	case "upgrade":
		handleUpgrade(ctx, backend)

	case "capabilities":
		handleCapabilities(ctx, backend)

	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Printf("Usage: %s <command> [args]\n\n", os.Args[0])
	fmt.Println("Commands:")
	fmt.Println("  search <query>       Search for packages")
	fmt.Println("  list                 List installed packages")
	fmt.Println("  install <package>... Install packages")
	fmt.Println("  uninstall <package>...Remove packages")
	fmt.Println("  update               Update package metadata")
	fmt.Println("  upgrade              Upgrade installed packages")
	fmt.Println("  capabilities         Show backend capabilities")
}

func handleSearch(ctx context.Context, backend pm.Manager, query string) {
	searcher, ok := backend.(pm.Searcher)
	if !ok {
		fmt.Fprintf(os.Stderr, "Backend does not support search\n")
		os.Exit(1)
	}

	results, err := searcher.Search(ctx, query, pm.SearchOptions{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Search failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Found %d packages:\n", len(results))
	for _, pkg := range results {
		fmt.Printf("  %s (%s)\n", pkg.Name, pkg.Kind)
	}
}

func handleList(ctx context.Context, backend pm.Manager) {
	lister, ok := backend.(pm.Lister)
	if !ok {
		fmt.Fprintf(os.Stderr, "Backend does not support list\n")
		os.Exit(1)
	}

	packages, err := lister.ListInstalled(ctx, pm.ListOptions{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "List failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Installed packages (%d):\n", len(packages))
	for _, pkg := range packages {
		namespace := pkg.Ref.Namespace
		if namespace == "" {
			namespace = "(unknown)"
		}
		fmt.Printf("  %-50s %-15s [%s]\n", pkg.Ref.Name, pkg.Version, namespace)
	}
}

func handleInstall(ctx context.Context, backend pm.Manager, packages []string) {
	installer, ok := backend.(pm.Installer)
	if !ok {
		fmt.Fprintf(os.Stderr, "Backend does not support install\n")
		os.Exit(1)
	}

	pkgs := make([]pm.PackageRef, len(packages))
	for i, name := range packages {
		pkgs[i] = pm.PackageRef{Name: name}
	}

	result, err := installer.Install(ctx, pkgs, pm.InstallOptions{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Install failed: %v\n", err)
		os.Exit(1)
	}

	if result.Changed {
		fmt.Printf("Successfully installed %d packages\n", len(result.PackagesInstalled))
	} else {
		fmt.Println("No changes made (packages already installed)")
	}
}

func handleUninstall(ctx context.Context, backend pm.Manager, packages []string) {
	uninstaller, ok := backend.(pm.Uninstaller)
	if !ok {
		fmt.Fprintf(os.Stderr, "Backend does not support uninstall\n")
		os.Exit(1)
	}

	pkgs := make([]pm.PackageRef, len(packages))
	for i, name := range packages {
		pkgs[i] = pm.PackageRef{Name: name}
	}

	result, err := uninstaller.Uninstall(ctx, pkgs, pm.UninstallOptions{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Uninstall failed: %v\n", err)
		os.Exit(1)
	}

	if result.Changed {
		fmt.Printf("Successfully uninstalled %d packages\n", len(result.PackagesUninstalled))
	} else {
		fmt.Println("No changes made (packages not installed)")
	}
}

func handleUpdate(ctx context.Context, backend pm.Manager) {
	updater, ok := backend.(pm.Updater)
	if !ok {
		fmt.Fprintf(os.Stderr, "Backend does not support update\n")
		os.Exit(1)
	}

	result, err := updater.Update(ctx, pm.UpdateOptions{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Update failed: %v\n", err)
		os.Exit(1)
	}

	if result.Changed {
		fmt.Println("Package metadata updated")
	} else {
		fmt.Println("Package metadata already up to date")
	}
}

func handleUpgrade(ctx context.Context, backend pm.Manager) {
	upgrader, ok := backend.(pm.Upgrader)
	if !ok {
		fmt.Fprintf(os.Stderr, "Backend does not support upgrade\n")
		os.Exit(1)
	}

	result, err := upgrader.Upgrade(ctx, pm.UpgradeOptions{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Upgrade failed: %v\n", err)
		os.Exit(1)
	}

	if result.Changed {
		fmt.Printf("Successfully upgraded %d packages\n", len(result.PackagesChanged))
		for _, pkg := range result.PackagesChanged {
			fmt.Printf("  - %s\n", pkg.Name)
		}
	} else {
		fmt.Println("All packages are up to date")
	}
}

func handleCapabilities(ctx context.Context, backend pm.Manager) {
	caps, err := backend.Capabilities(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get capabilities: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Backend Capabilities:")
	for _, cap := range caps {
		status := "✗"
		if cap.Supported {
			status = "✓"
		}
		notes := ""
		if cap.Notes != "" {
			notes = fmt.Sprintf(" (%s)", cap.Notes)
		}
		fmt.Printf("  %s %s%s\n", status, cap.Operation, notes)
	}
}

// progressReporter is a simple progress reporter that prints to stdout
type progressReporter struct{}

func (p *progressReporter) OnAction(action pm.ProgressAction) {
	if !action.StartedAt.IsZero() && action.EndedAt.IsZero() {
		fmt.Printf("→ %s\n", action.Name)
	}
}

func (p *progressReporter) OnTask(task pm.ProgressTask) {
	if !task.StartedAt.IsZero() && task.EndedAt.IsZero() {
		fmt.Printf("  • %s\n", task.Name)
	}
}

func (p *progressReporter) OnStep(step pm.ProgressStep) {
	if !step.StartedAt.IsZero() && step.EndedAt.IsZero() {
		fmt.Printf("    - %s\n", step.Name)
	}
}

func (p *progressReporter) OnMessage(msg pm.ProgressMessage) {
	prefix := ""
	switch msg.Severity {
	case pm.SeverityInfo:
		prefix = "ℹ"
	case pm.SeverityWarning:
		prefix = "⚠"
	case pm.SeverityError:
		prefix = "✗"
	}
	fmt.Printf("    %s %s\n", prefix, msg.Text)
}
