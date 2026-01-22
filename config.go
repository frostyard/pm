package pm

// BackendKind represents a package manager backend type.
type BackendKind string

const (
	// BackendBrew represents Homebrew.
	BackendBrew BackendKind = "brew"

	// BackendFlatpak represents Flatpak.
	BackendFlatpak BackendKind = "flatpak"

	// BackendSnap represents Snap/snapd.
	BackendSnap BackendKind = "snap"
)

// ConstructorOption is a function that configures a backend during construction.
type ConstructorOption func(config *backendConfig)

// backendConfig holds configuration for backend constructors.
type backendConfig struct {
	progress ProgressReporter
}

// WithProgress sets a progress reporter for a backend.
func WithProgress(p ProgressReporter) ConstructorOption {
	return func(config *backendConfig) {
		config.progress = p
	}
}
