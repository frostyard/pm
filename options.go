package pm

// UpdateOptions provides options for Update operations.
//
// Update operations refresh package metadata/indexes without modifying
// installed packages. This is analogous to 'apt update' or 'brew update'.
type UpdateOptions struct {
	// Progress is an optional progress reporter.
	Progress ProgressReporter
}

// UpdateResult is the result of an Update operation.
//
// Contract guarantees:
//   - Changed=false means no metadata updates were available or needed
//   - Changed=true means metadata/indexes were successfully refreshed
//   - No installed packages are ever modified by Update operations
type UpdateResult struct {
	// Changed indicates whether metadata was refreshed.
	// Will be false if metadata was already current or if operation failed.
	Changed bool

	// Messages contains summary messages from the operation.
	Messages []ProgressMessage
}

// UpgradeOptions provides options for Upgrade operations.
//
// Upgrade operations install newer versions of packages that are already installed.
// This is analogous to 'apt upgrade' or 'brew upgrade'.
type UpgradeOptions struct {
	// Progress is an optional progress reporter.
	Progress ProgressReporter
}

// UpgradeResult is the result of an Upgrade operation.
//
// Contract guarantees:
//   - Changed=false means no packages were upgraded (all current, or operation failed)
//   - Changed=true means one or more packages were successfully upgraded
//   - PackagesChanged lists the specific packages that were upgraded
//   - Upgrade operations may modify installed software and system state
type UpgradeResult struct {
	// Changed indicates whether any packages were changed.
	// Will be false if all packages were current or if operation failed.
	Changed bool

	// PackagesChanged lists packages that were upgraded.
	// Empty if Changed=false.
	PackagesChanged []PackageRef

	// Messages contains summary messages from the operation.
	Messages []ProgressMessage
}

// InstallOptions provides options for Install operations.
type InstallOptions struct {
	// Progress is an optional progress reporter.
	Progress ProgressReporter
}

// InstallResult is the result of an Install operation.
type InstallResult struct {
	// Changed indicates whether any packages were installed.
	Changed bool

	// PackagesInstalled lists packages that were installed.
	PackagesInstalled []PackageRef

	// Messages contains summary messages from the operation.
	Messages []ProgressMessage
}

// UninstallOptions provides options for Uninstall operations.
type UninstallOptions struct {
	// Progress is an optional progress reporter.
	Progress ProgressReporter
}

// UninstallResult is the result of an Uninstall operation.
type UninstallResult struct {
	// Changed indicates whether any packages were uninstalled.
	Changed bool

	// PackagesUninstalled lists packages that were uninstalled.
	PackagesUninstalled []PackageRef

	// Messages contains summary messages from the operation.
	Messages []ProgressMessage
}

// SearchOptions provides options for Search operations.
type SearchOptions struct {
	// Progress is an optional progress reporter.
	Progress ProgressReporter
}

// ListOptions provides options for ListInstalled operations.
type ListOptions struct {
	// Progress is an optional progress reporter.
	Progress ProgressReporter
}
