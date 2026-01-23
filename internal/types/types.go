// Package types contains shared types used by internal backend implementations.
// This breaks import cycles by avoiding backends importing the parent pm package.
package types

import (
	"errors"
	"fmt"

	"github.com/frostyard/pm/progress"
)

// Core errors that backends can return.
var (
	ErrNotSupported = errors.New("operation not supported")
	ErrNotAvailable = errors.New("backend not available")
)

// IsNotSupported checks if an error is a NotSupported error.
func IsNotSupported(err error) bool {
	return errors.Is(err, ErrNotSupported)
}

// IsNotAvailable checks if an error is a NotAvailable error.
func IsNotAvailable(err error) bool {
	return errors.Is(err, ErrNotAvailable)
}

// NotSupportedError wraps ErrNotSupported with additional context.
type NotSupportedError struct {
	Operation Operation
	Backend   string
	Reason    string
}

func (e *NotSupportedError) Error() string {
	if e.Reason != "" {
		return fmt.Sprintf("%s: %s operation not supported by %s: %s", ErrNotSupported, e.Operation, e.Backend, e.Reason)
	}
	return fmt.Sprintf("%s: %s operation not supported by %s", ErrNotSupported, e.Operation, e.Backend)
}

func (e *NotSupportedError) Unwrap() error {
	return ErrNotSupported
}

// NotAvailableError wraps ErrNotAvailable with additional context.
type NotAvailableError struct {
	Backend string
	Reason  string
}

func (e *NotAvailableError) Error() string {
	if e.Reason != "" {
		return fmt.Sprintf("%s: %s: %s", ErrNotAvailable, e.Backend, e.Reason)
	}
	return fmt.Sprintf("%s: %s", ErrNotAvailable, e.Backend)
}

func (e *NotAvailableError) Unwrap() error {
	return ErrNotAvailable
}

// ExternalFailureError represents a failure from an external command or API.
type ExternalFailureError struct {
	Operation Operation
	Backend   string
	Stdout    string
	Stderr    string
	Payload   map[string]interface{}
	Err       error
}

func (e *ExternalFailureError) Error() string {
	msg := fmt.Sprintf("external failure: %s operation on %s", e.Operation, e.Backend)
	if e.Err != nil {
		msg = fmt.Sprintf("%s: %v", msg, e.Err)
	}
	if e.Stderr != "" {
		msg = fmt.Sprintf("%s (stderr: %s)", msg, e.Stderr)
	}
	return msg
}

func (e *ExternalFailureError) Unwrap() error {
	return e.Err
}

// IsExternalFailure checks if an error is an ExternalFailure error.
func IsExternalFailure(err error) bool {
	var extErr *ExternalFailureError
	return errors.As(err, &extErr)
}

// PackageRef mirrors pm.PackageRef for internal use.
type PackageRef struct {
	Name      string
	Namespace string
	Channel   string
	Kind      string
}

// InstalledPackage mirrors pm.InstalledPackage for internal use.
type InstalledPackage struct {
	Ref     PackageRef
	Version string
	Status  string
}

// Operation mirrors pm.Operation for internal use.
type Operation string

const (
	OperationUpdateMetadata  Operation = "UpdateMetadata"
	OperationUpgradePackages Operation = "UpgradePackages"
	OperationInstall         Operation = "Install"
	OperationUninstall       Operation = "Uninstall"
	OperationSearch          Operation = "Search"
	OperationListInstalled   Operation = "ListInstalled"
)

// Capability mirrors pm.Capability for internal use.
type Capability struct {
	Operation Operation
	Supported bool
	Notes     string
}

// Progress reporter types from progress module.
type (
	// Severity represents the severity level of a progress message.
	Severity = progress.Severity

	// ProgressMessage is a message emitted during progress.
	ProgressMessage = progress.ProgressMessage

	// ProgressAction represents a high-level action in a long-running operation.
	ProgressAction = progress.ProgressAction

	// ProgressTask represents a task within an action.
	ProgressTask = progress.ProgressTask

	// ProgressStep represents a step within a task.
	ProgressStep = progress.ProgressStep

	// ProgressReporter is the interface for progress reporting.
	ProgressReporter = progress.ProgressReporter

	// ProgressHelper provides a convenient API for backends to emit progress updates.
	ProgressHelper = progress.ProgressHelper
)

// Re-export severity constants
const (
	SeverityInfo    = progress.SeverityInfo
	SeverityWarning = progress.SeverityWarning
	SeverityError   = progress.SeverityError
)

// NewProgressHelper creates a new progress helper with progress reporting.
func NewProgressHelper(defaultReporter, overrideReporter ProgressReporter) *ProgressHelper {
	return progress.NewProgressHelper(defaultReporter, overrideReporter)
}

// Result types for operations.
type UpdateResult struct {
	Changed  bool
	Messages []ProgressMessage
}

type UpgradeResult struct {
	Changed         bool
	PackagesChanged []PackageRef
	Messages        []ProgressMessage
}

type InstallResult struct {
	Changed           bool
	PackagesInstalled []PackageRef
	Messages          []ProgressMessage
}

type UninstallResult struct {
	Changed             bool
	PackagesUninstalled []PackageRef
	Messages            []ProgressMessage
}

// Options types for operations.
type UpdateOptions struct {
	Progress ProgressReporter
}

type UpgradeOptions struct {
	Progress ProgressReporter
}

type InstallOptions struct {
	Progress ProgressReporter
}

type UninstallOptions struct {
	Progress ProgressReporter
}

type SearchOptions struct {
	Progress ProgressReporter
}

type ListOptions struct {
	Progress ProgressReporter
}
