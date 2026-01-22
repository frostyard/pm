package pm

import (
	"context"
	"errors"

	"github.com/frostyard/pm/internal/backend/brew"
	"github.com/frostyard/pm/internal/backend/flatpak"
	"github.com/frostyard/pm/internal/backend/snap"
	"github.com/frostyard/pm/internal/runner"
	"github.com/frostyard/pm/internal/types"
)

// backendAdapter wraps internal backend types to expose pm package types.
type backendAdapter struct {
	backend interface {
		Available(ctx context.Context) (bool, error)
		Capabilities(ctx context.Context) ([]types.Capability, error)
		Update(ctx context.Context, opts types.UpdateOptions) (types.UpdateResult, error)
		Upgrade(ctx context.Context, opts types.UpgradeOptions) (types.UpgradeResult, error)
		Install(ctx context.Context, pkgs []types.PackageRef, opts types.InstallOptions) (types.InstallResult, error)
		Uninstall(ctx context.Context, pkgs []types.PackageRef, opts types.UninstallOptions) (types.UninstallResult, error)
		Search(ctx context.Context, query string, opts types.SearchOptions) ([]types.PackageRef, error)
		ListInstalled(ctx context.Context, opts types.ListOptions) ([]types.InstalledPackage, error)
	}
}

// convertError converts internal error types to public error types.
func convertError(err error) error {
	if err == nil {
		return nil
	}

	// Convert base errors
	if err == types.ErrNotSupported {
		return ErrNotSupported
	}
	if err == types.ErrNotAvailable {
		return ErrNotAvailable
	}

	// Convert wrapped errors
	if types.IsNotSupported(err) {
		var notSupportedErr *types.NotSupportedError
		if errors.As(err, &notSupportedErr) {
			return &NotSupportedError{
				Operation: Operation(notSupportedErr.Operation),
				Backend:   notSupportedErr.Backend,
				Reason:    notSupportedErr.Reason,
			}
		}
		return ErrNotSupported
	}

	if types.IsNotAvailable(err) {
		var notAvailableErr *types.NotAvailableError
		if errors.As(err, &notAvailableErr) {
			return &NotAvailableError{
				Backend: notAvailableErr.Backend,
				Reason:  notAvailableErr.Reason,
			}
		}
		return ErrNotAvailable
	}

	if types.IsExternalFailure(err) {
		var extFailErr *types.ExternalFailureError
		if errors.As(err, &extFailErr) {
			return &ExternalFailureError{
				Operation: Operation(extFailErr.Operation),
				Backend:   extFailErr.Backend,
				Stdout:    extFailErr.Stdout,
				Stderr:    extFailErr.Stderr,
				Payload:   extFailErr.Payload,
				Err:       extFailErr.Err,
			}
		}
	}

	// Return error as-is if not a known type
	return err
}

func (a *backendAdapter) Available(ctx context.Context) (bool, error) {
	available, err := a.backend.Available(ctx)
	return available, convertError(err)
}

func (a *backendAdapter) Capabilities(ctx context.Context) ([]Capability, error) {
	caps, err := a.backend.Capabilities(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]Capability, len(caps))
	for i, c := range caps {
		result[i] = Capability{
			Operation: Operation(c.Operation),
			Supported: c.Supported,
			Notes:     c.Notes,
		}
	}
	return result, nil
}

func (a *backendAdapter) Update(ctx context.Context, opts UpdateOptions) (UpdateResult, error) {
	internalOpts := types.UpdateOptions{Progress: convertProgressReporter(opts.Progress)}
	res, err := a.backend.Update(ctx, internalOpts)
	var messages []ProgressMessage
	for _, m := range res.Messages {
		messages = append(messages, ProgressMessage{
			Severity:  Severity(m.Severity),
			Text:      m.Text,
			Timestamp: m.Timestamp,
			ActionID:  m.ActionID,
			TaskID:    m.TaskID,
			StepID:    m.StepID,
		})
	}
	return UpdateResult{Changed: res.Changed, Messages: messages}, convertError(err)
}

func (a *backendAdapter) Upgrade(ctx context.Context, opts UpgradeOptions) (UpgradeResult, error) {
	internalOpts := types.UpgradeOptions{Progress: convertProgressReporter(opts.Progress)}
	res, err := a.backend.Upgrade(ctx, internalOpts)
	var messages []ProgressMessage
	var pkgs []PackageRef
	for _, m := range res.Messages {
		messages = append(messages, ProgressMessage{
			Severity:  Severity(m.Severity),
			Text:      m.Text,
			Timestamp: m.Timestamp,
			ActionID:  m.ActionID,
			TaskID:    m.TaskID,
			StepID:    m.StepID,
		})
	}
	for _, p := range res.PackagesChanged {
		pkgs = append(pkgs, PackageRef{
			Name:      p.Name,
			Namespace: p.Namespace,
			Channel:   p.Channel,
			Kind:      p.Kind,
		})
	}
	return UpgradeResult{Changed: res.Changed, PackagesChanged: pkgs, Messages: messages}, convertError(err)
}

func (a *backendAdapter) Install(ctx context.Context, pkgs []PackageRef, opts InstallOptions) (InstallResult, error) {
	internalPkgs := make([]types.PackageRef, len(pkgs))
	for i, p := range pkgs {
		internalPkgs[i] = types.PackageRef{
			Name:      p.Name,
			Namespace: p.Namespace,
			Channel:   p.Channel,
			Kind:      p.Kind,
		}
	}
	internalOpts := types.InstallOptions{Progress: convertProgressReporter(opts.Progress)}
	res, err := a.backend.Install(ctx, internalPkgs, internalOpts)
	var messages []ProgressMessage
	var installed []PackageRef
	for _, m := range res.Messages {
		messages = append(messages, ProgressMessage{
			Severity:  Severity(m.Severity),
			Text:      m.Text,
			Timestamp: m.Timestamp,
			ActionID:  m.ActionID,
			TaskID:    m.TaskID,
			StepID:    m.StepID,
		})
	}
	for _, p := range res.PackagesInstalled {
		installed = append(installed, PackageRef{
			Name:      p.Name,
			Namespace: p.Namespace,
			Channel:   p.Channel,
			Kind:      p.Kind,
		})
	}
	return InstallResult{Changed: res.Changed, PackagesInstalled: installed, Messages: messages}, convertError(err)
}

func (a *backendAdapter) Uninstall(ctx context.Context, pkgs []PackageRef, opts UninstallOptions) (UninstallResult, error) {
	internalPkgs := make([]types.PackageRef, len(pkgs))
	for i, p := range pkgs {
		internalPkgs[i] = types.PackageRef{
			Name:      p.Name,
			Namespace: p.Namespace,
			Channel:   p.Channel,
			Kind:      p.Kind,
		}
	}
	internalOpts := types.UninstallOptions{Progress: convertProgressReporter(opts.Progress)}
	res, err := a.backend.Uninstall(ctx, internalPkgs, internalOpts)
	var messages []ProgressMessage
	var uninstalled []PackageRef
	for _, m := range res.Messages {
		messages = append(messages, ProgressMessage{
			Severity:  Severity(m.Severity),
			Text:      m.Text,
			Timestamp: m.Timestamp,
			ActionID:  m.ActionID,
			TaskID:    m.TaskID,
			StepID:    m.StepID,
		})
	}
	for _, p := range res.PackagesUninstalled {
		uninstalled = append(uninstalled, PackageRef{
			Name:      p.Name,
			Namespace: p.Namespace,
			Channel:   p.Channel,
			Kind:      p.Kind,
		})
	}
	return UninstallResult{Changed: res.Changed, PackagesUninstalled: uninstalled, Messages: messages}, convertError(err)
}

func (a *backendAdapter) Search(ctx context.Context, query string, opts SearchOptions) ([]PackageRef, error) {
	internalOpts := types.SearchOptions{Progress: convertProgressReporter(opts.Progress)}
	internalRes, err := a.backend.Search(ctx, query, internalOpts)
	if err != nil {
		return nil, convertError(err)
	}
	result := make([]PackageRef, len(internalRes))
	for i, p := range internalRes {
		result[i] = PackageRef{
			Name:      p.Name,
			Namespace: p.Namespace,
			Channel:   p.Channel,
			Kind:      p.Kind,
		}
	}
	return result, nil
}

func (a *backendAdapter) ListInstalled(ctx context.Context, opts ListOptions) ([]InstalledPackage, error) {
	internalOpts := types.ListOptions{Progress: convertProgressReporter(opts.Progress)}
	internalRes, err := a.backend.ListInstalled(ctx, internalOpts)
	if err != nil {
		return nil, convertError(err)
	}
	result := make([]InstalledPackage, len(internalRes))
	for i, p := range internalRes {
		result[i] = InstalledPackage{
			Ref: PackageRef{
				Name:      p.Ref.Name,
				Namespace: p.Ref.Namespace,
				Channel:   p.Ref.Channel,
				Kind:      p.Ref.Kind,
			},
			Version: p.Version,
			Status:  p.Status,
		}
	}
	return result, nil
}

// convertProgressReporter wraps a pm.ProgressReporter to be a types.ProgressReporter.
func convertProgressReporter(pr ProgressReporter) types.ProgressReporter {
	if pr == nil {
		return nil
	}
	return &progressReporterAdapter{pr: pr}
}

type progressReporterAdapter struct {
	pr ProgressReporter
}

func (a *progressReporterAdapter) OnAction(action types.ProgressAction) {
	a.pr.OnAction(ProgressAction{
		ID:        action.ID,
		Name:      action.Name,
		StartedAt: action.StartedAt,
		EndedAt:   action.EndedAt,
	})
}

func (a *progressReporterAdapter) OnTask(task types.ProgressTask) {
	a.pr.OnTask(ProgressTask{
		ID:        task.ID,
		ActionID:  task.ActionID,
		Name:      task.Name,
		StartedAt: task.StartedAt,
		EndedAt:   task.EndedAt,
	})
}

func (a *progressReporterAdapter) OnStep(step types.ProgressStep) {
	a.pr.OnStep(ProgressStep{
		ID:        step.ID,
		TaskID:    step.TaskID,
		Name:      step.Name,
		StartedAt: step.StartedAt,
		EndedAt:   step.EndedAt,
	})
}

func (a *progressReporterAdapter) OnMessage(msg types.ProgressMessage) {
	a.pr.OnMessage(ProgressMessage{
		Severity:  Severity(msg.Severity),
		Text:      msg.Text,
		Timestamp: msg.Timestamp,
		ActionID:  msg.ActionID,
		TaskID:    msg.TaskID,
		StepID:    msg.StepID,
	})
}

// NewBrew creates a new Brew backend that implements Manager and other interfaces.
func NewBrew(opts ...ConstructorOption) Manager {
	cfg := &backendConfig{}
	for _, opt := range opts {
		opt(cfg)
	}

	return &backendAdapter{
		backend: brew.New(nil, runner.NewRealRunner(), convertProgressReporter(cfg.progress)),
	}
}

// NewFlatpak creates a new Flatpak backend that implements Manager and other interfaces.
func NewFlatpak(opts ...ConstructorOption) Manager {
	cfg := &backendConfig{}
	for _, opt := range opts {
		opt(cfg)
	}

	return &backendAdapter{
		backend: flatpak.New(nil, convertProgressReporter(cfg.progress)),
	}
}

// NewSnap creates a new Snap backend that implements Manager and other interfaces.
func NewSnap(opts ...ConstructorOption) Manager {
	cfg := &backendConfig{}
	for _, opt := range opts {
		opt(cfg)
	}

	return &backendAdapter{
		backend: snap.New(nil, runner.NewRealRunner(), convertProgressReporter(cfg.progress)),
	}
}
