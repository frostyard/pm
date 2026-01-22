// Package pm provides a unified Go API for interacting with multiple package managers.
//
// This library abstracts Homebrew (brew), Flatpak, and Snap behind a common set of
// interfaces, enabling applications to support multiple package managers without
// manager-specific code. Each backend implements shared contracts for operations like
// search, install, upgrade, and list.
//
// # Key Features
//
//   - Common interfaces: Manager, Searcher, Installer, Upgrader, Lister, etc.
//   - Empty implementations: Backends can return NotSupported for unimplemented operations
//   - Capability introspection: Check what operations a backend supports before calling
//   - Progress reporting: Optional structured progress (Actions → Tasks → Steps) for long operations
//   - Clear semantics: Update (metadata only) vs Upgrade (may change packages)
//
// # Example Usage
//
//	import (
//	    "context"
//	    "github.com/frostyard/pm"
//	)
//
//	func example() {
//	    ctx := context.Background()
//	    mgr := pm.NewBrew()
//
//	    // Check availability
//	    available, err := mgr.Available(ctx)
//	    if err != nil || !available {
//	        // Handle unavailable backend
//	    }
//
//	    // Check capabilities
//	    caps, err := mgr.Capabilities(ctx)
//	    // Use caps to determine what operations are supported
//	}
//
// # Backend Integration
//
// The library prefers SDK/API integration where available, then REST, then CLI as a last resort.
// For deterministic testing, all system interactions are abstracted behind injectable dependencies.
package pm
