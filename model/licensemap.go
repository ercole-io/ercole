package model

// LicenseMap holds information about Oracle database license
type LicenseMap map[string]interface{}

// Name getter
func (license *LicenseMap) Name() string {
	name := (*license)["Name"]
	return name.(string)
}

// Count getter
func (license *LicenseMap) Count() int {
	count := (*license)["Count"]
	switch val := count.(type) {
	case int:
		return val
	case float64:
		return int(val)
	default:
		panic("Invalid type")
	}
}
