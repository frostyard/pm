package pm

import "context"

// Manager provides core backend functionality: availability and capability introspection.
type Manager interface {
	// Available checks if the backend is available (installed/reachable).
	Available(ctx context.Context) (bool, error)

	// Capabilities returns the operations this backend supports.
	Capabilities(ctx context.Context) ([]Capability, error)
}

// Updater updates package metadata without changing installed packages.
//
// Semantics Contract:
//   - Update MUST only refresh metadata/indexes (e.g., apt update, brew update)
//   - Update MUST NOT install, remove, or upgrade any packages
//   - Update SHOULD set UpdateResult.Changed=true if metadata was refreshed
//   - Update SHOULD set UpdateResult.Changed=false if no changes were needed
//
// Examples:
//   - apt update: refreshes package lists
//   - brew update: updates Homebrew and formula metadata
//   - flatpak update --appstream: updates app metadata
//
// This operation is safe to run repeatedly and never modifies installed software.
type Updater interface {
	Update(ctx context.Context, opts UpdateOptions) (UpdateResult, error)
}

// Upgrader upgrades installed packages to newer versions.
//
// Semantics Contract:
//   - Upgrade MAY install newer versions of installed packages
//   - Upgrade MAY change system state (packages, configurations)
//   - Upgrade SHOULD set UpgradeResult.Changed=true if any packages were upgraded
//   - Upgrade SHOULD populate UpgradeResult.PackagesChanged with upgraded packages
//   - Upgrade SHOULD set UpgradeResult.Changed=false if no upgrades were available
//
// Examples:
//   - apt upgrade: upgrades installed packages to newer versions
//   - brew upgrade: upgrades outdated formulae
//   - flatpak update: updates installed applications
//
// Note: Some backends may require Update to be called first to refresh metadata.
// This operation modifies installed software and may require elevated privileges.
type Upgrader interface {
	Upgrade(ctx context.Context, opts UpgradeOptions) (UpgradeResult, error)
}

// Installer installs packages.
type Installer interface {
	Install(ctx context.Context, pkgs []PackageRef, opts InstallOptions) (InstallResult, error)
}

// Uninstaller uninstalls packages.
type Uninstaller interface {
	Uninstall(ctx context.Context, pkgs []PackageRef, opts UninstallOptions) (UninstallResult, error)
}

// Searcher searches for packages.
type Searcher interface {
	Search(ctx context.Context, query string, opts SearchOptions) ([]PackageRef, error)
}

// Lister lists packages.
type Lister interface {
	ListInstalled(ctx context.Context, opts ListOptions) ([]InstalledPackage, error)
}
