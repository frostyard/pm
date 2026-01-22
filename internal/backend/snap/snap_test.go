package snap

import (
	"context"
	"testing"

	"github.com/frostyard/pm/internal/types"
)

func TestBackend_Available(t *testing.T) {
	t.Run("Returns NotAvailable when API is unreachable", func(t *testing.T) {
		b := New(nil, nil, nil)
		ctx := context.Background()

		available, err := b.Available(ctx)
		if available {
			t.Error("Expected Available() to return false when API is unreachable")
		}
		if !types.IsNotAvailable(err) {
			t.Errorf("Expected NotAvailable error, got %v", err)
		}
	})
}

func TestBackend_Capabilities(t *testing.T) {
	b := New(nil, nil, nil)
	ctx := context.Background()

	caps, err := b.Capabilities(ctx)
	if err != nil {
		t.Fatalf("Capabilities() error = %v", err)
	}

	if caps == nil {
		t.Fatal("Capabilities() returned nil, expected non-nil slice")
	}

	// Verify all operations are marked as not supported
	for _, cap := range caps {
		if cap.Supported {
			t.Errorf("Expected %s to be unsupported, but it's marked as supported", cap.Operation)
		}
	}
}

func TestBackend_EmptyMethods(t *testing.T) {
	b := New(nil, nil, nil)
	ctx := context.Background()

	t.Run("Update", func(t *testing.T) {
		_, err := b.Update(ctx, types.UpdateOptions{})
		if !types.IsNotSupported(err) {
			t.Errorf("Update should return NotSupported, got %v", err)
		}
	})

	t.Run("Upgrade", func(t *testing.T) {
		_, err := b.Upgrade(ctx, types.UpgradeOptions{})
		if !types.IsNotSupported(err) {
			t.Errorf("Upgrade should return NotSupported, got %v", err)
		}
	})

	t.Run("Install", func(t *testing.T) {
		_, err := b.Install(ctx, []types.PackageRef{{Name: "test"}}, types.InstallOptions{})
		if !types.IsNotSupported(err) {
			t.Errorf("Install should return NotSupported, got %v", err)
		}
	})

	t.Run("Uninstall", func(t *testing.T) {
		_, err := b.Uninstall(ctx, []types.PackageRef{{Name: "test"}}, types.UninstallOptions{})
		if !types.IsNotSupported(err) {
			t.Errorf("Uninstall should return NotSupported, got %v", err)
		}
	})

	t.Run("Search", func(t *testing.T) {
		_, err := b.Search(ctx, "test", types.SearchOptions{})
		if !types.IsNotSupported(err) {
			t.Errorf("Search should return NotSupported, got %v", err)
		}
	})

	t.Run("ListInstalled", func(t *testing.T) {
		_, err := b.ListInstalled(ctx, types.ListOptions{})
		if !types.IsNotSupported(err) {
			t.Errorf("ListInstalled should return NotSupported, got %v", err)
		}
	})
}
