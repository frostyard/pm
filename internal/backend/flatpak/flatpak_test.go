package flatpak

import (
	"context"
	"testing"

	"github.com/frostyard/pm/internal/types"
)

// mockRunner is a test double for runner.Runner
type mockRunner struct {
	stdout string
	stderr string
	err    error
}

func (m *mockRunner) Run(ctx context.Context, name string, args ...string) (string, string, error) {
	return m.stdout, m.stderr, m.err
}

func TestBackend_Available(t *testing.T) {
	t.Run("Returns NotAvailable when runner is nil", func(t *testing.T) {
		b := New(nil, nil)
		ctx := context.Background()

		available, err := b.Available(ctx)
		if available {
			t.Error("Expected Available() to return false with nil runner")
		}
		if !types.IsNotAvailable(err) {
			t.Errorf("Expected NotAvailable error, got %v", err)
		}
	})
}

func TestBackend_Capabilities(t *testing.T) {
	b := New(nil, nil)
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
	b := New(nil, nil)
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

func TestBackend_ListInstalled(t *testing.T) {
	t.Run("Parses namespace from installation column", func(t *testing.T) {
		mockRnr := &mockRunner{
			stdout: "Discord\tcom.discordapp.Discord\t0.0.121\tsystem\n" +
				"Rnote\tcom.github.flxzt.rnote\t0.13.1\tuser\n" +
				"Flatseal\tcom.github.tchx84.Flatseal\t2.4.0\tuser\n" +
				"Extension Manager\tcom.mattjakeman.ExtensionManager\t0.6.5\tsystem\n",
		}

		b := New(mockRnr, nil)
		ctx := context.Background()

		packages, err := b.ListInstalled(ctx, types.ListOptions{})
		if err != nil {
			t.Fatalf("ListInstalled() error = %v", err)
		}

		if len(packages) != 4 {
			t.Fatalf("Expected 4 packages, got %d", len(packages))
		}

		// Verify first package (system installation)
		if packages[0].Ref.Name != "com.discordapp.Discord" {
			t.Errorf("Expected name 'com.discordapp.Discord', got '%s'", packages[0].Ref.Name)
		}
		if packages[0].Ref.Namespace != "system" {
			t.Errorf("Expected namespace 'system', got '%s'", packages[0].Ref.Namespace)
		}
		if packages[0].Version != "0.0.121" {
			t.Errorf("Expected version '0.0.121', got '%s'", packages[0].Version)
		}

		// Verify second package (user installation)
		if packages[1].Ref.Name != "com.github.flxzt.rnote" {
			t.Errorf("Expected name 'com.github.flxzt.rnote', got '%s'", packages[1].Ref.Name)
		}
		if packages[1].Ref.Namespace != "user" {
			t.Errorf("Expected namespace 'user', got '%s'", packages[1].Ref.Namespace)
		}

		// Verify third package (user installation)
		if packages[2].Ref.Name != "com.github.tchx84.Flatseal" {
			t.Errorf("Expected name 'com.github.tchx84.Flatseal', got '%s'", packages[2].Ref.Name)
		}
		if packages[2].Ref.Namespace != "user" {
			t.Errorf("Expected namespace 'user', got '%s'", packages[2].Ref.Namespace)
		}

		// Verify fourth package (system installation)
		if packages[3].Ref.Name != "com.mattjakeman.ExtensionManager" {
			t.Errorf("Expected name 'com.mattjakeman.ExtensionManager', got '%s'", packages[3].Ref.Name)
		}
		if packages[3].Ref.Namespace != "system" {
			t.Errorf("Expected namespace 'system', got '%s'", packages[3].Ref.Namespace)
		}
	})

	t.Run("Handles missing installation column gracefully", func(t *testing.T) {
		mockRnr := &mockRunner{
			stdout: "Discord\tcom.discordapp.Discord\t0.0.121\n",
		}

		b := New(mockRnr, nil)
		ctx := context.Background()

		packages, err := b.ListInstalled(ctx, types.ListOptions{})
		if err != nil {
			t.Fatalf("ListInstalled() error = %v", err)
		}

		if len(packages) != 1 {
			t.Fatalf("Expected 1 package, got %d", len(packages))
		}

		// Should still parse name and version, but namespace will be empty
		if packages[0].Ref.Name != "com.discordapp.Discord" {
			t.Errorf("Expected name 'com.discordapp.Discord', got '%s'", packages[0].Ref.Name)
		}
		if packages[0].Version != "0.0.121" {
			t.Errorf("Expected version '0.0.121', got '%s'", packages[0].Version)
		}
	})
}
