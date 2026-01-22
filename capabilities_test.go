package pm

import "testing"

func TestSupports(t *testing.T) {
	caps := []Capability{
		{Operation: OperationSearch, Supported: true},
		{Operation: OperationInstall, Supported: false, Notes: "not implemented"},
		{Operation: OperationListInstalled, Supported: true},
	}

	tests := []struct {
		name string
		op   Operation
		want bool
	}{
		{
			name: "supported operation",
			op:   OperationSearch,
			want: true,
		},
		{
			name: "unsupported operation",
			op:   OperationInstall,
			want: false,
		},
		{
			name: "missing operation",
			op:   OperationUpgradePackages,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Supports(caps, tt.op); got != tt.want {
				t.Errorf("Supports() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCapability(t *testing.T) {
	caps := []Capability{
		{Operation: OperationSearch, Supported: true},
		{Operation: OperationInstall, Supported: false, Notes: "not implemented"},
	}

	tests := []struct {
		name    string
		op      Operation
		wantNil bool
	}{
		{
			name:    "existing capability",
			op:      OperationSearch,
			wantNil: false,
		},
		{
			name:    "missing capability",
			op:      OperationUpgradePackages,
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetCapability(caps, tt.op)
			if (got == nil) != tt.wantNil {
				t.Errorf("GetCapability() nil = %v, wantNil %v", got == nil, tt.wantNil)
			}
			if got != nil && got.Operation != tt.op {
				t.Errorf("GetCapability() operation = %v, want %v", got.Operation, tt.op)
			}
		})
	}
}
