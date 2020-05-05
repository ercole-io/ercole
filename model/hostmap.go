package model

// HostMap contains info about the host
type HostMap map[string]interface{}

// CPUCores getter
func (host *HostMap) CPUCores() int {
	switch val := (*host)["CPUCores"].(type) {
	case int:
		return val
	case float64:
		return int(val)
	default:
		panic("Invalid CPUCores type")
	}

}
