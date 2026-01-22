package pm

// Operation represents a package manager operation type.
type Operation string

const (
	// OperationUpdateMetadata refreshes package metadata/indexes without changing installed packages.
	OperationUpdateMetadata Operation = "UpdateMetadata"

	// OperationUpgradePackages may change installed packages (upgrade versions).
	OperationUpgradePackages Operation = "UpgradePackages"

	// OperationInstall installs one or more packages.
	OperationInstall Operation = "Install"

	// OperationUninstall removes one or more packages.
	OperationUninstall Operation = "Uninstall"

	// OperationSearch searches for available packages.
	OperationSearch Operation = "Search"

	// OperationListInstalled lists installed packages.
	OperationListInstalled Operation = "ListInstalled"

	// OperationListAvailable lists available packages (if supported).
	OperationListAvailable Operation = "ListAvailable"
)

// PackageRef identifies a package in a backend-agnostic way.
type PackageRef struct {
	// Name is the package name (required).
	Name string

	// Namespace is an optional namespace/scope (e.g., flatpak remote, snap publisher).
	Namespace string

	// Channel is an optional channel (e.g., snap channel: stable, edge).
	Channel string

	// Kind is an optional package kind (e.g., brew cask vs formula, flatpak app vs runtime).
	Kind string
}

// InstalledPackage represents a package currently installed on the system.
type InstalledPackage struct {
	// Ref is the package reference.
	Ref PackageRef

	// Version is the installed version.
	Version string

	// Status is the installation status (e.g., "installed", "held", "disabled").
	Status string
}

// Capability represents an operation that a backend supports.
type Capability struct {
	// Operation is the operation type.
	Operation Operation

	// Supported indicates whether the operation is supported.
	Supported bool

	// Notes provides optional context (e.g., why unsupported, constraints).
	Notes string
}
