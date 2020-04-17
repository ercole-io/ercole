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
	return count.(int)
}
