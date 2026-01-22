package pm

// Supports checks if a list of capabilities supports a specific operation.
func Supports(caps []Capability, op Operation) bool {
	for _, c := range caps {
		if c.Operation == op && c.Supported {
			return true
		}
	}
	return false
}

// GetCapability returns the capability for a specific operation, or nil if not found.
func GetCapability(caps []Capability, op Operation) *Capability {
	for i := range caps {
		if caps[i].Operation == op {
			return &caps[i]
		}
	}
	return nil
}
